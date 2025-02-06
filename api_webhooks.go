package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/nohlachilders/bootdevserver/internal/auth"
)

func (cfg *apiConfig) updateRedHandler(w http.ResponseWriter, req *http.Request) {
	type updateRedRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	thisRequest := updateRedRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&thisRequest)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	key, err := auth.GetAPIKey(req.Header)
	if err != nil || key != cfg.polkaSecret {
		respondWithError(w, 401, fmt.Sprintf("Unauthorized"))
		return
	}

	if thisRequest.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	id, err := uuid.Parse(thisRequest.Data.UserID)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}
	user, err := cfg.db.UpdateUserRed(req.Context(), id)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	user.Email = ""

	respondWithJSON(w, 204, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		IsRed:     user.IsRed,
	})
}
