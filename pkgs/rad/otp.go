package rad

import (
	"log"
	"regexp"

	"github.com/Abbas-gheydi/radotp/pkgs/authentiate"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

const (
	label_otp_stage = "otp"
)

func otpHandler(w radius.ResponseWriter, r *radius.Request) {
	paket := r.Packet
	username := rfc2865.UserName_GetString(paket)
	password := rfc2865.UserPassword_GetString(paket)
	state := rfc2865.State_GetString(paket)
	var code radius.Code

	stateInPool, _ := inMemoryPool.Lookup(username)
	if isStateValied(state, stateInPool) {
		//log.Println("state is ok")
		if IsOtpCodeValied(username, password) {
			code = AcceptUser(w, r, label_otp_stage)
			//delete user from statepool
			inMemoryPool.Delete(username)
		} else {

			code = RejectUser(w, r, label_otp_stage)
		}
		//state mismatch
	} else {
		code = RejectUser(w, r, label_otp_stage)
		log.Println("Warning, state mismatch for user", username)
	}
	log.Printf("%v to %v for %v", code, r.RemoteAddr, username)
}

func IsOtpCodeValied(username string, password string) bool {

	//otp password must be only numbers
	if !isSafeInput(password) {
		return false
	}

	if !isSafeInput(username) {
		return false
	}

	if authentiate.IsOtpCodeCurrect(username, password) {

		return true
	} else {
		return false
	}

}

func isSafeInput(input string) bool {
	var reg = regexp.MustCompile(`^[0-9A-Za-z_.@]{1,30}$`)
	return reg.MatchString(input)

}
