package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	Token          string    `json:"token,omitempty"`
	RefreshToken   string    `json:"refresh_token,omitempty"`
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	thisRequest := userLoginRequest{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&thisRequest)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	user, err := cfg.db.GetUserByEmail(req.Context(), thisRequest.Email)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	match := auth.CheckPasswordHash(thisRequest.Password, user.HashedPassword)
	if match != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	timeOneHour := time.Duration(1) * time.Hour
	token, err := auth.MakeJWT(user.ID, cfg.secret, timeOneHour)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	refreshString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}
	params := database.CreateRefreshTokenParams{
		Token:     refreshString,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
		UserID:    user.ID,
	}
	refresh, err := cfg.db.CreateRefreshToken(req.Context(), params)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	respondWithJSON(w, 200, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refresh.Token,
	})
}

func (cfg *apiConfig) userRefreshHandler(w http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}
	refreshToken, err := cfg.db.GetRefreshToken(req.Context(), tokenString)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, 401, "Token is expired")
		return
	}
	if refreshToken.RevokedAt.Valid {
		respondWithError(w, 401, "Token is revoked")
		return
	}

	user, err := cfg.db.GetUserByID(req.Context(), refreshToken.UserID)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	timeOneHour := time.Duration(1) * time.Hour
	token, err := auth.MakeJWT(user.ID, cfg.secret, timeOneHour)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}

	type responseShape struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, 200, responseShape{Token: token})
}

func (cfg *apiConfig) userRevokeHandler(w http.ResponseWriter, req *http.Request) {
	tokenString, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}
	refreshToken, err := cfg.db.GetRefreshToken(req.Context(), tokenString)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Something went wrong: %s", err))
		return
	}
	err = cfg.db.RevokeRefreshToken(req.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Something went wrong: %s", err))
		return
	}
	// respond with 204
	w.WriteHeader(204)
}
