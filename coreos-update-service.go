package main

import (
	coreos "./coreos"
	"encoding/xml"
	"github.com/Sirupsen/logrus"
	omaha "github.com/coreos/go-omaha/omaha"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"encoding/json"
)

var updater coreos.CoreOsUpdater


func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var request omaha.Request

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	logrus.Infof("New update request: %s ", string(body))

	xml.Unmarshal(body, &request)
	if len(request.Apps) != 1 {
		logrus.Error("More that one app update in the request not suported")
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Support only CoreOs
	app := request.Apps[0]
	if app.Id == coreos.CoreOs {
		resp_body, err := xml.Marshal(updater.CoreOsUpdate(request))
		if err != nil {
			panic(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(resp_body)
	} else {
		logrus.Errorf("Application %s is not suported", app.Id)
		w.WriteHeader(http.StatusConflict)
	}
}

func NodesHandler(w http.ResponseWriter, r *http.Request){
	nodes := updater.GetNodes()
	data, _ := json.Marshal(nodes)
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func UpdateVersionHandler(w http.ResponseWriter, r *http.Request) {
	var version = coreos.CoreOSVersion{}
	vars := mux.Vars(r)
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &version)

	if err != nil {
		logrus.Error("%s: %s", err, string(body))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if vars["id"] != version.VersionId {
		logrus.Errorf("Inconsitency at the id %s and the metadata %s", vars["id"], version.VersionId)
		w.WriteHeader(http.StatusConflict)
		return
	}
	updater.UpdateVersion(version)
	w.WriteHeader(http.StatusCreated)
}

func UpdateTrackHandler(w http.ResponseWriter, r *http.Request) {
	var tracks map[string]string
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &tracks)
	if err != nil {
		logrus.Error("%s: %s", err, string(body))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	updater.UpdateTracks(tracks)
	w.WriteHeader(http.StatusCreated)
}

func GetTrackHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(updater.GetTracks())
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}


func main() {
	storage := coreos.NewLocalDB("update.db")
	updater = coreos.CoreOsUpdater{
		Storage: storage,
	}

	r := mux.NewRouter()
	r.HandleFunc("/v1/update/", UpdateHandler)

	r.Methods(http.MethodGet).Path("/nodes").HandlerFunc(NodesHandler)
	r.Methods(http.MethodPut).Path("/version/{id}").HandlerFunc(UpdateVersionHandler)
	r.Methods(http.MethodPut).Path("/tracks").HandlerFunc(UpdateTrackHandler)
	r.Methods(http.MethodGet).Path("/tracks").HandlerFunc(GetTrackHandler)


	http.ListenAndServe(":8000", handlers.LoggingHandler(os.Stdout, r))
}
