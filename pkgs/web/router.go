package web

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// StartRouter initializes the HTTP and HTTPS servers and routes.
func StartRouter() {
	// Initialize the secret for JWT
	jwtHmacSecret = []byte(generateRandomString())

	// Create a new router
	router := mux.NewRouter()

	// Static file handlers
	router.PathPrefix("/assets/").Handler(serveAssets(http.FileServer(http.FS(assets))))
	router.PathPrefix("/swager").Handler(http.FileServer(http.FS(swager)))

	// Authentication-related routes
	router.HandleFunc("/login/", login)
	router.HandleFunc("/sign_out/", signOut)

	// Protected routes
	router.Handle("/", MustAuth(manageUsers))
	router.Handle("/edit/", MustAuth(editAdminUser))
	router.Handle("/logs/", MustAuth(logs))
	router.Handle("/header/", MustAuth(serverHeader))

	// REST API routes (v1)
	apiBase := "/api/v1/user/{username}"
	router.Handle(apiBase, restApiMustAuth(apiGetUser)).Methods(http.MethodGet)
	router.Handle(apiBase, restApiMustAuth(apiCreateUser)).Methods(http.MethodPut)
	router.Handle(apiBase, restApiMustAuth(apiDeleteUser)).Methods(http.MethodDelete)
	router.Handle(apiBase, restApiMustAuth(apiUpdateUser)).Methods(http.MethodPost)

	// Log listening addresses
	log.Printf("HTTP Interface Listening on: %s\n", HTTPListenAddr)
	log.Printf("HTTPS Interface Listening on: %s\n", HTTPSListenAddr)

	// Start HTTP and HTTPS servers
	go startHTTPServer(router)
	startHTTPSServer(router)
}

// startHTTPServer starts the HTTP server.
// If HTTPS redirection is enabled, all HTTP requests are redirected to HTTPS.
func startHTTPServer(router *mux.Router) {
	if RedirectToHTTPS {
		redirectRouter := mux.NewRouter()
		redirectRouter.PathPrefix("/").HandlerFunc(redirectToHTTPS)
		err := http.ListenAndServe(HTTPListenAddr, redirectRouter)
		if err != nil {
			log.Printf("Failed to start HTTP server: %v\n", err)
		}
	} else {
		err := http.ListenAndServe(HTTPListenAddr, router)
		if err != nil {
			log.Printf("Failed to start HTTP server: %v\n", err)
		}
	}
}

// startHTTPSServer starts the HTTPS server.
func startHTTPSServer(router *mux.Router) {
	err := http.ListenAndServeTLS(HTTPSListenAddr, server_crt, server_key, router)
	if err != nil {
		log.Printf("Failed to start HTTPS server: %v\n", err)
	}
}

// redirectToHTTPS redirects HTTP requests to HTTPS.
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			host = r.Host
		} else {
			log.Printf("Error in redirectToHTTPS: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	// Build the HTTPS URL
	httpsURL := *r.URL
	httpsURL.Host = net.JoinHostPort(host, RedirectToHTTPSPortNumber)
	httpsURL.Scheme = "https"

	// Redirect the client to HTTPS
	http.Redirect(w, r, httpsURL.String(), http.StatusTemporaryRedirect)
}
