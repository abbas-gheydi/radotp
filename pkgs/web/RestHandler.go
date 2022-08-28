package web

import (
	"crypto/sha512"
	"crypto/subtle"
	"fmt"
	"net/http"
)

var (
	bearerPrefix  = "Bearer "
	ApiKey        = ""
	EnableRestApi bool
)

func restApiMustAuth(handler func(w http.ResponseWriter, r *http.Request)) *RestAuthHandler {

	return &RestAuthHandler{next: handler}
}

type RestAuthHandler struct {
	next func(w http.ResponseWriter, r *http.Request)
}

func (h *RestAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if isRestReqAuthorized(w, r) {
		h.next(w, r)

	} else {
		w.WriteHeader(http.StatusForbidden)

		fmt.Fprint(w, "Not Authorized")

	}

}

func secureCompare(apikey string, userkey string) bool {
	givenSha := sha512.Sum512([]byte(apikey))
	actualSha := sha512.Sum512([]byte(userkey))

	return subtle.ConstantTimeCompare(givenSha[:], actualSha[:]) == 1
}

func isRestReqAuthorized(w http.ResponseWriter, r *http.Request) bool {
	var bearerkey = bearerPrefix + ApiKey
	userToken := r.Header.Get("Authorization")
	return secureCompare(bearerkey, userToken) && ApiKey != "" && EnableRestApi
}
