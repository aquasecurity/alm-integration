package uiserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const (
	testCfgFile = "test.cfg"

	inputConfigJson = `[{"type":"common"}]`
)

func TestUpdateConfig(t *testing.T) {
	tests := []struct {
		input  string
		status int
	}{
		{inputConfigJson, http.StatusOK},
	}

	srv := &uiServer{
		cfgPath: testCfgFile,
	}
	os.Create(testCfgFile)
	defer os.RemoveAll(testCfgFile)

	for _, test := range tests {
		req := httptest.NewRequest("POST", "/update", strings.NewReader(test.input))
		w := httptest.NewRecorder()
		srv.updateConfig(w, req)
		response := w.Result()

		msg, err := io.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}

		if st := w.Result().StatusCode; st != test.status {
			t.Errorf("request to /update returns a wrong status %d, wanted %d.\nData: %q\nMessage: %q", st, test.status, test.input, string(msg))
		}
		response.Body.Close()
	}
}
