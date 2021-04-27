package regoservice

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
)

const module = `package postee

default allow = false

allow {
%s
}
`

func IsRegoCorrectScanResult(rule string, scanResult string) (bool, error) {
	var input interface{}
	if err := json.Unmarshal([]byte(scanResult), &input); err != nil {
		return false, err
	}
	return IsRegoCorrectInterface(input, rule)
}

func IsRegoCorrectInterface(input interface{}, rule string) (bool, error) {
	r := rego.New(
		rego.Query("x = data.postee.allow"),
		rego.Module("postee.rego", fmt.Sprintf(module, rule)),
	)

	ctx := context.Background()
	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return false, err
	}

	rs, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, err
	}

	if len(rs) > 0 {
		switch rs[0].Bindings["x"].(type) {
		case bool:
			return rs[0].Bindings["x"].(bool), nil
		}
	}
	return false, nil
}