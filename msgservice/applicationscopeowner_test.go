package msgservice

import (
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/aquasecurity/postee/dbservice"
	"github.com/aquasecurity/postee/routes"
)

var (
	scnWithOwners = `{
		"image":"Demo mock image1",
		"registry":"registry1",
		"vulnerability_summary":{"critical":0,"high":1,"medium":3,"low":4,"negligible":5},
		"image_assurance_results":{"disallowed":true},
		"application_scope_owners": ["recipient1@aquasec.com", "recipient1@aquasec.com"]
	}`
)

func TestApplicationScopeOwner(t *testing.T) {
	dbPathReal := dbservice.DbPath
	defer func() {
		os.Remove(dbservice.DbPath)
		dbservice.DbPath = dbPathReal
	}()
	dbservice.DbPath = "test_webhooks.db"

	demoEmailOutput := &DemoEmailOutput{
		emailCounts: 0,
	}

	srvUrl := ""
	demoRoute := &routes.InputRoute{}

	demoRoute.Name = "demo-route"
	demoRoute.Plugins.PolicyShowAll = true

	demoInptEval := &DemoInptEval{}

	demoEmailOutput.wg = &sync.WaitGroup{}
	demoEmailOutput.wg.Add(1)

	srv := new(MsgService)
	srv.MsgHandling([]byte(scnWithOwners), demoEmailOutput, demoRoute, demoInptEval, &srvUrl)

	demoEmailOutput.wg.Wait()

	if len(demoEmailOutput.payloads) != 1 {
		t.Errorf("Output Send method isn't called as expected! Number of invocation expected %d, got: %d", 1, len(demoEmailOutput.payloads))
	}
	sent := demoEmailOutput.payloads[0]

	ownersStr, ok := sent["owners"]
	if !ok {
		t.Errorf("Owners key is missed from output payload")
	}
	owners := strings.Split(ownersStr, ";")
	for _, own := range owners {
		if own != "recipient1@aquasec.com" && own != "recipient2@aquasec.com" {
			t.Errorf("Unexpected owner value: %s", own)
		}
	}
}
