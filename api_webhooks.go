package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

	if thisRequest.Event != "user.upgrade" {
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

	w.WriteHeader(204)
}
