package formatting

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

var (
	scan1 = `{
		"image":"Demo mock image1",
		"registry":"registry1",
		"digest":"abc",
		"vulnerability_summary":{"critical":0,"high":1,"medium":3,"low":4,"negligible":5},
		"image_assurance_results":{"disallowed":true}
	}`
)

func TestEval(t *testing.T) {
	expectedTitle := "Demo mock image1 vulnerability scan report"
	expectedDescription := `<p>Image name: Demo mock image1</p>
<p>Registry: registry1</p>
<p>Image is non-compliant</p>
<TABLE border='1' style='width: 100%; border-collapse: collapse;'>
<TR>
<TH style='padding: 5px;'>CRITICAL</TH><TH style='padding: 5px;'>HIGH</TH><TH style='padding: 5px;'>MEDIUM</TH><TH style='padding: 5px;'>LOW</TH><TH style='padding: 5px;'>NEGLIGIBLE</TH>
</TR>
<TR>
<TD style='padding: 5px;'><span style='color:#c00000'>0</span></TD><TD style='padding: 5px;'><span style='color:#e0443d'>1</span></TD><TD style='padding: 5px;'><span style='color:#f79421'>3</span></TD><TD style='padding: 5px;'><span style='color:#e1c930'>4</span></TD><TD style='padding: 5px;'><span style='color:green'>5</span></TD>
</TR>
</TABLE>
<p>See more: <a href=''></a></p>
`

	in := map[string]interface{}{}
	if err := json.Unmarshal([]byte(scan1), &in); err != nil {
		t.Fatalf("json.Unmarshal error for %s: %v\n", scan1, err)
	}
	e, err := BuildLegacyScnEvaluator("html")

	if err != nil {
		t.Fatalf("Unexpected error %v\n", err)
	}
	out, err := e.Eval(in, "")

	if out["title"] != expectedTitle {
		t.Errorf("Unexpected title value got %s, expected %s\n", out["title"], expectedTitle)
	}
	if out["description"] != expectedDescription {
		t.Errorf("Unexpected description value got %s, expected %s\n", out["description"], expectedDescription)
	}
}

func TestAggregationSupport(t *testing.T) {
	e := &legacyScnEvaluator{}
	if !e.IsAggregationSupported() {
		t.Errorf("Legacy Scan Evaluator should support aggregation by default\n")
	}
}

func TestBuildAggregatedContent(t *testing.T) {

	expectedTitle := "Vulnerability scan report"
	expectedDescription := `<h1>title1</h1>
description1<h1>title2</h1>
description2`

	expectedUrl := `url1
url2`
	expectedOwners := []string{"admin", "user"}

	e, err := BuildLegacyScnEvaluator("html")

	if err != nil {
		t.Fatalf("Unexpected error %v\n", err)
	}
	in := []map[string]string{
		{
			"title":       "title1",
			"description": "description1",
			"url":         "url1",
			"owners":      "admin",
		},
		{
			"title":       "title2",
			"description": "description2",
			"url":         "url2",
			"owners":      "user",
		},
	}
	out, err := e.BuildAggregatedContent(in)
	if err != nil {
		t.Fatalf("Unexpected error %v\n", err)
	}
	if out["title"] != expectedTitle {
		t.Errorf("Unexpected title value got %s, expected %s\n", out["title"], expectedTitle)
	}
	if out["description"] != expectedDescription {
		t.Errorf("Unexpected description value got %s, expected %s\n", out["description"], expectedDescription)
	}
	if out["url"] != expectedUrl {
		t.Errorf("Unexpected description value got %s, expected %s\n", out["url"], expectedUrl)
	}
	actualOwners := strings.Split(out["owners"], ";")
	if len(actualOwners) == len(expectedOwners) {
		for _, own := range actualOwners {
			found := false
			for _, expOwn := range expectedOwners {
				if own == expOwn {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Unexpected owner: %s\n", own)

			}
		}
	} else {
		t.Errorf("Unexpected owners value got %s, expected %s\n", out["owners"], expectedOwners)
	}
}

func TestBuildLegacyScnEvaluator(t *testing.T) {
	tests := []struct {
		layoutType          string
		expectedLayoutClass string
		shouldFail          bool
	}{
		{"html", "*formatting.HtmlProvider", false},
		{"jira", "*formatting.JiraLayoutProvider", false},
		{"slack", "*formatting.SlackMrkdwnProvider", false},
		{"xml", "", true},
	}
	for _, test := range tests {
		e, err := BuildLegacyScnEvaluator(test.layoutType)
		if err == nil && test.shouldFail {
			t.Fatalf("BuildLegacyScnEvaluator should fail for layout type %s but actually didn't return an error\n", test.layoutType)
		} else if err != nil && !test.shouldFail {
			t.Fatalf("Unexpected error %v\n", err)
		}
		if test.shouldFail {
			return
		}
		scnEvaluator, ok := e.(*legacyScnEvaluator)
		if !ok {
			t.Fatalf("Unexpected type of evaluator returned %T\n", e)
		}
		actualCls := fmt.Sprintf("%T", scnEvaluator.layoutProvider)
		if actualCls != test.expectedLayoutClass {
			t.Errorf("Invalid type of layout provider returned, expected %s, got %s\n", test.expectedLayoutClass, actualCls)
		}
	}
}
