package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nohlachilders/bootdevserver/internal/auth"
	"github.com/nohlachilders/bootdevserver/internal/database"
)

type Chirp struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Body      string        `json:"body"`
	UserID    uuid.NullUUID `json:"user_id"`
}

func (cfg *apiConfig) getAllChirpsHandler(w http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetAllChirps(req.Context())
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	thisChirps := make([]Chirp, len(chirps))
	for i, chirp := range chirps {
		thisChirps[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
	}
	respondWithJSON(w, 200, thisChirps)
}

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, req *http.Request) {
	thisId, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		respondWithError(w, 400, "Invalid ID")
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), thisId)
	if err != nil {
		respondWithError(w, 404, "Not Found")
		return
	}
	respondWithJSON(w, 200, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, req *http.Request) {
	type validationRequest struct {
		Body string `json:"body"`
	}
	thisRequest := validationRequest{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&thisRequest)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	if len(thisRequest.Body) >= 140 {
		respondWithError(w, 400, "Chirpy too long")
		return
	}

	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}
	parsedUUID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("%v", parsedUUID))
		return
	}

	packedNullUUID := uuid.NullUUID{
		UUID:  parsedUUID,
		Valid: true,
	}
	cleanedBody, _ := cleanBody(thisRequest.Body)
	params := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: packedNullUUID,
	}
	chirp, err := cfg.db.CreateChirp(req.Context(), params)
	respondWithJSON(w, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func cleanBody(body string) (cleaned_body string, was_cleaned bool) {
	wordBlacklist := []string{"Kerfuffle", "Sharbert", "Fornax"}

	split := strings.Split(body, " ")
	for i, word := range split {
		for _, badWord := range wordBlacklist {
			if word == badWord {
				split[i] = "****"
			}
		}
		for _, badWord := range wordBlacklist {
			if word == strings.ToLower(badWord) {
				split[i] = "****"
			}
		}
	}
	cleaned_body = strings.Join(split, " ")
	was_cleaned = !(cleaned_body == body)

	return cleaned_body, was_cleaned
}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, req *http.Request) {
	thisId, err := uuid.Parse(req.PathValue("id"))
	if err != nil {
		respondWithError(w, 400, "Invalid ID")
		return
	}
	chirp, err := cfg.db.GetChirp(req.Context(), thisId)
	if err != nil {
		respondWithError(w, 404, "Not Found")
		return
	}
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	parsedUUID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("%v", parsedUUID))
		return
	}
	if chirp.UserID.UUID != parsedUUID {
		w.WriteHeader(403)
		return
	}
	err = cfg.db.DeleteChirp(req.Context(), thisId)
	if err != nil {
		respondWithError(w, 404, "Not Found")
		return
	}

	w.WriteHeader(204)
}
