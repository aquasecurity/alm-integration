package settings

type Settings struct {
	PluginName             string
	PolicyMinVulnerability string
	PolicyRegistry         []string
	PolicyImageName        []string
	PolicyNonCompliant     bool
	IgnoreRegistry         []string
	IgnoreImageName        []string

	PolicyOPA []string

	AggregateIssuesNumber   int
	AggregateTimeoutSeconds int
	IsScheduleRun           chan struct{}
	PolicyOnlyFixAvailable  bool
	PolicyShowAll           bool
	AquaServer              string
}

func GetDefaultSettings() *Settings {
	return &Settings{
		PluginName:             "",
		PolicyMinVulnerability: "",
		PolicyRegistry:         []string{},
		PolicyImageName:        []string{},
		PolicyShowAll:          false,
		PolicyNonCompliant:     false,
		IgnoreRegistry:         []string{},
		IgnoreImageName:        []string{},
	}
}
