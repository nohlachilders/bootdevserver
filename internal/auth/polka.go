package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	_, exists := headers["Authorization"]
	if exists != true {
		return "", errors.New("No Authorization Header found")
	}

	authHeaderContents := strings.Split(headers.Get("Authorization"), " ")
	if len(authHeaderContents) != 2 {
		return "", errors.New("Malformed Authorization Header Contents")
	}
	if authHeaderContents[0] != "ApiKey" {
		return "", errors.New("Malformed Authorization Header Contents")
	}

	return authHeaderContents[1], nil
}
