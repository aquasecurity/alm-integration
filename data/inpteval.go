package data

type Inpteval interface {
	Eval(in map[string]interface{}, serverUrl string) (map[string]string, error)
	BuildAggregatedContent(items []map[string]string) (map[string]string, error)
	IsAggregationSupported() bool
}
