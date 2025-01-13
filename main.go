package main

import (
	//"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	servemux := http.ServeMux{}
	server := http.Server{
		Handler: &servemux,
		Addr:    ":8080",
	}
	cfg := apiConfig{}

	filesystemHandler := http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir("."))))
	servemux.Handle("/app/", filesystemHandler)

	servemux.HandleFunc("GET /api/healthz", healthResponseHandler)
	servemux.HandleFunc("POST /api/validate_chirp", validationResponseHandler)

	servemux.HandleFunc("POST /admin/reset", cfg.metricsResetHandler)
	servemux.HandleFunc("GET /admin/metrics", cfg.metricsResponseHandler)

	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	fileserverHits atomic.Int32
}
