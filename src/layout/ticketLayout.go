package layout

import (
	"bytes"
	"data"
	"strconv"
)

func GenTicketDescription(provider LayoutProvider, scanInfo, prevScan *data.ScanImageInfo) string {
	var builder bytes.Buffer
	builder.WriteString(provider.P("Image name: " + scanInfo.Image))
	builder.WriteString(provider.P("Registry: " + scanInfo.Registry  ))
	if scanInfo.Disallowed {
		builder.WriteString(provider.P("Image is non-compliant"))
	} else {
		builder.WriteString(provider.P("Image is compliant"))
	}

	if scanInfo.ScanMalware {
		if scanInfo.Malware > 0 {
			builder.WriteString( provider.P("Malware found: Yes"))
		} else {
			builder.WriteString(provider.P("Malware found: No"))
		}
	}

	if scanInfo.ScanSensitiveData {
		if scanInfo.Sensitive > 0 {
			builder.WriteString(provider.P("Sensitive data found: Yes"))
		} else {
			builder.WriteString( provider.P("Sensitive data found: No"))
		}
	}

	builder.WriteString( VulnerabilitiesTable(provider, [2][]string {
		{"CRITICAL","HIGH","MEDIUM","LOW","NEGLIGIBLE"},
		{strconv.Itoa(scanInfo.Critical), strconv.Itoa(scanInfo.High), strconv.Itoa(scanInfo.Medium), strconv.Itoa(scanInfo.Low), strconv.Itoa(scanInfo.Negligible)},
	}))

	// Rendering Assurances
	if len(scanInfo.ImageAssuranceResults.ChecksPerformed) > 0 {
		builder.WriteString( provider.TitleH2("Assurance controls"))
		builder.WriteString( RenderAssurances(provider, scanInfo.ImageAssuranceResults))
	}

	// Rendering Found vulnerabilities
	if len(scanInfo.Resources) > 0 {
		builder.WriteString( provider.TitleH2("Found vulnerabilities"))
		RenderVulnerabilities( scanInfo.Resources, provider, &builder)
	}

	// Discovered vulnerabilities from last scan:
	if prevScan != nil && len(prevScan.Resources) > 0 {
		builder.WriteString("\n")
		builder.WriteString( provider.TitleH2("Discovered vulnerabilities from last scan" ))
		RenderVulnerabilities(prevScan.Resources, provider, &builder)
	}
	return builder.String()
}



