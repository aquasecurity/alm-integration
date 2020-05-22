package scanservice

import (
	"dbservice"
	"formatting"
	"layout"
	"os"
	"plugins"
	"settings"
	"testing"
)

func TestRemoveLowLevelVulnerabilities(t *testing.T) {
	var tests = []struct{
		input      string
		severities map[string]bool
	} {
		/*{
			string(Bkmorrow),
			map[string]int{
				"critical":0,
				"high":0,
				"medium":2,
				"low":2,
				"negligible":2,
			},
		},
		 */
		{
			string(AlpineImageSource),
			map[string]bool{
				"critical": false,
				"high": false,
				"medium":true,
				"low":true,
				"negligible":true,
			},
		},
	}

	dbPathReal := dbservice.DbPath
	defer func() {
		dbservice.DbPath = dbPathReal
	}()
	dbservice.DbPath = "test_" + dbPathReal

	setting1 :=  &settings.Settings{
		PolicyMinVulnerability: "",
		PolicyRegistry:         nil,
		PolicyImageName:        nil,
		PolicyNonCompliant:     false,
		IgnoreRegistry:         nil,
		IgnoreImageName:        nil,
	}

	demoWithSettings := &DemoPlugin{
		name: "Demo plugin with settings",
		lay:  new(formatting.HtmlProvider),
		sets: setting1,
		t:    t,
	}

	for _, test := range tests {
		for severity, needSending := range test.severities {
			setting1.PolicyMinVulnerability = severity
			plgs := map[string]plugins.Plugin {}
			demoWithSettings.Sent = false
			plgs["demoSettings"] = demoWithSettings

			service := new(ScanService)
			service.ResultHandling( test.input,  plgs )

			if needSending != demoWithSettings.Sent {
				t.Errorf("The notify was sent with wrong severity %q for %q\n",
					severity, service.scanInfo.GetUniqueId())
			}
			os.Remove(dbservice.DbPath)
		}
	}

	demoWithoutSettings := &DemoPlugin{
		name: "Demo without settings",
		lay:   new(formatting.JiraLayoutProvider),
		sets: nil,
		t:    t,
	}
	for _, test := range tests {
		for range test.severities {
			plgs := map[string]plugins.Plugin {}
			demoWithoutSettings.Sent = false
			plgs["demoWithoutSettings"]= demoWithoutSettings
			service := new(ScanService)
			service.ResultHandling( test.input,  plgs )
			if !demoWithoutSettings.Sent {
				t.Errorf("The notify wasn't sent for plugin without settings for %q\n",
					service.scanInfo.GetUniqueId())
			}
			os.Remove(dbservice.DbPath)
		}
	}
}

type DemoPlugin struct {
	Sent bool
	name string
	lay  layout.LayoutProvider
	sets *settings.Settings
	t    *testing.T
}
func (plg *DemoPlugin) Init() error {	return nil}
func (plg *DemoPlugin) Send(data map[string]string) error {
	plg.Sent = true
	plg.t.Logf("Sending data via %q\n", plg.name)
//	plg.t.Logf("Title: %q\n", data["title"])
//	plg.t.Logf("Description: %q\n", data["description"])
	return nil
}

func (plg *DemoPlugin) Terminate() error { return nil}
func (plg *DemoPlugin) GetLayoutProvider() layout.LayoutProvider {
	return plg.lay
}
func (plg *DemoPlugin) GetSettings() *settings.Settings {
	return plg.sets
}