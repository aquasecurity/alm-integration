package router

import "testing"

func TestAnonymizeSettings(t *testing.T) {
	tests := []struct {
		original *OutputSettings
		expected *OutputSettings
	}{{
		&OutputSettings{
			User: "admin",
		},
		&OutputSettings{
			User: "<hidden>",
		},
	}, {
		&OutputSettings{
			User: "",
		},
		&OutputSettings{
			User: "",
		},
	}, {
		&OutputSettings{
			Password: "secret",
		},
		&OutputSettings{
			Password: "<hidden>",
		},
	}, {
		&OutputSettings{
			Url: "http://localhost",
		},
		&OutputSettings{
			Url: "<hidden>",
		},
	},
	}

	for _, test := range tests {
		anonymized := anonymizeSettings(test.original)
		if anonymized == test.original {
			t.Errorf("Anonymized settings weren't cloned")
		}
		if anonymized.User != test.expected.User {
			t.Errorf("Settings anonymization is incorrect: expected User %s, got %s", test.expected.User, anonymized.User)
		}
		if anonymized.Password != test.expected.Password {
			t.Errorf("Settings anonymization is incorrect: expected Password %s, got %s", test.expected.Password, anonymized.Password)
		}
		if anonymized.Url != test.expected.Url {
			t.Errorf("Settings anonymization is incorrect: expected Url %s, got %s", test.expected.Url, anonymized.Url)
		}
	}
}
