package uiserver

import (
	"net/http"
	"strings"
)

func (srv *uiServer) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !strings.HasPrefix(r.RequestURI, "/api") {
			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(r.RequestURI, "/api/login") {
			next.ServeHTTP(w, r)
			return
		}

		if user, err := srv.getUserFromRequest(r); err == nil && user != "" {
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			http.Error(w, "Forbidden", http.StatusUnauthorized)
		}
	})
}

func (srv *uiServer) getUserFromRequest(r *http.Request) (string, error) {
	session, err := srv.store.Get(r, sessioncookiename)
	if err != nil {
		return "", err
	}
	userObj := session.Values["user"]
	if userObj == nil {
		return "", nil
	} else {
		return userObj.(string), nil
	}
}
