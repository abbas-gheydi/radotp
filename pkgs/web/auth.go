package web

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Abbas-gheydi/radotp/pkgs/storage"

	"github.com/golang-jwt/jwt/v4"
)

var jwtHmacSecret []byte

type WebAdminUserPass struct {
	User string
	Pass string
}

// generateRandomString creates a random 16-character string using MD5 hash
func generateRandomString() string {
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	hash := md5.Sum(randomBytes)
	return hex.EncodeToString(hash[:])
}

// validateJWT checks if the provided JWT token is valid and extracts the username from it
func validateJWT(tokenString string) (string, bool) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtHmacSecret, nil
	})

	if token == nil {
		return "", false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := fmt.Sprint(claims["user"])
		return user, true
	}
	return "", false
}

// generateJWT generates a JWT token for the specified user
func generateJWT(user string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(8 * time.Hour).Unix(), // Token expires in 8 hours
	})
	tokenString, _ := token.SignedString(jwtHmacSecret)
	return tokenString
}

// isCookieValid checks if the provided cookie contains a valid JWT token
func isCookieValid(cookieValue string) bool {
	if cookieValue == "" {
		return false
	}
	_, isValid := validateJWT(cookieValue)
	return isValid
}

// CheckWebAdminPass validates the provided username and password against the storage
func CheckWebAdminPass(username, password string) bool {
	return storage.ShaGenerator(password) == storage.GetAdminPassword(username)
}

// setCookie sets a secure, HTTP-only cookie for authentication
func setCookie(w http.ResponseWriter, username string) {
	cookie := &http.Cookie{
		Name:     "auth",
		Value:    generateJWT(username),
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
}

// setExpiredCookie sets an expired cookie to log the user out
func setExpiredCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:  "auth",
		Value: "expired",
		Path:  "/",
	}
	http.SetCookie(w, cookie)
}

// login handles the login functionality for web admins
func login(w http.ResponseWriter, r *http.Request) {
	loginTemplate := &templateHandler{filename: "login.gohtml"} // Template for login page

	if r.Method == http.MethodPost {
		username := r.FormValue("user")
		password := r.FormValue("pass")
		clientIP := strings.Split(r.RemoteAddr, ":")[0] // Extract client IP address

		if CheckWebAdminPass(username, password) {
			log.Printf("User '%s' logged in successfully from IP: %s\n", username, clientIP)
			setCookie(w, username) // Set auth cookie
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		log.Printf("Failed login attempt for user '%s' from IP: %s\n", username, clientIP)
	}

	// Render the login template
	loginTemplate.ServeHTTP(w, r)
}

// signOut handles user logout by setting an expired cookie
func signOut(w http.ResponseWriter, r *http.Request) {
	setExpiredCookie(w) // Expire the auth cookie
	http.Redirect(w, r, "/", http.StatusFound)
}
