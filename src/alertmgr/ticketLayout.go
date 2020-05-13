package alertmgr

import (
	"bytes"
	"encoding/json"
)

func GenTicketDescription(scanInfo *ScanImageInfo) string {
	var builder bytes.Buffer
	builder.WriteString("Image name: " + scanInfo.Image + "\n")
	builder.WriteString("Registry: " + scanInfo.Registry )
	if scanInfo.Disallowed {
		builder.WriteString("Image is non-compliant\n")
	} else {
		builder.WriteString("Image is compliance\n")
	}
	builder.WriteString(
		RenderVulnerabilitiesCounts(
			scanInfo.Critical, scanInfo.Negligible, scanInfo.Medium, scanInfo.Low, scanInfo.Negligible ))

	if scanInfo.ScanMalware {
		if scanInfo.Malware > 0 {
			builder.WriteString("Malware found: YES\n")
		} else {
			builder.WriteString("Malware found: no\n")
		}
	}

	if scanInfo.ScanSensitiveData {
		if scanInfo.Sensitive > 0 {
			builder.WriteString("Sensitive data found: yes\n")
		} else {
			builder.WriteString("Sensitive data found: no\n")
		}
	}

	// Rendering Assurances
	builder.WriteString( RenderAssurances(scanInfo.ImageAssuranceResults))

	// Rendering Found vulnerabilities
	builder.WriteString( "\nh2. Found vulnerabilities\n")
	for _, r := range scanInfo.Resources {
		v := RenderVulnerabilities( r.Name, r.Vulnerabilities)
		builder.WriteString( v )
	}

	builder.WriteString("" + "\n")
	return builder.String()
}

func GenTicketDescriptionSimpleAdd(scanInfo *ScanImageInfo) string {
	builder := "Image name: " + scanInfo.Image + "\n"
	builder += "Registry: " + "\n"
	builder += "" + "\n"
	return builder
}

func ParseImageInfo(source []byte) (*ScanImageInfo, error) {
	scanInfo := new(ScanImageInfo)
	err := json.Unmarshal(source, scanInfo)
	if err != nil {
		return nil, err
	}
	return scanInfo, nil
}
