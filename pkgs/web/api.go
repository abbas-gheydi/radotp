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

type apiActions func(*userCode)

func apiActionsfunc(w http.ResponseWriter, r *http.Request, action apiActions, okResponseCode int) {

	user := getUserNameParamFromUrl(r)
	action(&user)
	respCode := createUserResponseHandler(&user, okResponseCode)
	userInJson := newjsonUser(user.UserName, user.Result, user.Code)
	makeJsonResponse(w, userInJson, respCode)

}

func apiGetUser(w http.ResponseWriter, r *http.Request) {
	apiActionsfunc(w, r, searchuser, http.StatusOK)
}

func apiCreateUser(w http.ResponseWriter, r *http.Request) {
	apiActionsfunc(w, r, createuser, http.StatusCreated)

}

func apiDeleteUser(w http.ResponseWriter, r *http.Request) {
	apiActionsfunc(w, r, deleteuser, http.StatusCreated)
}

func apiUpdateUser(w http.ResponseWriter, r *http.Request) {
	apiActionsfunc(w, r, updateuser, http.StatusCreated)
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
		respCode = http.StatusNotFound

	case disabled_OTP_Code_for_User:
		respCode = okResponseCode

	case already_exists:
		respCode = http.StatusMethodNotAllowed

	default:
		user.Result = "error"
		respCode = http.StatusInternalServerError

	}
	return
}
