package api

import (
	"encoding/json"
	version2 "github.com/Phill93/DoorManager/version"
	"net/http"
)

type Version struct {
	BuildDate string
	GitCommit string
	Version   string
	GoVersion string
	OsArch    string
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
