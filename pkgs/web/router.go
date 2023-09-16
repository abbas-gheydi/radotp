package web

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func StartRouter() {

	jwtHmacSecret = []byte(generateRandomString())
	router := mux.NewRouter()

	//http.Handle("/assets/", http.FileServer(http.FS(assets)))
	fs := http.FileServer(http.FS(assets))
	router.PathPrefix("/assets/").Handler(serveAssets(fs))

	swgr := http.FileServer(http.FS(swager))
	router.PathPrefix("/swager").Handler(swgr)

	router.HandleFunc("/login/", login)
	router.HandleFunc("/sign_out/", signOut)

	router.Handle("/", MustAuth(manageUsers))
	router.Handle("/edit/", MustAuth(editAdminUser))
	router.Handle("/logs/", MustAuth(logs))
	router.Handle("/header/", MustAuth(serverHeader))

	router.Handle("/api/v1/user/{username}", restApiMustAuth(apiGetUser)).Methods(http.MethodGet)
	router.Handle("/api/v1/user/{username}", restApiMustAuth(apiCreateUser)).Methods(http.MethodPut)
	router.Handle("/api/v1/user/{username}", restApiMustAuth(apiDeleteUser)).Methods(http.MethodDelete)
	router.Handle("/api/v1/user/{username}", restApiMustAuth(apiUpdateUser)).Methods(http.MethodPost)
	log.Println("HTTP Interface Listen on:", HTTPListenAddr)
	log.Println("HTTPS Interface Listen on:", HTTPSListenAddr)
	go httpListenAndServe(router)
	err := http.ListenAndServeTLS(HTTPSListenAddr, server_crt, server_key, router)
	if err != nil {
		log.Println(err)
	}
}

func httpListenAndServe(defaultRouter *mux.Router) {
	if RedirectToHTTPS {
		redirecrtRouter := mux.NewRouter()
		redirecrtRouter.PathPrefix("/").HandlerFunc(redirectToHttps)
		http.ListenAndServe(HTTPListenAddr, redirecrtRouter)
	} else {
		err := http.ListenAndServe(HTTPListenAddr, defaultRouter)
		if err != nil {
			log.Println(err)
		}
	}
}

func redirectToHttps(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.Host)
	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			host = r.Host
		} else {
			log.Println("redirectToHttps", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))

		}
	}

	u := r.URL
	u.Host = net.JoinHostPort(host, RedirectToHTTPSPortNumber)
	u.Scheme = "https"
	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)

}
