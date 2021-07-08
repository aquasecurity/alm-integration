package routes

import (
	"log"
	"strconv"
	"strings"
)

func ConfigureAggrTimeout(route *InputRoute) *InputRoute {
	var timeout int
	var err error

	times := map[string]int{
		"s": 1,
		"m": 60,
		"h": 3600,
	}

	if len(route.Plugins.AggregateIssuesTimeout) > 0 {
		wasConvert := false
		for suffix, k := range times {
			if strings.HasSuffix(strings.ToLower(route.Plugins.AggregateIssuesTimeout), suffix) {
				timeout, err = strconv.Atoi(strings.TrimSuffix(route.Plugins.AggregateIssuesTimeout, suffix))
				timeout *= k
				wasConvert = true
				break
			}
		}
		if !wasConvert {
			timeout, err = strconv.Atoi(route.Plugins.AggregateIssuesTimeout)
		}
		if err != nil {
			log.Printf("%q settings: Can't convert 'AggregateIssuesTimeout'(%q) to seconds.",
				route.Name, route.Plugins.AggregateIssuesTimeout)
		}
	}
	route.Plugins.AggregateTimeoutSeconds = timeout
	return route
}
