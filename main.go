package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nohlachilders/bootdevserver/internal/database"
)

//"fmt"

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Environment variable DB_URL must be set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	dbQueries := database.New(db)

	servemux := http.ServeMux{}
	server := http.Server{
		Handler: &servemux,
		Addr:    ":8080",
	}
	cfg := apiConfig{
		db: dbQueries,
	}

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
	db             *database.Queries
}
