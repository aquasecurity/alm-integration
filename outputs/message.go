package outputs

import (
	"bytes"
	"strings"

	"github.com/aquasecurity/postee/layout"
)

func buildShortMessage(server, urls string, provider layout.LayoutProvider) string {
	var builder bytes.Buffer
	builder.WriteString(provider.P("This message is too long to display here. Please visit the link to read the content."))
	links := strings.Split(urls, "\n")
	for _, link := range links {
		builder.WriteString(provider.P(provider.A(server+link, link)))
	}
	return builder.String()
}
