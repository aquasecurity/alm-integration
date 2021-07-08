package router

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

var (
	regoRule = "package postee.slack"
)

func TestInitTemplate(t *testing.T) {
	savedGetHttpClient := getHttpClient
	getHttpClient = getMockedHttpClient
	defaultRegoFolder := "rego-templates"
	commonRegoFolder := defaultRegoFolder + "/common"
	testRego := defaultRegoFolder + "/rego1.rego"
	err := os.Mkdir(defaultRegoFolder, 0777)
	if err != nil {
		t.Fatalf("Can't create rego folder: %v", err)
	}
	err = os.Mkdir(commonRegoFolder, 0777)
	if err != nil {
		t.Fatalf("Can't create rego folder: %v", err)
	}

	err = ioutil.WriteFile(testRego, []byte(regoRule), 0644)

	if err != nil {
		t.Fatalf("Can't write rego: %v", err)
	}

	defer func() {
		os.Remove(testRego)
		os.Remove(commonRegoFolder)
		os.Remove(defaultRegoFolder)
		getHttpClient = savedGetHttpClient
	}()

	tests := []struct {
		template          *Template
		caseDesc          string
		expectedCls       string
		shouldReturnError bool
	}{
		{
			template: &Template{
				Name:               "legacy-html",
				LegacyScanRenderer: "html",
			},
			caseDesc:    "Legacy mode test",
			expectedCls: "*formatting.legacyScnEvaluator",
		},
		{
			template: &Template{
				Name:        "built-in",
				RegoPackage: "postee.slack",
			},
			caseDesc:    "Built-in rego package",
			expectedCls: "*regoservice.regoEvaluator",
		},
		{
			template: &Template{
				Name: "from-url",
				Url:  "http://localhost/slack.rego",
			},
			caseDesc:    "Loading rego from url",
			expectedCls: "*regoservice.regoEvaluator",
		},
		{
			template: &Template{
				Name: "not-found",
				Url:  "http://localhost/wrong.rego",
			},
			caseDesc:          "Loading rego from not existing url",
			expectedCls:       "*regoservice.regoEvaluator",
			shouldReturnError: true,
		},
		{
			template: &Template{
				Name: "from-invalid-url",
				Url:  "invalid-url",
			},
			caseDesc:          "Loading rego from invalid url",
			expectedCls:       "*regoservice.regoEvaluator",
			shouldReturnError: true,
		},
		{
			template: &Template{
				Name: "inline",
				Body: "package postee.inline",
			},
			caseDesc:    "Loading rego from yaml config",
			expectedCls: "*regoservice.regoEvaluator",
		},
	}
	for _, test := range tests {
		doInitTemplate(t, test.caseDesc, test.template, test.expectedCls, test.shouldReturnError)
	}

}
func doInitTemplate(t *testing.T, caseDesc string, template *Template, expectedCls string, shouldReturnError bool) {
	demoCtx := Instance()
	err := demoCtx.initTemplate(template)
	if err != nil && !shouldReturnError {
		t.Fatalf("[%s] Unexpected error: %v", caseDesc, err)
	}
	if err == nil && shouldReturnError {
		t.Fatalf("Test case [%s] should return an error", caseDesc)
	}
	if shouldReturnError {
		return
	}

	initialized, ok := demoCtx.templates[template.Name]
	if !ok {
		t.Fatalf("[%s] template %s is not initialized", caseDesc, template.Name)
	}
	actualCls := fmt.Sprintf("%T", initialized)
	if actualCls != expectedCls {
		t.Errorf("[%s] Unexpected type of input evaluator. Expected %s, got %s \n", caseDesc, expectedCls, actualCls)
	}
}

//stuff for mocking http requests goes below

func getMockedHttpClient() *http.Client {
	return NewTestClient(responseWithRego)
}

// RoundTripFunc
type RoundTripFunc func(req *http.Request) (*http.Response, error)

// RoundTrip
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { //this is kind of wrapper where original function is used in interface implementation
	return f(req)
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}
func responseWithRego(req *http.Request) (*http.Response, error) {
	if req.URL.String() == "http://localhost/wrong.rego" {
		fmt.Printf("response status is %d\n", 404)
		return newTestResponse(404, "<html>not found</html>"), nil
	} else if req.URL.String() == "invalid-url" {
		return nil, errors.New("invalid url")
	} else {
		fmt.Printf("response status is %d\n", 200)
		return newTestResponse(200, "package custom1"), nil
	}
}

func newTestResponse(status int, response string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(bytes.NewBufferString(response)),
		// Must be set to non-nil value or it panics
		Header: make(http.Header),
	}
}
