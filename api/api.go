package api

import (
	"encoding/json"
	"github.com/Phill93/DoorManager/log"
	version2 "github.com/Phill93/DoorManager/version"
	"github.com/gorilla/mux"
	"net/http"
)

type Version struct {
	BuildDate string
	GitCommit string
	Version   string
	GoVersion string
	OsArch    string
}

type Response struct {
	Success bool
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	log.Infof("got version request from %s", r.Host)
	version := Version{
		BuildDate: version2.BuildDate,
		GitCommit: version2.GitCommit,
		Version:   version2.Version,
		GoVersion: version2.GoVersion,
		OsArch:    version2.OsArch,
	}

	js, err := json.Marshal(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handleBeep(w http.ResponseWriter, r *http.Request) {
	log.Infof("got beep request from %s", r.Host)
	response := Response{
		Success: true,
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func handleLed(w http.ResponseWriter, r *http.Request) {
	log.Infof("got led request from %s", r.Host)
	response := Response{
		Success: true,
	}

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func ServeAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/version", handleVersion)
	r.HandleFunc("/beep", handleBeep)
	r.HandleFunc("/led", handleLed)
	log.Fatal(http.ListenAndServe(":8080", r))
}
