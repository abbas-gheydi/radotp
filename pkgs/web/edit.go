package web

import (
	"log"
	"net/http"

	"github.com/Abbas-gheydi/radotp/pkgs/storage"
)

func editAdminUser(w http.ResponseWriter, r *http.Request) {
	loginTemplate := &templateHandler{filename: "edit.gohtml"}

	if r.Method == http.MethodPost {
		currenPassword := r.FormValue("current")
		newPassword := r.FormValue("pwd")
		//loginTemplate.options = opts
		if CheckWebAdminPass("admin", currenPassword) && newPassword != "" {
			storage.SetAdminPassword(newPassword)
			log.Println("webadmin password updated")
			setExpiredCookie(w)
			http.Redirect(w, r, "/", http.StatusFound)
		}
	}

	loginTemplate.ServeHTTP(w, r)

}
