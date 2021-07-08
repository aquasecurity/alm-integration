package routes

type InputRoute struct {
	Name       string   `json:"name"`
	Input      string   `json:"input"`
	Outputs    []string `json:"outputs"`
	Plugins    Plugins  `json:"plugins"`
	Template   string   `json:"template"`
	Scheduling chan struct{}
}

type Plugins struct {
	AggregateIssuesNumber   int    `json:"aggregate-issues-number"`
	AggregateIssuesTimeout  string `json:"aggregate-issues-timeout"`
	PolicyShowAll           bool   `json:"policy-show-all"`
	AggregateTimeoutSeconds int
}

func (route *InputRoute) IsSchedulerRun() bool {
	return route.Scheduling != nil
}
func (route *InputRoute) StartScheduler() {
	route.Scheduling = make(chan struct{})
}

func (route *InputRoute) StopScheduler() {
	if route.Scheduling != nil {
		close(route.Scheduling)
	}
}
