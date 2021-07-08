package data

import "testing"

const (
	dist     = "dist"
	image    = "image"
	register = "reg"
	want     = "dist-image-reg"
)

func TestBuildUniqueId(t *testing.T) {
	if got := BuildUniqueId(dist, image, register); got != want {
		t.Errorf("BuildUniqueId(%q, %q, %q) == %q, want %q", dist, image, register, got, want)
	}
}

func TestGetUniqueId(t *testing.T) {
	si := &ScanImageInfo{
		Image:    image,
		Registry: register,
		Digest:   dist,
	}
	if got := si.GetUniqueId(); got != want {
		t.Errorf(".GetUniqueId() == %q, want %q", got, want)
	}
}
func TestHasUniqueId(t *testing.T) {
	tests := []struct {
		image          string
		registry       string
		digest         string
		expectedResult bool
	}{
		{"1111", "2222", "3333", true},
		{"1111", "", "3333", true},
		{"", "", "3333", true},
		{"", "", "", false},
	}

	for _, test := range tests {
		si := &ScanImageInfo{
			Image:    test.image,
			Registry: test.registry,
			Digest:   test.digest,
		}
		if got := si.HasId(); got != test.expectedResult {
			t.Errorf(`for scan (image: "%s",registry: "%s",digest: "%s") HasUniqueId()== %t, want %t`, test.image, test.registry, test.digest, got, test.expectedResult)
		}
	}
}
