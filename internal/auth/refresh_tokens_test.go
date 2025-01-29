package auth

import "testing"

func TestRefreshTokenGen(t *testing.T) {
	token, err := MakeRefreshToken()
	if err != nil {
		t.Error("error creating token")
	}
	if len(token) == 0 {
		t.Error("empty token")
	}
}
