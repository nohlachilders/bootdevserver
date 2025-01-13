package main

import (
	"fmt"
	"net/http"
)

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
