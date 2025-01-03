package main

import (
	//"fmt"
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
    servemux := http.ServeMux{}
    server := http.Server{
        Handler: &servemux,
        Addr: ":8080",
    }
    cfg := apiConfig{
    }

    servemux.Handle("/", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(".")))))
    servemux.HandleFunc("/healthz", healthResponseHandler)
    servemux.HandleFunc("/metrics", cfg.middlewareMetricsHandler)

    log.Fatal(server.ListenAndServe())
}

func healthResponseHandler(w http.ResponseWriter, req *http.Request) {
    w.Header().Add("Content-Type","text/plain; charset=utf-8")
    w.WriteHeader(200)
    w.Write([]byte("OK"))
}

type apiConfig struct {
    fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
    cfg.fileserverHits.Add(4)
    return next
}

func (cfg *apiConfig) middlewareMetricsHandler(w http.ResponseWriter, req *http.Request) {
    hits := cfg.fileserverHits.Load()
    str := fmt.Sprintf("Hits: %d", hits)
    w.Write([]byte(str))
}
