package web

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Abbas-gheydi/radotp/pkgs/authentiate"
	"github.com/Abbas-gheydi/radotp/pkgs/storage"
)

const (
	UserHasOtpCode         string = "user has otp code"
	UserNotFound           string = storage.User_not_found
	AlreadyExists          string = "already exists"
	DisabledOtpCodeForUser string = "Disabled OTP Code for User"
)

var (
	HTTPListenAddr = "0.0.0.0:8080"
	QrIssuer       = "radotp"
	//go:embed templates
	templates embed.FS
	//go:embed assets
	assets embed.FS
	//go:embed swager
	swager embed.FS
)

type userCode struct {
	UserName string
	Code     string
	Qr       string
	Err      error
	Result   string
}

// authHandler is a custom handler for authentication checks
type authHandler struct {
	next func(w http.ResponseWriter, r *http.Request)
}

// ServeHTTP handles the authentication for each request
func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("auth")
	if err == nil && isCookieValid(token.Value) {
		// Proceed to next handler if authentication is successful
		h.next(w, r)
	} else {
		// Redirect to login if not authenticated
		http.Redirect(w, r, "/login/", http.StatusFound)
	}
}

// MustAuth returns an authHandler to wrap around routes requiring authentication
func MustAuth(handler func(w http.ResponseWriter, r *http.Request)) *authHandler {
	return &authHandler{next: handler}
}

// templateHandler handles the loading and rendering of templates
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
	options  interface{}
}

// ServeHTTP renders the template and serves the response
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		// Parse the template only once
		t.templ = template.Must(template.ParseFS(templates, "templates/"+t.filename))
	})
	// Execute the template with options passed to it
	t.templ.Execute(w, t.options)
}

// createUser creates a new OTP user and stores it
func createUser(user *userCode) {
	user.Code, user.Qr = authentiate.NewOtpUser(user.UserName, QrIssuer)
	user.Err = storage.Set(user.UserName, user.Code)

	if user.Err != nil {
		user.Code = ""
		user.Qr = ""
		if strings.Contains(user.Err.Error(), "duplicate key value violates unique constraint \"otps_username_key\"") {
			user.Err = fmt.Errorf(AlreadyExists)
		}
		log.Println("createUser", user.Err)
	}
}

// updateUser updates an existing OTP user
func updateUser(user *userCode) {
	user.Code, user.Qr = authentiate.NewOtpUser(user.UserName, QrIssuer)
	user.Err = storage.Update(user.UserName, user.Code)
	if user.Err != nil {
		user.Code = ""
		user.Qr = ""
		log.Println("updateUser", user.Err)
	}
}

// deleteUser deletes the OTP code for a user
func deleteUser(user *userCode) {
	user.Err = storage.Delete(user.UserName)
	if user.Err != nil {
		log.Println("deleteUser", user.Err)
	} else {
		user.Result = DisabledOtpCodeForUser
	}
}

// searchUser searches for an OTP code for a user
func searchUser(user *userCode) {
	searchResult, getErr := storage.Get(user.UserName)
	if searchResult == "" || getErr != nil {
		log.Println("searchUser error", getErr)
		user.Result = UserNotFound
		if !strings.Contains(getErr.Error(), "record not found") {
			user.Err = getErr
		}
	} else {
		user.Result = UserHasOtpCode
	}
}

// manageUsers handles user management actions (create, update, delete, search)
func manageUsers(w http.ResponseWriter, r *http.Request) {
	templ := templateHandler{filename: "index.gohtml"}
	var opts userCode
	opts.UserName = r.FormValue("username")

	if r.Method == http.MethodPost {
		// Determine action based on form submission
		switch r.FormValue("tasks") {
		case "createuser":
			createUser(&opts)
		case "updateuser":
			updateUser(&opts)
		case "deleteuser":
			deleteUser(&opts)
		case "searchuser":
			searchUser(&opts)
		}
		templ.options = opts
	}

	// Render the template
	templ.ServeHTTP(w, r)
}

// serveAssets handles requests for static assets
func serveAssets(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure no directory listing
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		// Serve the static asset
		next.ServeHTTP(w, r)
	})
}
