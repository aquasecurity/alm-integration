package webserver

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aquasecurity/postee/router"
	"github.com/aquasecurity/postee/utils"
	"github.com/gorilla/mux"
)

func (ctx *WebServer) tenantHandler(w http.ResponseWriter, r *http.Request) {
	route, ok := mux.Vars(r)["route"]
	if !ok || len(route) == 0 {
		log.Printf("Failed route: %q", route)
		ctx.writeResponse(w, http.StatusBadRequest, "failed route")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed ioutil.ReadAll: %s", err)
		ctx.writeResponseError(w, http.StatusInternalServerError, err)
		return
	}

	defer r.Body.Close()
	utils.Debug("%s\n\n", string(body))
	router.Instance().HandleRoute(route, body)
	ctx.writeResponse(w, http.StatusOK, "")
}
