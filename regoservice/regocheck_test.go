package regoservice

import (
	"encoding/json"
	"testing"
)

func TestOpaRego(t *testing.T) {
	rego := `contains(input.image, "alpine")`
	incorrectRego := `default input = false`
	emptyRego := ""

	scanResult := `{"image":"alpine:26"}`
	scanNoJson := "simple text"
	scanWithoutResult := `{"image":"1science:latest"}`

	tests := []struct {
		rules              string
		scan               string
		result             bool
		shouldTriggerError bool
	}{
		{rego, scanResult, true, false},
		{rego, scanNoJson, false, true},
		{rego, scanWithoutResult, false, false},
		{incorrectRego, scanResult, false, true},
		{emptyRego, scanResult, true, false},
	}

	for _, test := range tests {
		intr := map[string]interface{}{}
		if err := json.Unmarshal([]byte(test.scan), &intr); err != nil && !test.shouldTriggerError {
			t.Errorf("json.Unmarshal(%q) error: %v", test.scan, err)
			continue
		}

		got, err := DoesMatchRegoCriteria(intr, test.rules)
		if err != nil && !test.shouldTriggerError {
			t.Errorf("received an unexpected error: %v", err)
			continue
		}
		if got != test.result {
			t.Errorf("isRegoCorrect(%q, %q) == %t, wanted %t", test.scan, test.rules, got, test.result)
		}
	}
}
