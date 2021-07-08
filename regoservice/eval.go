package regoservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aquasecurity/postee/data"
	"github.com/open-policy-agent/opa/rego"
)

const (
	result_prop          = "result"
	title_prop           = "title"
	aggregation_pkg_prop = "aggregation_pkg"
)

var (
	buildinRegoTemplates = []string{"./rego-templates"}
	commonRegoTemplates  = []string{"./rego-templates/common"}
)

type regoEvaluator struct {
	prepQuery        *rego.PreparedEvalQuery
	aggrQuery        *rego.PreparedEvalQuery
	isPackageDefined bool
}

func (regoEvaluator *regoEvaluator) IsAggregationSupported() bool {
	return regoEvaluator.aggrQuery != nil
}

func (regoEvaluator *regoEvaluator) Eval(in map[string]interface{}, serverUrl string) (map[string]string, error) {
	ctx := context.Background()
	rs, err := regoEvaluator.prepQuery.Eval(ctx, rego.EvalInput(in))

	if err != nil {
		return nil, err
	}

	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return nil, errors.New("no results") //TODO error definition
	}

	var expr interface{}
	if regoEvaluator.isPackageDefined {
		expr = rs[0].Expressions[0].Value
	} else {
		expr = getFirstElement(rs[0].Expressions[0].Value.(map[string]interface{}), result_prop)
		if expr == nil {
			return nil, errors.New("invalid rego template structure")
		}
	}

	data := expr.(map[string]interface{})

	title, err := asStringOrJson(data, title_prop)
	if err != nil {
		return nil, err
	}

	description, err := asStringOrJson(data, result_prop)

	if err != nil {
		return nil, err
	}

	return map[string]string{
		"title":       title,
		"description": description,
		"url":         serverUrl,
	}, nil

}

func getFirstElement(context map[string]interface{}, key string) interface{} {
	for _, v := range context {
		log.Printf("checking: %s ...\n", key)
		childCtx, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}
		if childCtx[key] != nil {
			return v
		} else {
			found := getFirstElement(childCtx, key)
			if found != nil {
				return found
			}
		}
	}
	return nil
}

func asStringOrJson(data map[string]interface{}, prop string) (string, error) {
	expr, ok := data[prop]
	if !ok {
		return "", errors.New(fmt.Sprintf("property %s is not found", prop))
	}
	fmt.Printf("value: %q", expr)
	switch v := expr.(type) {
	case string:
		return v, nil
	default:
		val, err := json.Marshal(expr)
		if err != nil {
			return "", err
		}
		return string(val), nil
	}
}
func (regoEvaluator *regoEvaluator) BuildAggregatedContent(scans []map[string]string) (map[string]string, error) {
	aggregatedJson := make([]map[string]interface{}, len(scans), len(scans))

	for _, scan := range scans {
		desc := scan["description"]
		var in []map[string]interface{}

		item := make(map[string]interface{})

		if err := json.Unmarshal([]byte(desc), &in); err != nil {
			item["description"] = desc //description is not json, so it's passed as string
		} else {
			item["description"] = in
		}

		item["title"] = scan["title"]

		aggregatedJson = append(aggregatedJson, item)
	}

	ctx := context.Background()
	rs, err := regoEvaluator.aggrQuery.Eval(ctx, rego.EvalInput(aggregatedJson))

	if err != nil {
		return nil, err
	}

	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return nil, errors.New("no results") //TODO error definition
	}

	expr := rs[0].Expressions[0].Value

	data := expr.(map[string]interface{})

	title, err := asStringOrJson(data, title_prop)

	if err != nil {
		return nil, err
	}

	description, err := asStringOrJson(data, result_prop)

	if err != nil {
		return nil, err
	}

	return map[string]string{
		"title":       title,
		"description": description,
	}, nil
}

func BuildBundledRegoEvaluator(rego_package string) (data.Inpteval, error) {
	r, err := buildBundledRegoForPackage(rego_package)

	if err != nil {
		return nil, err
	}

	aggrQuery, err := buildAggregatedRego(r)

	if err != nil {
		return nil, err
	}

	return &regoEvaluator{
		prepQuery:        r,
		isPackageDefined: true,
		aggrQuery:        aggrQuery,
	}, nil
}
func buildBundledRegoForPackage(rego_package string) (*rego.PreparedEvalQuery, error) {
	ctx := context.Background()
	query := fmt.Sprintf("data.%s", rego_package)

	r, err := rego.New(
		rego.Query(query),
		rego.Load(buildinRegoTemplates, nil),
	).PrepareForEval(ctx)

	if err != nil {
		return nil, err
	}

	return &r, nil
}
func buildAggregatedRego(query *rego.PreparedEvalQuery) (*rego.PreparedEvalQuery, error) {
	ctx := context.Background()

	//execute query with empty input and check if aggregation package is defined
	rs, err := query.Eval(ctx, rego.EvalInput(make(map[string]interface{})))

	if len(rs) == 0 || len(rs[0].Expressions) == 0 {
		return nil, errors.New("no results") //TODO error definition
	}

	expr := rs[0].Expressions[0].Value.(map[string]interface{})

	aggregation_pkg_val := expr[aggregation_pkg_prop]

	var aggrQuery *rego.PreparedEvalQuery

	if aggregation_pkg_val != nil {
		aggregation_pkg := aggregation_pkg_val.(string)
		aggrQuery, err = buildBundledRegoForPackage(aggregation_pkg)
		if err != nil {
			return nil, err
		}
	} else {
		//it's ok skip aggregation package - no aggregation features will be available
		log.Printf("No aggregation package configured!!!")
	}
	return aggrQuery, nil
}

func BuildExternalRegoEvaluator(filename string, body string) (data.Inpteval, error) {
	ctx := context.Background()

	r, err := rego.New(
		rego.Query("data"),
		rego.Load(commonRegoTemplates, nil), //only common modules
		rego.Module(filename, body),
	).PrepareForEval(ctx)

	if err != nil {
		return nil, err
	}

	aggrQuery, err := buildAggregatedRego(&r)

	if err != nil {
		return nil, err
	}

	return &regoEvaluator{
		prepQuery:        &r,
		isPackageDefined: false,
		aggrQuery:        aggrQuery,
	}, nil
}
