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
	user_has_otp_code          string = "user has otp code"
	user_not_found             string = storage.User_not_found
	already_exists             string = "already exists"
	disabled_OTP_Code_for_User string = "Disabled OTP Code for User"
)

var ListenAddr = "0.0.0.0:8080"
var QrIssuer = "radotp"

//go:embed templates
var templates embed.FS

//go:embed assets
var assets embed.FS

type userCode struct {
	UserName string
	Code     string
	Qr       string
	Err      error
	Result   string
}

type authHandler struct {
	//next http.Handler
	next func(w http.ResponseWriter, r *http.Request)
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, readCookieError := r.Cookie("auth")

	if readCookieError == nil && isCookieValied(token.Value) {
		//authenticatin is successful
		h.next(w, r)

	} else {

		//log.Println("read cookie error", readCookieError)

		http.Redirect(w, r, "/login/", http.StatusFound)

	}

}

func MustAuth(handler func(w http.ResponseWriter, r *http.Request)) *authHandler {

	return &authHandler{next: handler}
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
	options  interface{}
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFS(templates, "templates/"+t.filename))
	})
	t.templ.Execute(w, t.options)
}

func createuser(user *userCode) {
	user.Code, user.Qr = authentiate.NewOtpUser(user.UserName, QrIssuer)
	user.Err = storage.Set(user.UserName, user.Code)
	if user.Err != nil {
		user.Code = ""
		user.Qr = ""
		if strings.Contains(user.Err.Error(), "duplicate key value violates unique constraint \"otps_username_key\"") {
			user.Err = fmt.Errorf(already_exists)
		}

		log.Println("createUser", user.Err)
	}
}

func updateuser(user *userCode) {
	user.Code, user.Qr = authentiate.NewOtpUser(user.UserName, QrIssuer)
	user.Err = storage.Update(user.UserName, user.Code)
	if user.Err != nil {
		user.Code = ""
		user.Qr = ""
		log.Println("updateUser", user.Err)
	}
}

func deleteuser(user *userCode) {
	user.Err = storage.Delete(user.UserName)
	if user.Err != nil {
		log.Println("deleteUser", user.Err)
	} else {
		user.Result = disabled_OTP_Code_for_User
	}
}
func searchuser(user *userCode) {
	SearchResualt, getErr := storage.Get(user.UserName)
	if SearchResualt == "" || getErr != nil {
		log.Println("searchUser error", getErr)
		user.Result = user_not_found
		if !strings.Contains(getErr.Error(), "record not found") {
			user.Err = getErr

		}

	} else {
		user.Result = user_has_otp_code
	}
}

func manageUsers(w http.ResponseWriter, r *http.Request) {
	templ := templateHandler{filename: "index.gohtml"}
	var opts userCode
	opts.UserName = r.FormValue("username")

	if r.Method == http.MethodPost {

		switch r.FormValue("tasks") {
		case "createuser":
			createuser(&opts)

		case "updateuser":
			updateuser(&opts)
		case "deleteuser":
			deleteuser(&opts)
		case "searchuser":
			searchuser(&opts)
		}

		templ.options = opts

	}

	templ.ServeHTTP(w, r)

}

func serveAssets(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
