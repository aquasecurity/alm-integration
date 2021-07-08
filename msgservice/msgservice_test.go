package msgservice

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/aquasecurity/postee/dbservice"
	"github.com/aquasecurity/postee/routes"
)

var (
	invalidJson = `{
	image : "My Image"
	}`
)

type FailingInptEval struct {
	expectedError     error
	expectedAggrError error
}

func (inptEval *FailingInptEval) Eval(in map[string]interface{}, serverUrl string) (map[string]string, error) {
	if inptEval.expectedError != nil {
		return nil, inptEval.expectedError
	} else {
		return map[string]string{
			"title":       "some title",
			"description": "some description",
		}, nil
	}
}
func (inptEval *FailingInptEval) BuildAggregatedContent(items []map[string]string) (map[string]string, error) {

	return nil, inptEval.expectedAggrError
}
func (inptEval *FailingInptEval) IsAggregationSupported() bool {
	return inptEval.expectedAggrError != nil
}

func TestInputs(t *testing.T) {
	tests := []struct {
		input      []byte
		caseDesc   string
		shouldPass bool
	}{
		{
			input:      nil,
			caseDesc:   "Empty input",
			shouldPass: false,
		},
		{
			input:      []byte(invalidJson),
			caseDesc:   "Invalid Json",
			shouldPass: false,
		},
	}
	for _, test := range tests {
		validateInputValue(t, test.caseDesc, test.input, test.shouldPass)
	}

}
func validateInputValue(t *testing.T, caseDesc string, input []byte, shouldPass bool) {
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
	expected := 0
	if shouldPass {
		expected = 1
	}

	demoRoute := &routes.InputRoute{}

	demoRoute.Name = "demo-route"

	demoInptEval := &DemoInptEval{}

	demoEmailOutput.wg = &sync.WaitGroup{}
	demoEmailOutput.wg.Add(expected)

	srv := new(MsgService)
	srv.MsgHandling([]byte(input), demoEmailOutput, demoRoute, demoInptEval, &srvUrl)

	demoEmailOutput.wg.Wait()

	if demoEmailOutput.getEmailsCount() != expected {
		t.Errorf("[%s] Wrong number of Send method calls: expected %d, got %d", caseDesc, expected, demoEmailOutput.getEmailsCount())
	}

}
func TestEvalError(t *testing.T) {
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
	evalError := errors.New("eval error")

	demoInptEval := &FailingInptEval{
		expectedError: evalError,
	}

	srv := new(MsgService)
	srv.MsgHandling([]byte(mockScan1), demoEmailOutput, demoRoute, demoInptEval, &srvUrl)

	if demoEmailOutput.getEmailsCount() > 0 {
		t.Errorf("Output shouldn't be called when evaluation is failed")
	}
}

func TestAggrEvalError(t *testing.T) {
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
	demoRoute.Plugins.AggregateIssuesNumber = 2
	demoRoute.Plugins.PolicyShowAll = true

	aggrEvalError := errors.New("aggregation eval error")

	demoInptEval := &FailingInptEval{
		expectedAggrError: aggrEvalError,
	}

	for i := 0; i < 2; i++ {
		srv := new(MsgService)
		srv.MsgHandling([]byte(mockScan1), demoEmailOutput, demoRoute, demoInptEval, &srvUrl)
	}

	if demoEmailOutput.getEmailsCount() > 0 {
		t.Errorf("Output shouldn't be called when evaluation is failed")
	}
}
func TestEmptyInput(t *testing.T) {
	dbPathReal := dbservice.DbPath
	defer func() {
		os.Remove(dbservice.DbPath)
		dbservice.DbPath = dbPathReal
	}()
	dbservice.DbPath = "test_webhooks.db"

	srvUrl := ""

	demoRoute := &routes.InputRoute{}

	demoRoute.Name = "demo-route"

	demoInptEval := &DemoInptEval{}

	srv := new(MsgService)
	srv.MsgHandling([]byte("{}"), nil, demoRoute, demoInptEval, &srvUrl)

	if demoInptEval.renderCnt != 0 {
		t.Errorf("Eval() shouldn't be called if no output is passed to ResultHandling()")
	}
}
