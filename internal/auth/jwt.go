package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	thisToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			Subject:   userID.String(),
		},
	)
	signed, err := thisToken.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func ValidateJWT(tokenString string, tokenSecret string) (uuid.UUID, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, err
	}
	id, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuid.Parse(id)
}

func GetBearerToken(headers http.Header) (string, error) {
	_, exists := headers["Authorization"]
	if exists != true {
		return "", errors.New("No Authorization Header found")
	}

	authHeaderContents := strings.Split(headers.Get("Authorization"), " ")
	if len(authHeaderContents) != 2 {
		return "", errors.New("Malformed Authorization Header Contents")
	}
	if authHeaderContents[0] != "Bearer" {
		return "", errors.New("Malformed Authorization Header Contents")
	}

	return authHeaderContents[1], nil
}
