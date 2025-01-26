package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/nohlachilders/bootdevserver/internal/auth"
	"github.com/nohlachilders/bootdevserver/internal/database"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Token          string    `json:"token"`
}

func (cfg *apiConfig) userCreationHandler(w http.ResponseWriter, req *http.Request) {
	type userCreationRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	thisRequest := userCreationRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&thisRequest)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	hashed, err := auth.HashPassword(thisRequest.Password)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Something went wrong: %s", err))
		return
	}
	err = auth.CheckPasswordHash(thisRequest.Password, hashed)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	params := database.CreateUserParams{
		Email:          thisRequest.Email,
		HashedPassword: hashed,
	}
	user, err := cfg.db.CreateUser(req.Context(), params)
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

func (cfg *apiConfig) userLoginHandler(w http.ResponseWriter, req *http.Request) {
	type userLoginRequest struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds string `json:"expires_in_seconds"`
	}
	thisRequest := userLoginRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&thisRequest)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), thisRequest.Email)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	match := auth.CheckPasswordHash(thisRequest.Password, user.HashedPassword)
	if match != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	timeExpires, err := strconv.Atoi(thisRequest.ExpiresInSeconds)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Something went wrong: %s", err))
	}
	if 0 > timeExpires || timeExpires > 3600 {
		timeExpires = 3600
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(timeExpires)*time.Second)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Something went wrong: %s", err))
	}

	respondWithJSON(w, 200, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}
