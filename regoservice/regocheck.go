package regoservice

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

const module = `package postee

default allow = false

allow {
%s
}
`

func DoesMatchRegoCriteria(input interface{}, rule string) (bool, error) {

	if rule == "" {
		return true, nil //no rule defined - any input allowed
	}

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
