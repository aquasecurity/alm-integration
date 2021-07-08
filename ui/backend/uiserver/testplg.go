package uiserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aquasecurity/postee/layout"
	"github.com/aquasecurity/postee/router"
)

func (srv *uiServer) testSettings(w http.ResponseWriter, r *http.Request) {
	plgSettings := &router.OutputSettings{}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal(body, plgSettings); err != nil {
		http.Error(w, fmt.Sprintf("Can't read JSON string %s", err), http.StatusBadRequest)
		return
	}

	plg := router.BuildAndInitOtpt(plgSettings, "")

	testPayload := make(map[string]string)

	testPayload["title"] = "Postee test title"
	testPayload["description"] = layout.GenTestDescription(plg.GetLayoutProvider(), "Postee test description")

	log.Printf("description is: %s \n", testPayload["description"])

	err = plg.Send(testPayload)

	if err != nil {
		//TODO provide method to write error response as JSON
		http.Error(w, fmt.Sprintf("Can't test output: %s \n", err), http.StatusBadRequest)
		return
	}

}
