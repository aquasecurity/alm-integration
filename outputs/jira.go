package outputs

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/aquasecurity/postee/formatting"
	"github.com/aquasecurity/postee/layout"

	"net/http"
	"os"
	"strings"

	"github.com/andygrunwald/go-jira"
)

type JiraAPI struct {
	Name            string
	Url             string
	User            string
	Password        string
	TlsVerify       bool
	Issuetype       string
	ProjectKey      string
	Priority        string
	Assignee        []string
	Description     string
	Summary         string
	SprintName      string
	SprintId        int
	FixVersions     []string
	AffectsVersions []string
	Labels          []string
	Unknowns        map[string]string
	BoardName       string
	boardId         int
	boardType       string
}

func (ctx *JiraAPI) GetName() string {
	return ctx.Name
}

func (ctx *JiraAPI) fetchBoardId(boardName string) {
	client, err := ctx.createClient()
	if err != nil {
		log.Printf("unable to create Jira client: %s, please check your credentials.", err)
		return
	}

	boardlist, _, err := client.Board.GetAllBoards(&jira.BoardListOptions{ProjectKeyOrID: ctx.ProjectKey})
	if err != nil {
		log.Printf("failed to get boards from Jira API GetAllBoards with ProjectID %s. %s", ctx.ProjectKey, err)
		return
	}
	var matches int
	for _, board := range boardlist.Values {
		if board.Name == boardName {
			ctx.boardId = board.ID
			ctx.boardType = board.Type
			matches++
		}
	}

	if matches > 1 {
		log.Printf("found more than one boards with name %q, working with board id %d", boardName, ctx.boardId)
	} else if matches == 0 {
		log.Printf("no boards found with name %s when getting all boards for User", boardName)
		return
	} else {
		log.Printf("using board ID %d with Name %q", ctx.boardId, boardName)
	}
}

func (ctx *JiraAPI) fetchSprintId(client jira.Client) {
	sprints, _, err := client.Board.GetAllSprintsWithOptions(ctx.boardId, &jira.GetAllSprintsOptions{State: "active"})
	if err != nil {
		log.Printf("failed to get active sprint for board ID %d from Jira API. %s", ctx.boardId, err)
		return
	}
	if len(sprints.Values) > 1 {
		ctx.SprintId = len(sprints.Values) - 1
		log.Printf("Found more than one active sprint, using sprint id %d as the active sprint", ctx.SprintId)
	} else if len(sprints.Values) == 1 {
		if sprints.Values[0].ID != ctx.SprintId {
			ctx.SprintId = sprints.Values[0].ID
			log.Printf("using sprint id %d as the active sprint", ctx.SprintId)
		}
	} else {
		log.Printf("no active sprints exist in board ID %d Name %s", ctx.boardId, ctx.ProjectKey)
	}
}

func (ctx *JiraAPI) Terminate() error {
	log.Printf("Jira output terminated\n")
	return nil
}

func (ctx *JiraAPI) Init() error {
	if ctx.BoardName == "" {
		ctx.BoardName = fmt.Sprintf("%s board", ctx.ProjectKey)
	}
	ctx.fetchBoardId(ctx.BoardName)

	log.Printf("Starting Jira output %q....", ctx.Name)
	if len(ctx.Password) == 0 {
		ctx.Password = os.Getenv("JIRA_PASSWORD")
	}
	return nil
}

func (jira *JiraAPI) GetLayoutProvider() layout.LayoutProvider {
	return new(formatting.JiraLayoutProvider)
}

func (ctx *JiraAPI) createClient() (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: ctx.User,
		Password: ctx.Password,
	}

	if !ctx.TlsVerify {
		tp.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	client, err := jira.NewClient(tp.Client(), ctx.Url)
	if err != nil {
		return client, fmt.Errorf("unable to create new JIRA client. %v", err)
	}
	return client, nil
}

