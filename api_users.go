package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) userCreationHandler(w http.ResponseWriter, req *http.Request) {
	type userCreationRequest struct {
		Email string `json:"email"`
	}
	thisRequest := userCreationRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&thisRequest)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), thisRequest.Email)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Something went wrong: %s", err))
		return
	}
	respondWithJSON(w, 201, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
