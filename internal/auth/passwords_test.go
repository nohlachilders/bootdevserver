package auth

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	pass := "test"
	hashed, _ := HashPassword(pass)
	t.Log(hashed)
}

func TestCheckPassword(t *testing.T) {
	pass := "test"
	hashed := "$2a$10$5CFy00EAVTrbdRNbLHR0z.dnTbGp1X9E.5STOeKn95h/kcjjAFlcC"
	t.Log(CheckPasswordHash(pass, hashed))
}
