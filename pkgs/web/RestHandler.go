package web

import (
	"crypto/sha512"
	"crypto/subtle"
	"fmt"
	"net/http"
)

var (
	bearerPrefix  = "Bearer "
	ApiKey        string
	EnableRestApi bool
)

// restApiMustAuth wraps a handler with REST API authentication.
func restApiMustAuth(handler func(w http.ResponseWriter, r *http.Request)) *RestAuthHandler {
	return &RestAuthHandler{next: handler}
}

// RestAuthHandler is a middleware struct for handling REST API authentication.
type RestAuthHandler struct {
	next func(w http.ResponseWriter, r *http.Request)
}

// ServeHTTP checks if the request is authorized and either proceeds to the next handler
// or responds with a 403 Forbidden status.
func (h *RestAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if isRestRequestAuthorized(r) {
		h.next(w, r) // Call the next handler
	} else {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Not Authorized")
	}
}

// secureCompare performs a constant-time comparison of two SHA-512 hashed strings.
// This helps prevent timing attacks.
func secureCompare(apiKey string, userKey string) bool {
	expectedHash := sha512.Sum512([]byte(apiKey))
	providedHash := sha512.Sum512([]byte(userKey))
	return subtle.ConstantTimeCompare(expectedHash[:], providedHash[:]) == 1
}

// isRestRequestAuthorized checks if the incoming request is authorized
// by comparing the "Authorization" header with the expected API key.
func isRestRequestAuthorized(r *http.Request) bool {
	// Construct the expected Bearer token
	expectedBearerToken := bearerPrefix + ApiKey

	// Extract the "Authorization" header from the request
	userToken := r.Header.Get("Authorization")

	// Verify the token and ensure REST API is enabled
	return secureCompare(expectedBearerToken, userToken) && ApiKey != "" && EnableRestApi
}
