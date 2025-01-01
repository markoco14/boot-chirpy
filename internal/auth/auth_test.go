package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name string
		header http.Header
		expectedToken string
		expectingError bool
	} {
		{
			name: "valid header",
			header: http.Header{"Authorization": {"Bearer mytoken"}},
			expectedToken: "mytoken",
			expectingError: false,
		},
		{
			name: "missing header",
			header: http.Header{},
			expectedToken: "",
			expectingError: true,
		},
		{
			name: "no bearer prefix",
			header: http.Header{"Authorization": {"mytoken"}},
			expectedToken: "",
			expectingError: true,
		},
		{
			name: "empty token",
			header: http.Header{"Authorization": {"Bearer "}},
			expectedToken: "",
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.header)
			if token != tt.expectedToken {
				t.Errorf("expected token %s, got %s", tt.expectedToken, token)
			}
			if err != nil && !tt.expectingError {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}