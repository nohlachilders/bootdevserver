package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

// test for basic usage of MakeJWT and ValidateJWT
func TestJWT(t *testing.T) {
	id, err := uuid.NewUUID()
	if err != nil {
		t.Error("error creating uuid")
	}
	tokenSecret := "test"
	token, err := MakeJWT(id, tokenSecret, time.Minute)
	if err != nil {
		t.Error("error creating token")
	}
	parsedId, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Errorf("error validating token: %v", err)
	}
	if parsedId != id {
		t.Error("parsed id does not match")
	}
}

func TestExpiration(t *testing.T) {
	id, err := uuid.NewUUID()
	if err != nil {
		t.Error("error creating uuid")
	}
	tokenSecret := "test"
	token, err := MakeJWT(id, tokenSecret, time.Microsecond)
	if err != nil {
		t.Error("error creating token")
	}
	time.Sleep(3 * time.Microsecond)
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Errorf("expected error, found none. expired token accepted")
	}
}

func TestSecret(t *testing.T) {
	id, err := uuid.NewUUID()
	if err != nil {
		t.Error("error creating uuid")
	}
	tokenSecret := "test"
	token, err := MakeJWT(id, tokenSecret, time.Minute)
	if err != nil {
		t.Error("error creating token")
	}

	incorrectSecret := "test2"
	_, err = ValidateJWT(token, incorrectSecret)
	if err == nil {
		t.Errorf("expected error, found none. incorrect token accepted")
	}
}

func TestBearerTokenExtraction(t *testing.T) {
	goodRequest, err := http.NewRequest("GET", "localhost", nil)
	if err != nil {
		t.Errorf("error found, expected none: %s", err)
	}

	goodRequest.Header.Add("Authorization", "Bearer token_string")
	_, err = GetBearerToken(goodRequest.Header)
	if err != nil {
		t.Errorf("error found, expected none: %s", err)
	}

	badRequest, err := http.NewRequest("GET", "localhost", nil)
	if err != nil {
		t.Errorf("error found, expected none: %s", err)
	}
	badRequest.Header.Add("Authorization", "token_string")
	_, err = GetBearerToken(badRequest.Header)
	if err == nil {
		t.Errorf("no error found, expected one")
	}

	emptyRequest, err := http.NewRequest("GET", "localhost", nil)
	if err != nil {
		t.Errorf("error found, expected none: %s", err)
	}
	_, err = GetBearerToken(emptyRequest.Header)
	if err == nil {
		t.Errorf("no error found, expected one")
	}
}
