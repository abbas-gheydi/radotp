package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func apiGetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userName := params["userName"]
	user := userCode{UserName: userName}

	searchuser(&user)
	fmt.Fprint(w, user.Result)

}
