package webserver

import (
	"net/http"

	"github.com/aquasecurity/postee/router"
)

func (web *WebServer) reload(w http.ResponseWriter, r *http.Request) {
	router.Instance().ReloadConfig()
}
