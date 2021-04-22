package uiserver

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type uiServer struct {
	port       string
	cfgPath    string
	boltDbPath string
	updateUrl  string
	updateKey  string
	admusr     string
	admpwd     string
	router     *mux.Router
	store      *sessions.CookieStore
}

func Instance(webLocalPath, port, cfg, updateUrl, updateKey string, admusr string, admpwd string) *uiServer {
	server := &uiServer{
		port:      port,
		cfgPath:   cfg,
		updateUrl: updateUrl,
		updateKey: updateKey,
		admusr:    admusr,
		admpwd:    admpwd,
		router:    mux.NewRouter().StrictSlash(true),
	}
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	server.store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	server.store.Options = &sessions.Options{
		MaxAge:   60 * 60 * 24, //one day
		HttpOnly: true,
	}

	server.router.Use(server.authenticationMiddleware)

	server.router.HandleFunc("/api/login", server.login).Methods("POST")
	server.router.HandleFunc("/api/logout", server.logout).Methods("GET")
	server.router.HandleFunc("/api/update", server.updateConfig).Methods("POST")
	server.router.HandleFunc("/api/plugins", server.pluginList).Methods("GET")
	server.router.HandleFunc("/api/test", server.testSettings).Methods("POST")
	server.router.HandleFunc("/api/plugins/stats", server.plgnStats).Methods("GET")

	web := &localWebServer{
		localPath: webLocalPath,
		url:       "/",
	}
	server.router.PathPrefix("/").Handler(web)
	return server
}

func (srv *uiServer) Start() {
	log.Print("UI Postee server starting...")
	http.ListenAndServe(":"+srv.port, srv.router)
}

func (srv *uiServer) Stop() {
	log.Print("UI Postee server stopped!")
}