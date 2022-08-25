package web

import "net/http"

func serverHeader(w http.ResponseWriter, r *http.Request) {
	templ := templateHandler{filename: "header.gohtml"}
	templ.ServeHTTP(w, r)

}
