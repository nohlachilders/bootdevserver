package main

import (
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed on dev platform"))
		return
	}
	cfg.fileserverHits.Store(0)
	cfg.db.Reset(req.Context())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics hits reset and db reset"))
}
