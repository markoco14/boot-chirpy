package auth

import (
	"testing"
	"time"
	"github.com/google/uuid"
)

// Test creating a valid JWT and successfully validating it.
func TestJWTCreationAndValidation(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := time.Hour

	signedToken, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("error making JWT: %v", err)
	}

	parsedID, err := ValidateJWT(signedToken, tokenSecret)
	if err != nil {
		t.Fatalf("error validating JWT: %v", err)
	}

	if parsedID != userID {
		t.Fatalf("expected %v, got %v", userID, parsedID)
	}
}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "mysecret"
	expiresIn := -time.Minute

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("error making JWT: %v", err)
	}

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}