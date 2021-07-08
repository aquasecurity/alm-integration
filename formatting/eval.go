package formatting

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aquasecurity/postee/data"
	"github.com/aquasecurity/postee/layout"
)

type legacyScnEvaluator struct {
	layoutProvider layout.LayoutProvider
}

func (legacyScnEvaluator *legacyScnEvaluator) Eval(in map[string]interface{}, serverUrl string) (map[string]string, error) {
	scan, err := toScanImage(in)
	if err != nil {
		return nil, err
	}
	title := fmt.Sprintf("%s vulnerability scan report", in["image"])

	return map[string]string{

		"title":       title,
		"description": layout.GenTicketDescription(legacyScnEvaluator.layoutProvider, scan, nil, serverUrl),
		"url":         serverUrl,
	}, nil
}
func (legacyScnEvaluator *legacyScnEvaluator) IsAggregationSupported() bool {
	return true
}

func (legacyScnEvaluator *legacyScnEvaluator) BuildAggregatedContent(scans []map[string]string) (map[string]string, error) {
	var descr bytes.Buffer
	var urls bytes.Buffer
	owners := []string{}
	for _, scan := range scans {
		descr.WriteString(legacyScnEvaluator.layoutProvider.TitleH1(scan["title"]))
		descr.WriteString(scan["description"])
		if urls.Len() > 0 {
			urls.WriteByte('\n')
		}
		urls.WriteString(scan["url"])
		if len(scan["owners"]) > 0 {
			owners = append(owners, scan["owners"])
		}
	}
	title := "Vulnerability scan report"

	r := map[string]string{
		"title":       title,
		"description": descr.String(),
		"url":         urls.String(), //TODO this is strange ...
	}

	if len(owners) > 0 {
		r["owners"] = strings.Join(owners, ";")
	}
	return r, nil
}

func toScanImage(in map[string]interface{}) (*data.ScanImageInfo, error) {
	source, err := json.Marshal(in) //back to bytes
	if err != nil {
		return nil, err
	}

	scanInfo := new(data.ScanImageInfo)

	err = json.Unmarshal(source, scanInfo)
	if err != nil {
		return nil, err
	}
	return scanInfo, nil
}

func BuildLegacyScnEvaluator(layoutType string) (data.Inpteval, error) {
	switch layoutType {
	case "slack":
		return &legacyScnEvaluator{
			layoutProvider: &SlackMrkdwnProvider{},
		}, nil
	case "html":
		return &legacyScnEvaluator{
			layoutProvider: &HtmlProvider{},
		}, nil
	case "jira":
		return &legacyScnEvaluator{
			layoutProvider: &JiraLayoutProvider{},
		}, nil
	default:
		return nil, errors.New("unknown layout type")
	}
}
