package rad

import (
	"crypto/hmac"
	"crypto/md5"

	"github.com/Abbas-gheydi/radotp/pkgs/monitoring"
	"github.com/Abbas-gheydi/radotp/pkgs/rad/vendors"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"layeh.com/radius/rfc2869"
)

var EnableMetrics bool

func AcceptUser(w radius.ResponseWriter, r *radius.Request, stage string) (code radius.Code) {
	code = radius.CodeAccessAccept
	packet := r.Packet.Response(code)

	// Add Message-Authenticator placeholder to the response packet at the first place
	placeholder := make([]byte, 16) // 16 bytes of zeros
	rfc2869.MessageAuthenticator_Set(packet, placeholder)

	if state := rfc2865.State_GetString(r.Packet); state != "" {
		rfc2865.State_SetString(packet, state)
	}
	AddVendorGroup(packet, r)

	// DO NOT Add another attribute after Message-Authenticator, it may cause the Message-Authenticator to be invalid
	AddMessageAuthenticator(packet)

	w.Write(packet)
	go sendToMonitoring(rfc2865.UserName_GetString(r.Packet), stage, "Accept")

	return
}

func RejectUser(w radius.ResponseWriter, r *radius.Request, stage string) (code radius.Code) {
	code = radius.CodeAccessReject
	packet := r.Packet.Response(code)

	// Add Message-Authenticator placeholder to the response packet at the first place
	placeholder := make([]byte, 16) // 16 bytes of zeros
	rfc2869.MessageAuthenticator_Set(packet, placeholder)

	//set state in response packet
	if state := rfc2865.State_GetString(r.Packet); state != "" {
		rfc2865.State_SetString(packet, state)
	}

	// DO NOT Add another attribute after Message-Authenticator, it may cause the Message-Authenticator to be invalid
	AddMessageAuthenticator(packet)

	w.Write(packet)
	go sendToMonitoring(rfc2865.UserName_GetString(r.Packet), stage, "Reject")

	return
}

func SendForChalenge(w radius.ResponseWriter, r *radius.Request, stage string) (code radius.Code) {
	packet := r.Packet.Response(radius.CodeAccessChallenge)
	username := rfc2865.UserName_GetString(r.Packet)
	stateInOurPool, _ := states.Lookup(username)

	// Add Message-Authenticator placeholder to the response packet at the first place
	placeholder := make([]byte, 16) // 16 bytes of zeros
	rfc2869.MessageAuthenticator_Set(packet, placeholder)

	//stateInOurPool, groupInOurpoll := inMemoryPool.Lookup(username)

	//pool.Insert(username)
	//log.Println(username, " state ", stateInOurPool, " user groups is ", groupInOurpoll)

	rfc2865.State_SetString(packet, stateInOurPool)
	rfc2865.ReplyMessage_Set(packet, []byte("ENTER OTP CODE"))

	// DO NOT Add another attribute after Message-Authenticator, it may cause the Message-Authenticator to be invalid
	AddMessageAuthenticator(packet)

	w.Write(packet)
	go sendToMonitoring(rfc2865.UserName_GetString(r.Packet), stage, "Challenge")

	return radius.CodeAccessChallenge
}

func AddVendorGroup(p *radius.Packet, r *radius.Request) {
	_, usergroup := states.Lookup(rfc2865.UserName_GetString(r.Packet))
	if RadiusConfigs.Enable_Fortinet_Group_Name {
		vendors.FortinetGroupName_SetString(p, usergroup)
	}
}

// AddMessageAuthenticator adds the Message-Authenticator attribute to a RADIUS packet.
// Message-Authenticator = HMAC-MD5 (Type, Identifier, Length, Request Authenticator, Attributes)
func AddMessageAuthenticator(packet *radius.Packet) error {
	// Step 1: Set the Message-Authenticator field to zero.
	placeholder := make([]byte, 16) // 16 bytes of zeros
	err := rfc2869.MessageAuthenticator_Set(packet, placeholder)
	if err != nil {
		return err
	}

	// Step 2: Marshal the packet to get the binary representation
	b, err := packet.MarshalBinary()
	if err != nil {
		return err
	}

	// Step 3: Calculate HMAC-MD5 with the shared secret
	mac := hmac.New(md5.New, []byte(RadiusConfigs.Secret))
	mac.Write(b)

	// Step 4: Get the authenticator (HMAC result)
	authenticator := mac.Sum(nil)

	// Step 5: Set the Message-Authenticator field
	err = rfc2869.MessageAuthenticator_Set(packet, authenticator)
	if err != nil {
		return err
	}

	return nil
}

func sendToMonitoring(username string, stage string, stat string) {
	if EnableMetrics {
		switch stat {
		case "Accept":
			monitoring.Accepted_users.Append(username, stage)
		case "Reject":
			monitoring.Rejected_users.Append(username, stage)
		case "Challenge":
			monitoring.Chalenged_users.Append(username, stage)

		}
	}

}
