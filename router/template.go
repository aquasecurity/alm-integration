package router

type Template struct {
	Name               string `json:"name"`
	Body               string `json:"body"`
	RegoPackage        string `json:"rego-package"`
	LegacyScanRenderer string `json:"legacy-scan-renderer"`
	Url                string `json:"url"`
}
