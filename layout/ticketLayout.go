package layout

import (
	"bytes"
	"strconv"

	"github.com/aquasecurity/postee/data"
)

func GenTestDescription(provider LayoutProvider, raw string) string {
	var builder bytes.Buffer
	builder.WriteString(provider.P(raw))

	return builder.String()
}
func GenTicketDescription(provider LayoutProvider, scanInfo, prevScan *data.ScanImageInfo, link string) string {
	var builder bytes.Buffer
	builder.WriteString(provider.P("Image name: " + scanInfo.Image))
	builder.WriteString(provider.P("Registry: " + scanInfo.Registry))
	if scanInfo.Disallowed {
		builder.WriteString(provider.P("Image is non-compliant"))
	} else {
		builder.WriteString(provider.P("Image is compliant"))
	}

	if scanInfo.ScanMalware {
		if scanInfo.Malware > 0 {
			builder.WriteString(provider.P("Malware found: Yes"))
		} else {
			builder.WriteString(provider.P("Malware found: No"))
		}
	}

	if scanInfo.ScanSensitiveData {
		if scanInfo.Sensitive > 0 {
			builder.WriteString(provider.P("Sensitive data found: Yes"))
		} else {
			builder.WriteString(provider.P("Sensitive data found: No"))
		}
	}

	builder.WriteString(VulnerabilitiesTable(provider, [2][]string{
		{"CRITICAL", "HIGH", "MEDIUM", "LOW", "NEGLIGIBLE"},
		{strconv.Itoa(scanInfo.Critical), strconv.Itoa(scanInfo.High), strconv.Itoa(scanInfo.Medium), strconv.Itoa(scanInfo.Low), strconv.Itoa(scanInfo.Negligible)},
	}))

	// Rendering Assurances
	if len(scanInfo.ImageAssuranceResults.ChecksPerformed) > 0 {
		builder.WriteString(provider.TitleH2("Assurance controls"))
		builder.WriteString(RenderAssurances(provider, scanInfo.ImageAssuranceResults))
	}

	// Rendering Found vulnerabilities
	if len(scanInfo.Resources) > 0 {
		builder.WriteString(provider.TitleH2("Found vulnerabilities"))
		RenderVulnerabilities(scanInfo.Resources, provider, &builder)
	}

	// Discovered vulnerabilities from last scan:
	if prevScan != nil && len(prevScan.Resources) > 0 {
		builder.WriteString("\n")
		builder.WriteString(provider.TitleH2("Discovered vulnerabilities from last scan"))
		RenderVulnerabilities(prevScan.Resources, provider, &builder)
	}
	if len(scanInfo.Malwares) > 0 {
		builder.WriteString("\n")
		builder.WriteString(provider.TitleH2("Malware"))
		RenderMalware(scanInfo.Malwares, provider, &builder)
	}
	if len(scanInfo.SensitiveData) > 0 {
		builder.WriteString("\n")
		builder.WriteString(provider.TitleH2("Sensitive Data"))
		RenderSensitiveData(scanInfo.SensitiveData, provider, &builder)
	}
	builder.WriteString(provider.P("See more: " + provider.A(link, link)))
	return builder.String()
}
