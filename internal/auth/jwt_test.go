package auth

import (
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
