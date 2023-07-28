package rad

import (
	"log"

	"github.com/Abbas-gheydi/radotp/pkgs/authentiate"
	"github.com/Abbas-gheydi/radotp/pkgs/storage"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

// set ldap as auth provider
var Auth_Provider authentiate.LdapProvider

const (
	label_ldap_stage = "active directory"
)

func User_PassHandler(w radius.ResponseWriter, r *radius.Request) {

	paket := r.Packet
	username := rfc2865.UserName_GetString(paket)
	password := rfc2865.UserPassword_GetString(paket)
	var code radius.Code

	if IsUserPassValied(Auth_Provider, username, password) {
		//check config
		switch RadiusConfigs.Authentication_Mode {
		case only_password:
			code = AcceptUser(w, r, label_ldap_stage)
		case two_fa_optional_otp:
			if !storage.IsUserExist(username) {
				code = AcceptUser(w, r, label_ldap_stage)
			} else {
				code = SendForChalenge(w, r, label_ldap_stage)
			}

		case two_fa:
			code = SendForChalenge(w, r, label_ldap_stage)

		}

	} else {
		//Wrong user and pass
		code = RejectUser(w, r, label_ldap_stage)
	}

	log.Printf("%v to %v for %v stage %v", code, r.RemoteAddr, username, label_ldap_stage)

}

func IsUserPassValied(auth_provider authentiate.Auth_Provider, username string, password string) bool {
	if !isSafeInput(username) {
		return false
	}
	var usergroup string
	authe_state, group := auth_provider.IsUserAuthenticated(username, password)
	//log.Printf("auth state %v group %v", authe_state, group)
	if authe_state {
		//user pass is ok
		//ToDO: Check if its necessary to have multiple groups value
		if len(group) == 0 {
			usergroup = ""
		} else {
			usergroup = group[0]
		}
		if RadiusConfigs.Authentication_Mode == two_fa || RadiusConfigs.Authentication_Mode == two_fa_optional_otp {
			states.Insert(username, usergroup)
		}
		return true

	}
	// user pass is invalied
	return false

}
