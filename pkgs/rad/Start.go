package rad

import (
	"log"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

type RadConfs struct {
	ListenAddress              string
	Secret                     string
	Enable_Fortinet_Group_Name bool
	Authentication_Mode        string //only_otp, only_password, two_fa
}

const (
	only_password       = "only_password"
	only_otp            = "only_otp"
	two_fa              = "two_fa"
	two_fa_optional_otp = "two_fa_optional_otp"
)

var RadiusConfigs RadConfs

func StartRadius() {

	inMemoryPool.Init()

	handler := func(w radius.ResponseWriter, r *radius.Request) {
		log.Print("received a radius packet for user: ", rfc2865.UserName_GetString(r.Packet))
		if rfc2865.UserPassword_GetString(r.Packet) == "" {
			log.Println("password is empty for user: ", rfc2865.UserName_GetString(r.Packet))
		}

		if mustCheckPassword() {
			User_PassHandler(w, r)

		} else {
			otpHandler(w, r)

		}

	}

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(RadiusConfigs.Secret)),
		Addr:         RadiusConfigs.ListenAddress,
	}

	log.Printf("Radius server Listen on %v", RadiusConfigs.ListenAddress)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func mustCheckPassword() bool {
	// check setting
	if RadiusConfigs.Authentication_Mode == only_otp {
		return false
	}
	return true
}
