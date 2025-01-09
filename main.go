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
		Addr:    ":8080",
	}
	cfg := apiConfig{}

	filesystemHandler := http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir("."))))
	servemux.Handle("/app/", filesystemHandler)

	servemux.HandleFunc("GET /api/healthz", healthResponseHandler)
	servemux.HandleFunc("POST /admin/reset", cfg.metricsResetHandler)
	servemux.HandleFunc("GET /admin/metrics", cfg.metricsResponseHandler)

	log.Fatal(server.ListenAndServe())
}

func healthResponseHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsResponseHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	hits := cfg.fileserverHits.Load()
	str := fmt.Sprintf(`
        <html>
            <body>
                <h1>Welcome, Chirpy Admin</h1>
                <p>Chirpy has been visited %d times!</p>
            </body>
        </html>
        `, hits)
	w.Write([]byte(str))
}

func (cfg *apiConfig) metricsResetHandler(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	//w.Write([]byte(str))
}
