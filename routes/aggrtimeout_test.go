package routes

import "testing"

func TestAggrTimeout(t *testing.T) {
	tests := []struct {
		caseDesc               string
		aggregateIssuesTimeout string
		expctdSeconds          int
	}{
		{
			"One minute",
			"1m",
			60,
		},
		{
			"Six hundredths seconds",
			"600s",
			600,
		},
		{
			"Two hours",
			"2h",
			7200,
		},
		{
			"Exact number of seconds",
			"300",
			300,
		},
		{
			"Invalid format",
			"xxxl",
			0,
		},
	}
	for _, test := range tests {
		route := &InputRoute{}
		route.Plugins.AggregateIssuesTimeout = test.aggregateIssuesTimeout
		route = ConfigureAggrTimeout(route)
		if route.Plugins.AggregateTimeoutSeconds != test.expctdSeconds {
			t.Errorf("[%s] Invalid number of seconds, expected %d, got %d \n", test.caseDesc, test.expctdSeconds, route.Plugins.AggregateTimeoutSeconds)
		}
	}

}
