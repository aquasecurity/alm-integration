package outputs

import (
	"encoding/json"
	"log"

	"github.com/aquasecurity/postee/formatting"
	"github.com/aquasecurity/postee/layout"
	servicenow "github.com/aquasecurity/postee/servicenow"
)

type ServiceNowOutput struct {
	Name           string
	User           string
	Password       string
	Instance       string
	Table          string
	layoutProvider layout.LayoutProvider
}

func (sn *ServiceNowOutput) GetName() string {
	return sn.Name
}

func (sn *ServiceNowOutput) Init() error {
	log.Printf("Starting ServiceNow output %q....", sn.Name)
	log.Printf("Your ServiceNow Table is %q on '%s.%s'", sn.Table, sn.Instance, servicenow.BaseServer)
	sn.layoutProvider = new(formatting.HtmlProvider)
	return nil
}

func (sn *ServiceNowOutput) Send(content map[string]string) error {
	log.Printf("Sending via ServiceNow %q", sn.Name)
	d := &servicenow.ServiceNowData{
		ShortDescription: content["title"],
		WorkNotes:        "[code]" + content["description"] + "[/code]",
	}
	body, err := json.Marshal(d)
	if err != nil {
		log.Println("ServiceNow Error:", err)
		return err
	}
	err = servicenow.InsertRecordToTable(sn.User, sn.Password, sn.Instance, sn.Table, body)
	if err != nil {
		log.Println("ServiceNow Error:", err)
		return err
	}
	log.Printf("Sending via ServiceNow %q was successful!", sn.Name)
	return nil
}

func (sn *ServiceNowOutput) Terminate() error {
	log.Printf("ServiceNow output %q terminated", sn.Name)
	return nil
}

func (sn *ServiceNowOutput) GetLayoutProvider() layout.LayoutProvider {
	return sn.layoutProvider
}
