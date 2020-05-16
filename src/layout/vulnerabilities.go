package layout

import (
	"data"
	"strconv"
)

func RenderVulnerabilities(provider LayoutProvider, title string, vulns []data.Vulnerability) string {
	const empty = "none"
	var table [][]string
	table = append(table, []string{"#", "Name", "Version", "Fix version",})
	for i, v := range vulns {
		var name, version, fixVersion string
		if v.Name == "" {name = empty} else { name = v.Name}
		if v.Version == "" { version = empty} else {version =v.Version}
		if v.FixVersion == "" { fixVersion = empty} else { fixVersion = data.ClearField(v.FixVersion)}

		table = append(table, []string{
			strconv.Itoa(i+1),
			name,
			version,
			fixVersion,
		})
	}
	return provider.Table(table)
}

func VulnerabilitiesTable(provider LayoutProvider, rows [2][]string) string  {
	if len(rows) != 2 && len(rows[1]) != 5 {
		return ""
	}
	var table [][]string
	table = append(table, rows[0])
	var r []string
	r = append(r, provider.ColourText(rows[1][0], data.CriticalColor()))
	r = append(r, provider.ColourText(rows[1][1], data.HighColor()))
	r = append(r, provider.ColourText(rows[1][2], data.MediumColor()))
	r = append(r, provider.ColourText(rows[1][3], data.LowColor()))
	r = append(r, provider.ColourText(rows[1][4], data.NegligibleColor()))
	table = append(table, r)
	return provider.Table(table)
}


