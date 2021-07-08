package uiserver

import (
	"encoding/json"
	"net/http"

	"github.com/aquasecurity/postee/ui/backend/dbservice"
)

func (srv *uiServer) plgnStats(w http.ResponseWriter, r *http.Request) {
	stats, err := dbservice.GetPlgnStats()
	if err != nil {
		handleErr(w, err)
		return
	}
	data, err := json.Marshal(stats)
	if err != nil {
		handleErr(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(data))
}
func handleErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
