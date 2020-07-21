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

func handleGate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		vars := mux.Vars(r)
		log.Debugf("Request to %d gate %d from %s", vars["action"], vars["id"], r.Host)
		response := Response{Success: true}
		js, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	default:
		http.Error(w, "Not Implemented", http.StatusNotImplemented)
	}
}

func ServeAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/version", handleVersion)
	r.HandleFunc("/gate/{id}/{action}", handleGate)
}
