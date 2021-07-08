package msgservice

import (
	"log"
	"os"
	"testing"

	"github.com/aquasecurity/postee/data"
	"github.com/aquasecurity/postee/dbservice"
	"github.com/aquasecurity/postee/outputs"
	"github.com/aquasecurity/postee/routes"
)

func TestAggregateByTimeout(t *testing.T) {
	const aggregationSeconds = 3

	dbPathReal := dbservice.DbPath
	savedRunScheduler := RunScheduler
	schedulerInvctCnt := 0
	defer func() {
		os.Remove(dbservice.DbPath)
		dbservice.DbPath = dbPathReal
		RunScheduler = savedRunScheduler
	}()
	RunScheduler = func(
		route *routes.InputRoute,
		fnSend func(plg outputs.Output, cnt map[string]string),
		fnAggregate func(outputName string, currentContent map[string]string, counts int, ignoreLength bool) []map[string]string,
		inpteval data.Inpteval,
		name *string,
		output outputs.Output,
	) {
		log.Printf("Mocked Scheduler is activated for route %q. Period: %d sec", route.Name, route.Plugins.AggregateTimeoutSeconds)
		route.StartScheduler()

		schedulerInvctCnt++
	}

	dbservice.DbPath = "test_webhooks.db"

	demoRoute := &routes.InputRoute{
		Name: "demo-route1",
		Plugins: routes.Plugins{
			PolicyShowAll:           true,
			AggregateTimeoutSeconds: aggregationSeconds,
		},
	}

	demoEmailPlg := &DemoEmailOutput{}

	demoInptEval := &DemoInptEval{}

	srvUrl := ""

	srv1 := new(MsgService)
	srv1.MsgHandling([]byte(mockScan1), demoEmailPlg, demoRoute, demoInptEval, &srvUrl)
	srv1.MsgHandling([]byte(mockScan2), demoEmailPlg, demoRoute, demoInptEval, &srvUrl)
	srv1.MsgHandling([]byte(mockScan3), demoEmailPlg, demoRoute, demoInptEval, &srvUrl)

	expectedSchedulerInvctCnt := 1

	if schedulerInvctCnt != expectedSchedulerInvctCnt {
		t.Errorf("Unexpected plugin invocation count %d, expected %d \n", schedulerInvctCnt, expectedSchedulerInvctCnt)
	}

	demoRoute.StopScheduler()
}
