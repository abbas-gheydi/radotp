package web

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func StartRouter() {

	jwtHmacSecret = []byte(generateRandomString())
	router := mux.NewRouter()

	//http.Handle("/assets/", http.FileServer(http.FS(assets)))
	fs := http.FileServer(http.FS(assets))
	router.PathPrefix("/assets/").Handler(serveAssets(fs))
	router.HandleFunc("/login/", login)
	router.HandleFunc("/sign_out/", signOut)

	router.Handle("/", MustAuth(manageUsers))
	router.Handle("/edit/", MustAuth(editAdminUser))
	router.Handle("/logs/", MustAuth(logs))
	router.Handle("/header/", MustAuth(serverHeader))

	router.Handle("/api/v1/{username}", restApiMustAuth(apiGetUser)).Methods(http.MethodGet)
	router.Handle("/api/v1/{username}", restApiMustAuth(apiCreateUser)).Methods(http.MethodPut)

	log.Println("Web Interface Listen on:", ListenAddr)

	http.ListenAndServe(ListenAddr, router)

}