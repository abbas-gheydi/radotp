package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type JsonUSer struct {
	UserName     string `json:"username"`
	Status       string `json:"status"`
	OtpCode      string `json:"otp_code,omitempty"`
	ResponseCode int    `json:"-"`
}

func getUserNameParamFromurl(r *http.Request) userCode {
	params := mux.Vars(r)
	userName := params["username"]
	user := userCode{UserName: userName}
	return user

}
func makeJsonResponse(w http.ResponseWriter, juser JsonUSer, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(juser)
}

func newjsonUser(userName string, status string, otpCode string) JsonUSer {
	return JsonUSer{UserName: userName, Status: status, OtpCode: otpCode}

}

func apiGetUser(w http.ResponseWriter, r *http.Request) {

	respCode := http.StatusOK

	user := getUserNameParamFromurl(r)
	searchuser(&user)

	userInJson := newjsonUser(user.UserName, user.Result, user.Qr)
	if user.Err != nil {
		respCode = http.StatusInternalServerError
	}

	makeJsonResponse(w, userInJson, respCode)

	//if user
	//fmt.Fprint(w, user.Result)

}
