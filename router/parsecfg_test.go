package router

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseCfgWithInvalidFilename(t *testing.T) {
	invalidfn := "not-a-cfg.yaml"
	_, err := Parsev2cfg(invalidfn)

	if err == nil {
		t.Errorf("Error is expected")
	}
}

func TestParseCfgWithInvalidYaml(t *testing.T) {
	cfgfn := "cfg.yaml"
	invalidYaml := `
playing_song_artist: Playing song, {{ song_name }} by {{ artist }}

playing_playlist: {{ action }} playlist {{ playlist_name }}`
	defer func() {
		os.Remove(cfgfn)
	}()

	ioutil.WriteFile(cfgfn, []byte(invalidYaml), 0644)
	_, err := Parsev2cfg(cfgfn)

	if err == nil {
		t.Errorf("Error is expected")
	}

}
