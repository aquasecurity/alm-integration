package msgservice

import (
	"sync"
	"testing"

	"github.com/aquasecurity/postee/outputs"
	"github.com/aquasecurity/postee/routes"
)

func TestSheduler(t *testing.T) {
	routeName := "test-schedule"
	demoRoute := &routes.InputRoute{}
	demoRoute.Name = routeName

	demoRoute.Plugins.AggregateTimeoutSeconds = 3

	demoSend := func(plg outputs.Output, cnt map[string]string) {
		plg.Send(cnt)
	}
	demoAggregate := func(outputName string, currentContent map[string]string, counts int, ignoreLength bool) []map[string]string {
		return []map[string]string{
			{
				"title":       "title1",
				"description": "description1",
			},
			{
				"title":       "title2",
				"description": "description2",
			},
		}
	}
	demoInptEval := &DemoInptEval{}

	demoEmailOutput := &DemoEmailOutput{
		emailCounts: 0,
	}

	demoEmailOutput.wg = &sync.WaitGroup{}
	demoEmailOutput.wg.Add(1)

	RunScheduler(demoRoute, demoSend, demoAggregate, demoInptEval, &routeName, demoEmailOutput)

	demoEmailOutput.wg.Wait()
	demoRoute.StopScheduler()

}
