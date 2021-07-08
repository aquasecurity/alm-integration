package uiserver

import (
	"net/http"
	"os"
	"path/filepath"
)

type localWebServer struct {
	localPath string
	url       string
}

func (web *localWebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	path = filepath.Join(web.localPath, path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(web.localPath, web.url))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.FileServer(http.Dir(web.localPath)).ServeHTTP(w, r)
}
