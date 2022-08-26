package web

import (
	"fmt"
	"net/http"
)

func isRestReqAuthorized(w http.ResponseWriter, r *http.Request) bool {
	//check for api key
	return true
}

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
