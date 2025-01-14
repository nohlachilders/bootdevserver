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
	platform := os.Getenv("CHIRPY_PLATFORM")
	if dbURL == "" {
		log.Fatal("Environment variable CHIRPY_PLATFORM must be set")
	}

	servemux := http.ServeMux{}
	server := http.Server{
		Handler: &servemux,
		Addr:    ":8080",
	}
	cfg := apiConfig{
		db:       dbQueries,
		platform: platform,
	}

	filesystemHandler := http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir("."))))
	servemux.Handle("/app/", filesystemHandler)

	servemux.HandleFunc("POST /api/users", cfg.userCreationHandler)
	servemux.HandleFunc("GET /api/healthz", healthResponseHandler)
	servemux.HandleFunc("POST /api/validate_chirp", validationResponseHandler)

	servemux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	servemux.HandleFunc("GET /admin/metrics", cfg.metricsResponseHandler)

	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	platform       string
	fileserverHits atomic.Int32
	db             *database.Queries
}
