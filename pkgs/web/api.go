package web

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type JsonUSer struct {
	UserName     string `json:"username"`
	Result       string `json:"result"`
	OtpCode      string `json:"otp_code,omitempty"`
	ResponseCode int    `json:"-"`
}

func apiGetUser(w http.ResponseWriter, r *http.Request) {

	user := getUserNameParamFromUrl(r)
	searchuser(&user)
	respCode := createUserResponseHandler(&user, http.StatusOK)
	userInJson := newjsonUser(user.UserName, user.Result, user.Code)
	makeJsonResponse(w, userInJson, respCode)

}

func apiCreateUser(w http.ResponseWriter, r *http.Request) {

	user := getUserNameParamFromUrl(r)
	createuser(&user)
	respCode := createUserResponseHandler(&user, http.StatusCreated)
	userInJson := newjsonUser(user.UserName, user.Result, user.Code)
	makeJsonResponse(w, userInJson, respCode)

}

func getUserNameParamFromUrl(r *http.Request) userCode {
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
	return JsonUSer{UserName: userName, Result: status, OtpCode: otpCode}

}

func createUserResponseHandler(user *userCode, okResponseCode int) (respCode int) {
	if user.Err != nil {
		user.Result = user.Err.Error()
	}
	switch user.Result {

	case "":
		user.Result = "true"
		respCode = okResponseCode

	case user_has_otp_code:
		respCode = okResponseCode

	case user_not_found:
		respCode = okResponseCode

	case already_exists:
		respCode = http.StatusMethodNotAllowed

	default:
		user.Result = "false"
		respCode = http.StatusInternalServerError

	}
	return
}
