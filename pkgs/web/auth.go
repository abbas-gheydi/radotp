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

	"github.com/golang-jwt/jwt"
)

var jwtHmacSecret []byte

type WebAdminUserPass struct {
	User string
	Pass string
}

func generateRandomString() string {
	randomByte := make([]byte, 8)
	rand.Read(randomByte)
	hash := md5.Sum([]byte(randomByte))
	return hex.EncodeToString(hash[:])

}

func validate_jwt(tokenString string) (user string, isValied bool) {

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
		//fmt.Println(claims["user"], claims["nbf"])
		user = fmt.Sprint(claims["user"])
		return user, true
	} else {
		//log.Println(err)
		return "", false
	}

}

func generate_jwt(user string) string {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour * 8).Unix(),
	})

	tokenString, _ := token.SignedString(jwtHmacSecret)

	//fmt.Println(tokenString, err)
	return tokenString
}

func isCookieValied(postedCookie string) bool {
	if postedCookie == "" {
		return false
	} else {
		_, status := validate_jwt(postedCookie)
		return status

	}
}

func CheckWebAdminPass(username string, pssword string) bool {
	if storage.ShaGenerator(pssword) == storage.GetAdminPassword(username) {
		return true

	} else {
		return false
	}
}

func setCookie(w http.ResponseWriter, uname string) {
	cookie := &http.Cookie{Name: "auth",
		Value: generate_jwt(uname),
		Path:  "/",
	}
	http.SetCookie(w, cookie)

}

func setExpiredCookie(w http.ResponseWriter, uname string) {
	cookie := &http.Cookie{Name: "auth",
		Value: "expired",
		Path:  "/",
	}
	http.SetCookie(w, cookie)

}

func login(w http.ResponseWriter, r *http.Request) {
	loginTemplate := &templateHandler{filename: "login.gohtml"}

	//var opts WebAdminUserPass
	if r.Method == http.MethodPost {
		user := r.FormValue("user")
		pass := r.FormValue("pass")
		ip := strings.Split(r.RemoteAddr, ":")[0]
		//loginTemplate.options = opts
		if CheckWebAdminPass(user, pass) {
			log.Println(user, "is logged in from:", ip)
			setCookie(w, user)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			log.Println(user, "failed login attempts from:", ip)

		}
	}

	loginTemplate.ServeHTTP(w, r)

}

func editAdminUser(w http.ResponseWriter, r *http.Request) {
	loginTemplate := &templateHandler{filename: "edit.gohtml"}

	if r.Method == http.MethodPost {
		currenPassword := r.FormValue("current")
		newPassword := r.FormValue("pwd")
		//loginTemplate.options = opts
		if CheckWebAdminPass("admin", currenPassword) {
			storage.SetAdminPassword(newPassword)
			log.Println("webadmin password updated")
			setExpiredCookie(w, "user")
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	loginTemplate.ServeHTTP(w, r)

}

func signOut(w http.ResponseWriter, r *http.Request) {
	setExpiredCookie(w, "user")
	http.Redirect(w, r, "/", http.StatusFound)

}
