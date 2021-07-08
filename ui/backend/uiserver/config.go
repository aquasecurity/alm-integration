package uiserver

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	hookDbService "github.com/aquasecurity/postee/dbservice"
	"github.com/aquasecurity/postee/router"
)

func (srv *uiServer) getConfig(w http.ResponseWriter, r *http.Request) {
	log.Printf("configured config path %s", srv.cfgPath)

	_, err := router.Parsev2cfg(srv.cfgPath)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Invalid config file format: %s", err.Error())))
		return
	}

	d, err := ioutil.ReadFile(srv.cfgPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "text/yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(d)
}

func (srv *uiServer) updateConfig(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	inputYaml, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can't read JSON string", http.StatusBadRequest)
		return
	}
	if err := os.Rename(srv.cfgPath, srv.cfgPath+".copy"); err != nil {
		log.Printf("rename file error %v", err)
		http.Error(w, "Can't remove data from the config file for overwrite", http.StatusBadRequest)
		return
	}
	f, err := os.Create(srv.cfgPath)
	if err != nil {
		log.Printf("create file error %v", err)
		http.Error(w, "Can't open the config file for overwrite", http.StatusBadRequest)
		return
	}
	defer f.Close()
	_, err = f.Write(inputYaml)
	if err != nil {
		log.Printf("write file error %v", err)
		http.Error(w, "Can't write to the config file for overwrite", http.StatusBadRequest)
		return
	}
	os.RemoveAll(srv.cfgPath + ".copy")

	apikey, err := hookDbService.GetApiKey()

	if err != nil {
		log.Printf("Can not load api key from bolt %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = reloadWebhookCfg(srv.webhookUrl, apikey)

	if err != nil {
		log.Printf("Can not reload webhook server %v", err)
		http.Error(w, "Can not reload webhook server", http.StatusBadRequest)
		return
	}
}

func reloadWebhookCfg(url string, key string) error {
	u := fmt.Sprintf("%s/reload?key=%s", url, key)
	resp, err := http.Get(u)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
