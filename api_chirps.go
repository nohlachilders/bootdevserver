package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nohlachilders/bootdevserver/internal/database"
)

type Chirp struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Body      string        `json:"body"`
	UserID    uuid.NullUUID `json:"user_id"`
}

func validationResponseHandler(w http.ResponseWriter, req *http.Request) {
	type validationRequest struct {
		Body string `json:"body"`
	}
	type validationResponse struct {
		Cleaned_body string `json:"cleaned_body,omitempty"`
	}
	thisRequest := validationRequest{}
	thisResponse := validationResponse{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&thisRequest)

	if err != nil {
		respondWithError(w, 400, "Something went wrong")
	}

	if len(thisRequest.Body) >= 140 {
		respondWithError(w, 400, "Chirpy too long")
	}

	cleanedBody, _ := cleanBody(thisRequest.Body)
	thisResponse.Cleaned_body = cleanedBody
	respondWithJSON(w, 200, thisResponse)
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, req *http.Request) {
	type validationRequest struct {
		Body   string        `json:"body"`
		UserId uuid.NullUUID `json:"user_id"`
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

	if thisRequest.UserId.Valid == false {
		respondWithError(w, 400, "User ID not valid")
		return
	}

	cleanedBody, _ := cleanBody(thisRequest.Body)
	params := database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: thisRequest.UserId,
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
