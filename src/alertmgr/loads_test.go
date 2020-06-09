package alertmgr

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLoads(t *testing.T) {
	cfgData := `
---
- AquaServer: https://demolab.aquasec.com
- name: jira
  type: jira
  enable: true
  url: "http://localhost:2990/jira"
  user: admin
  password: admin
  tls_verify: false
  project_key: KEY
  description:
  summary:
  issuetype: "Bug"
  priority: Medium
  assignee: 
  Policy-Min-Vulnerability: Critical
  labels: ["label1", "label2"]
  Policy-Min-Vulnerability: high

- name: my-slack
  type: slack
  enable: true
  url: "https://hooks.slack.com/services/TT/BBB/WWWW"

- name: email
  type: email
  enable: true
  user: EMAILUSER
  password: EMAILPASS
  host: smtp.gmail.com
  port: 587
  recipients: ["demo@gmail.com"]

- name: ms-team
  type: teams
  enable: true
  url: https://outlook.office.com/webhook/.... # Webhook's url

- name: faild
  enable: true
  type: nextplugin

- name: my-servicenow
  type: serviceNow
  enable: true
  user: SERVICENOWUSER
  password: SERVICENOWPASS
  instance: dev00000
  board: incident
`
	cfgName :="cfg_test.yaml"
	ioutil.WriteFile(cfgName, []byte(cfgData),0644)
	defer func() {
		os.Remove(cfgName)
	}()


	demoCtx := Instance()
	demoCtx.Start(cfgName)
	pluginsNumber := 5
	if len(demoCtx.plugins) != pluginsNumber {
		t.Errorf("There are stopped plugins\nWaited: %d\nResult: %d", pluginsNumber, len(demoCtx.plugins))
	}

	_, ok := demoCtx.plugins["ms-team"]
	if !ok {
		t.Errorf("'ms-team' plugin didn't start!")
	}

	aquaWaiting := "https://demolab.aquasec.com/#/images/"
	if aquaServer != aquaWaiting {
		t.Errorf("Wrong init of AquaServer link.\nWait: %q\nGot: %q", aquaWaiting, aquaServer)
	}

	if _, ok := demoCtx.plugins["my-servicenow"]; !ok {
		t.Errorf("Plugin 'my-servicenow' didn't run!")
	}
}
