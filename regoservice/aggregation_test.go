package regoservice

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	regoWithAggregation = `
package rego1
title:="Audit event received"
result:=sprintf("Audit event received from %s", [input.user])
aggregation_pkg:="rego1.aggr"
`
	aggregationRego = `
package rego1.aggr
import data.postee.flat_array


title := "Vulnerability scan report"
result := res {
    scans := [ scan | 
            item:=input[i].description

            scan:=[sprintf("<h1>%s</h1>", [input[i].title]), item]
    ] 

    res:= concat("\n", flat_array(scans))
}
`
)

func TestAggregation(t *testing.T) {
	tests := []struct {
		regoRule            *string
		aggregationRegoRule *string
		caseDesc            string
		items               []map[string]string
		regoPackage         string
		expectedValues      map[string]string
	}{
		{
			regoRule:            &regoWithAggregation,
			aggregationRegoRule: &aggregationRego,
			caseDesc:            "simple case",
			items: []map[string]string{{
				"title":       "title1",
				"description": "description1",
			}, {
				"title":       "title2",
				"description": "description2",
			}},
			regoPackage: "rego1",
			expectedValues: map[string]string{
				"title": "Vulnerability scan report",
				"description": `<h1>title1</h1>
description1
<h1>title2</h1>
description2`,
			},
		},
	}
	for _, test := range tests {
		aggregateBuildinRego(t, test.caseDesc, test.regoRule, test.aggregationRegoRule, test.items, test.regoPackage, test.expectedValues)
	}
}

func aggregateBuildinRego(t *testing.T, caseDesc string, regoRule *string, aggregationRegoRule *string, items []map[string]string, regoPackage string, expectedValues map[string]string) {
	buildinRegoTemplatesSaved := buildinRegoTemplates
	testRego := "rego1.rego"
	aggrRego := "aggr1.rego"
	commonRegoFilename := "common.rego"
	buildinRegoTemplates = []string{commonRegoFilename, testRego, aggrRego} //common part goes in single bundle

	ioutil.WriteFile(commonRegoFilename, []byte(commonRego), 0644)
	ioutil.WriteFile(testRego, []byte(*regoRule), 0644)
	ioutil.WriteFile(aggrRego, []byte(*aggregationRegoRule), 0644)

	defer func() {
		buildinRegoTemplates = buildinRegoTemplatesSaved
		os.Remove(testRego)
		os.Remove(commonRegoFilename)
		os.Remove(aggrRego)
	}()
	demo, err := BuildBundledRegoEvaluator(regoPackage)
	if err != nil {
		t.Errorf("received an unexpected error: %v\n", err)
	}
	if !demo.IsAggregationSupported() {
		t.Errorf("Should support aggregation")
		return
	}
	r, err := demo.BuildAggregatedContent(items)
	if err != nil {
		t.Errorf("received an unexpected error: %v\n", err)
	}

	for key, expected := range expectedValues {
		if r[key] != expected {
			t.Errorf("Incorrect %s: expected %s, got %s\n", key, expected, r[key])
		}

	}
}