func (ctx *JiraAPI) Send(content map[string]string) error {
	client, err := ctx.createClient()
	if err != nil {
		log.Printf("unable to create Jira client: %s", err)
		return err
	}

	if ctx.boardType == "scrum" {
		ctx.fetchSprintId(*client)
	}

	metaProject, err := createMetaProject(client, ctx.ProjectKey)
	if err != nil {
		return fmt.Errorf("Failed to create meta project: %s\n", err)
	}

	metaIssueType, err := createMetaIssueType(metaProject, ctx.Issuetype)
	if err != nil {
		return fmt.Errorf("Failed to create meta issue type: %s", err)
	}

	ctx.Summary = content["title"]
	ctx.Description = content["description"]

	assignee := ctx.User
	if len(ctx.Assignee) > 0 {
		assignees := getHandledRecipients(ctx.Assignee, &content, ctx.Name)
		if len(assignees) > 0 {
			assignee = assignees[0]
		}
	}

	fieldsConfig := map[string]string{
		"Issue Type":  ctx.Issuetype,
		"Project":     ctx.ProjectKey,
		"Priority":    ctx.Priority,
		"Assignee":    assignee,
		"Description": ctx.Description,
		"Summary":     ctx.Summary,
	}
	if ctx.SprintId > 0 {
		fieldsConfig["Sprint"] = strconv.Itoa(ctx.SprintId)
	}

	//Add all custom fields that are unknown to fieldsConfig. Unknown are fields that are custom User defined in jira.
	for k, v := range ctx.Unknowns {
		fieldsConfig[k] = v
	}
	if len(ctx.Unknowns) > 0 {
		log.Printf("added %d custom fields to issue.", len(ctx.Unknowns))
	}

	type Version struct {
		Name string `json:"name"`
	}

	issue, err := InitIssue(client, metaProject, metaIssueType, fieldsConfig)
	if err != nil {
		log.Printf("Failed to init issue: %s\n", err)
		return err
	}

	if len(ctx.Labels) > 0 {
		for _, l := range ctx.Labels {
			issue.Fields.Labels = append(issue.Fields.Labels, l)
		}
	}

	if len(ctx.FixVersions) > 0 {
		for _, v := range ctx.FixVersions {
			issue.Fields.FixVersions = append(issue.Fields.FixVersions, &jira.FixVersion{
				Name: v,
			})
		}
	}

	if len(ctx.AffectsVersions) > 0 {
		affectsVersions := []*Version{}
		for _, v := range ctx.AffectsVersions {
			affectsVersions = append(affectsVersions, &Version{
				Name: v,
			})
		}
		issue.Fields.Unknowns["versions"] = affectsVersions
		log.Printf("added %d affected versions into Versions field", len(ctx.AffectsVersions))
	}

	i, err := ctx.openIssue(client, issue)
	if err != nil {
		log.Printf("Failed to open jira issue, %s\n", err)
		return err
	}
	log.Printf("Created new jira issue %s", i.ID)
	return nil
}

func (ctx *JiraAPI) login(client *jira.Client) error {
	_, err := client.Authentication.AcquireSessionCookie(ctx.User, ctx.Password)
	return err
}

func (ctx *JiraAPI) openIssue(client *jira.Client, issue *jira.Issue) (*jira.Issue, error) {
	i, res, err := client.Issue.Create(issue)

	defer res.Body.Close()
	resp, _ := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(string(resp))
	}
	return i, nil
}

func createMetaProject(c *jira.Client, project string) (*jira.MetaProject, error) {
	meta, _, err := c.Issue.GetCreateMeta(project)
	if err != nil {
		return nil, fmt.Errorf("failed to get create meta : %s", err)
	}

	// get right project
	metaProject := meta.GetProjectWithKey(project)
	if metaProject == nil {
		return nil, fmt.Errorf("could not find project with key %s", project)
	}

	return metaProject, nil
}

