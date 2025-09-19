package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header missing")
	}

	// Check if it's a ApiKey token
	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", errors.New("authorization header is not a ApiKey token")
	}

	// Extract the token
	token := strings.TrimPrefix(authHeader, "ApiKey ")
	return strings.TrimSpace(token), nil
}
