package alertmgr

import (
	"io/ioutil"
	"log"
	"plugins"
	"scanservice"
	"sync"

	"github.com/ghodss/yaml"
	"utils"
)

const (
	IssueTypeDefault = "Task"
	PriorityDefault = "High"
)

type PluginSettings struct {
	Name            string `json:"name"`
	Enable          bool   `json:"enable"`
	Url             string `json:"url"`
	User            string `json:"user"`
	Password        string `json:"password"`
	TlsVerify       bool   `json:"tls_verify"`
	ProjectKey      string `json:"project_key,omitempty" structs:"project_key,omitempty"`
	IssueType       string `json:"issuetype" structs:"issuetype"`
	BoardName       string `json:"board,omitempty" structs:"board,omitempty"`
	Priority        string `json:"priority,omitempty"`
	Assignee        string `json:"assignee,omitempty"`
	Description     string
	Summary         string            `json:"summary,omitempty"`
	FixVersions     []string          `json:"fixVersions,omitempty"`
	AffectsVersions []string          `json:"affectsVersions,omitempty"`
	Labels          []string          `json:"labels,omitempty"`
	Sprint          string            `json:"sprint,omitempty"`
	Unknowns        map[string]string `json:"unknowns" structs:"unknowns,omitempty"`

	Host string `json:"host"`
	Port string `json:"port"`
	Recipients []string `json:"recipients"`
	Sender string `json:"sender"`

	PolicyMinVulnerability string `json:"Policy-Min-Vulnerability"`
	PolicyRegistry []string `json:"Policy-Registry"`
	PolicyImageName []string `json:"Policy-Image-Name"`
	PolicyNonCompliant bool `json:"Policy-Non-Compliant"`
}

type AlertMgr struct {
	mutex   sync.Mutex
	quit    chan struct{}
	queue   chan string
	cfgfile string
	plugins map[string]plugins.Plugin
	serviceSetting scanservice.ScanSettings
}

var initCtx sync.Once
var alertmgrCtx *AlertMgr

func NewEmailPlugin(settings PluginSettings) *plugins.EmailPlugin {
	em := new(plugins.EmailPlugin)
	em.User = settings.User
	em.Password = settings.Password
	em.Host = settings.Host
	em.Port = settings.Port
	em.Recipients = settings.Recipients
	em.Sender = settings.Sender
	return em
}

func NewJiraAPI(settings PluginSettings) *plugins.JiraAPI {
	jiraApi := &plugins.JiraAPI{
		Url:                    settings.Url,
		User:                   settings.User,
		Password:               settings.Password,
		TlsVerify:              settings.TlsVerify,
		Issuetype:              settings.IssueType,
		ProjectKey:             settings.ProjectKey,
		Priority:               settings.Priority,
		Assignee:               settings.Assignee,
		FixVersions:            settings.FixVersions,
		AffectsVersions:        settings.AffectsVersions,
		Labels:                 settings.Labels,
		Unknowns:               settings.Unknowns,
		SprintName:             settings.Sprint,
		SprintId:               -1,
		BoardName:              settings.BoardName,
		PolicyMinVulnerability: settings.PolicyMinVulnerability,
		PolicyRegistry:         settings.PolicyRegistry,
		PolicyImageName:        settings.PolicyImageName,
		PolicyNonCompliant:     settings.PolicyNonCompliant,
	}
	if jiraApi.Issuetype == "" {
		jiraApi.Issuetype = IssueTypeDefault
	}

	if jiraApi.Priority == "" {
		jiraApi.Priority = PriorityDefault
	}

	if jiraApi.Assignee == "" {
		jiraApi.Assignee = jiraApi.User
	}
	return jiraApi
}

func Instance() *AlertMgr {
	initCtx.Do(func() {
		alertmgrCtx = &AlertMgr{
			mutex:   sync.Mutex{},
			quit:    make(chan struct{}),
			queue:   make(chan string, 1000),
			plugins: make(map[string]plugins.Plugin),
		}
	})
	return alertmgrCtx
}

func (ctx *AlertMgr) Start(cfgfile string) {
	log.Printf("Starting AlertMgr....")
	ctx.cfgfile = cfgfile
	ctx.load()
	go ctx.listen()
}

func (ctx *AlertMgr) Terminate() {
	log.Printf("Terminating AlertMgr....")
	close(ctx.quit)
	for _, plugin := range ctx.plugins {
		if plugin != nil {
			plugin.Terminate()
		}
	}
}

func (ctx *AlertMgr) Send(data string) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	ctx.queue <- data
}

func (ctx *AlertMgr) load() error {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()
	log.Printf("Loading alerts configuration file %s ....\n", ctx.cfgfile)
	data, err := ioutil.ReadFile(ctx.cfgfile)
	if err != nil {
		log.Printf("Failed to open file %s, %s", ctx.cfgfile, err)
		return err
	}
	pluginSettings := []PluginSettings{}
	err = yaml.Unmarshal(data, &pluginSettings)
	if err != nil {
		log.Printf("Failed yaml.Unmarshal, %s", err)
		return err
	}
	for name, plugin := range ctx.plugins {
		if plugin != nil {
			ctx.plugins[name] = nil
			plugin.Terminate()
		}
	}
	for _, settings := range pluginSettings {
		utils.Debug("%#v\n", settings)
		if settings.Enable {
			utils.Debug("Starting Plugin %s\n", settings.Name)
			switch settings.Name {
			case "jira":
				plugin := NewJiraAPI(settings)
				plugin.Init()
				ctx.plugins["jira"] = plugin
			case "email":
				plugin := NewEmailPlugin(settings)
				plugin.Init()
				ctx.plugins["email"] = plugin
			}
		}
	}
	return nil
}

func (ctx *AlertMgr) listen() {
	for {
		select {
		case <-ctx.quit:
			return
		case data := <-ctx.queue:
			service := new(scanservice.ScanService)
			go service.ResultHandling(data, ctx.serviceSetting, ctx.plugins)
		}
	}
}