func createMetaIssueType(metaProject *jira.MetaProject, issueType string) (*jira.MetaIssueType, error) {
	metaIssuetype := metaProject.GetIssueTypeWithName(issueType)
	if metaIssuetype == nil {
		return nil, fmt.Errorf("could not find issuetype %s", issueType)
	}

	return metaIssuetype, nil
}

func InitIssue(c *jira.Client, metaProject *jira.MetaProject, metaIssuetype *jira.MetaIssueType, fieldsConfig map[string]string) (*jira.Issue, error) {
	issue := new(jira.Issue)
	issueFields := new(jira.IssueFields)
	issueFields.Unknowns = make(map[string]interface{})

	// map the field names the User presented to jira's internal key
	allFields, _ := metaIssuetype.GetAllFields()
	for key, value := range fieldsConfig {

		jiraKey, found := allFields[key]
		if !found {
			return nil, fmt.Errorf("key %s is not found in the list of fields", key)
		}

		valueType, err := metaIssuetype.Fields.String(jiraKey + "/schema/type")
		if err != nil {
			return nil, err
		}

		switch strings.ToLower(valueType) {
		case "array":
			// split value (string) into slice by delimiter
			elements := strings.Split(value, ",")

			elemType, err := metaIssuetype.Fields.String(jiraKey + "/schema/items")
			if err != nil {
				return nil, err
			}
			switch elemType {
			case "component":
				issueFields.Unknowns[jiraKey] = []jira.Component{{Name: value}}
			case "option":
				optionsMap := make([]map[string]string, 0)

				for _, element := range elements {
					optionsMap = append(optionsMap, map[string]string{"value": element})
				}
				issueFields.Unknowns[jiraKey] = optionsMap
			default:
				if key == "Sprint" {
					num, err := strconv.Atoi(value)
					if err != nil {
						return nil, err
					}
					issueFields.Unknowns[jiraKey] = num // Due to Jira REST API behavior, needed to specify not a slice but a number.
				} else {
					issueFields.Unknowns[jiraKey] = []string{value}
				}
			}
		case "number":
			val, err := strconv.Atoi(value)
			if err != nil {
				fmt.Printf("Failed convert value(string) to int: %s\n", err)
			}
			issueFields.Unknowns[jiraKey] = val

		// TODO: Handle Cascading Select List
		//case "option-with-child":
		//	type CustomField struct {
		//		Value string `json:"value"`
		//	}
		//	type CustomFieldCascading struct {
		//		Value string `json:"value"`
		//		Child CustomField `json:"child"`
		//	}
		//
		//	a := CustomFieldCascading{ Value: "1", Child: CustomField{Value: "a"}}

		case "string":
			issueFields.Unknowns[jiraKey] = value
		case "date":
			issueFields.Unknowns[jiraKey] = value
		case "datetime":
			issueFields.Unknowns[jiraKey] = value
		case "any":
			// Treat any as string
			issueFields.Unknowns[jiraKey] = value
		case "project":
			issueFields.Unknowns[jiraKey] = jira.Project{
				Name: metaProject.Name,
				ID:   metaProject.Id,
			}
		case "priority":
			issueFields.Unknowns[jiraKey] = jira.Priority{Name: value}
		case "user":
			users, resp, err := c.User.Find(value)
			if err != nil {
				log.Printf("Get Jira User info error: %v", err)
				continue
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("http response failed: %q", resp.Status)
				continue
			}
			if len(users) == 0 {
				log.Printf("There is no user for %q", value)
				continue
			}
			issueFields.Unknowns[jiraKey] = users[0]
		case "issuetype":
			issueFields.Unknowns[jiraKey] = jira.IssueType{
				Name: value,
			}
		case "option":
			issueFields.Unknowns[jiraKey] = jira.Option{
				Value: value,
			}

		default:
			return nil, fmt.Errorf("Unknown issue type encountered: %s for %s", valueType, key)
		}
	}
	issue.Fields = issueFields
	return issue, nil
}
