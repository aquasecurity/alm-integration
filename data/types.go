package data

import "fmt"

type ScanImageInfo struct {
	Image                  string `json:"image"`
	Registry               string `json:"registry"`
	Digest                 string `json:"digest"`
	PreviousDigest         string `json:"previous_digest"`
	ImageAssuranceResults  `json:"image_assurance_results,omitempty"`
	VulnerabilitySummary   `json:"vulnerability_summary,omitempty"`
	ScanOptions            `json:"scan_options,omitempty"`
	Resources              []InfoResources `json:"resources,omitempty"`
	ApplicationScopeOwners []string        `json:"application_scope_owners,omitempty"`
	Malwares               []MalwareData   `json:"malware,omitempty"`
	SensitiveData          []SensitiveData `json:"sensitive_data,omitempty"`
}

type SensitiveData struct {
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Type     string `json:"type"`
	Hash     string `json:"hash"`
}

type MalwareData struct {
	Malware string `json:"malware"`
	Path    string `json:"path"`
	Hash    string `json:"hash"`
}

type ImageAssuranceResults struct {
	Disallowed      bool           `json:"disallowed"`
	ChecksPerformed []ControlCheck `json:"checks_performed"`
}

type ControlCheck struct {
	Control    string `json:"control"`
	PolicyName string `json:"policy_name"`
	Failed     bool   `json:"failed"`
}

type ScanOptions struct {
	ScanSensitiveData bool `json:"scan_sensitive_data"`
	ScanMalware       bool `json:"scan_malware"`
}

type VulnerabilitySummary struct {
	Total      int `json:"total"`
	Critical   int `json:"critical"`
	High       int `json:"high"`
	Medium     int `json:"medium"`
	Low        int `json:"low"`
	Negligible int `json:"negligible"`
	Sensitive  int `json:"sensitive"`
	Malware    int `json:"malware"`
}

type InfoResources struct {
	Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	ResourceDetails `json:"resource"`
}
type ResourceDetails struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Vulnerability struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	FixVersion string `json:"fix_version"`
	Severity   string `json:"aqua_severity"` //`json:""`nvd_severity
}

func BuildUniqueId(digest, image, registry string) string {
	return fmt.Sprintf("%s-%s-%s", digest, image, registry)
}

func (si *ScanImageInfo) GetUniqueId() string {
	return BuildUniqueId(si.Digest, si.Image, si.Registry)
}

func (si *ScanImageInfo) HasId() bool {
	return si.Digest != "" || si.Image != "" || si.Registry != ""
}
