package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/nohlachilders/bootdevserver/internal/database"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Environment variable DB_URL must be set")
	}
	secret := os.Getenv("CHIRPY_SECRET")
	if dbURL == "" {
		log.Fatal("Environment variable CHIRPY_SECRET must be set")
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

	const port string = ":8080"
	const fileSystemRoot = "."

	servemux := http.ServeMux{}
	server := http.Server{
		Handler: &servemux,
		Addr:    port,
	}
	cfg := apiConfig{
		db:       dbQueries,
		platform: platform,
		secret:   secret,
	}
	filesystemHandler := http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(http.Dir(fileSystemRoot))))
	servemux.Handle("/app/", filesystemHandler)

	servemux.HandleFunc("POST /api/users", cfg.userCreationHandler)
	servemux.HandleFunc("PUT /api/users", cfg.userUpdateHandler)
	servemux.HandleFunc("POST /api/login", cfg.userLoginHandler)
	servemux.HandleFunc("POST /api/refresh", cfg.userRefreshHandler)
	servemux.HandleFunc("POST /api/revoke", cfg.userRevokeHandler)
	servemux.HandleFunc("GET /api/healthz", healthResponseHandler)
	servemux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
	servemux.HandleFunc("GET /api/chirps", cfg.getAllChirpsHandler)
	servemux.HandleFunc("GET /api/chirps/{id}", cfg.getChirpHandler)
	servemux.HandleFunc("DELETE /api/chirps/{id}", cfg.deleteChirpHandler)

	servemux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	servemux.HandleFunc("GET /admin/metrics", cfg.metricsResponseHandler)

	fmt.Printf("Serving on port %s...\n", port)
	log.Fatal(server.ListenAndServe())
}

type apiConfig struct {
	platform       string
	fileserverHits atomic.Int32
	db             *database.Queries
	secret         string
}
