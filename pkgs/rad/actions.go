package rad

import (
	"github.com/Abbas-gheydi/radotp/pkgs/monitoring"
	"github.com/Abbas-gheydi/radotp/pkgs/rad/vendors"

	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
)

var EnableMetrics bool

func AcceptUser(w radius.ResponseWriter, r *radius.Request, stage string) (code radius.Code) {
	code = radius.CodeAccessAccept
	paket := r.Packet.Response(code)
	if state := rfc2865.State_GetString(r.Packet); state != "" {
		rfc2865.State_SetString(paket, state)

	}
	AddVendorGroup(paket, r)
	w.Write(paket)
	go sendToMonitoring(rfc2865.UserName_GetString(r.Packet), stage, "Accept")

	return
}

func RejectUser(w radius.ResponseWriter, r *radius.Request, stage string) (code radius.Code) {
	code = radius.CodeAccessReject
	paket := r.Packet.Response(code)
	//set state in response packet
	if state := rfc2865.State_GetString(r.Packet); state != "" {
		rfc2865.State_SetString(paket, state)
	}

	w.Write(paket)
	go sendToMonitoring(rfc2865.UserName_GetString(r.Packet), stage, "Reject")

	return
}

func SendForChalenge(w radius.ResponseWriter, r *radius.Request, stage string) (code radius.Code) {
	paket := r.Packet.Response(radius.CodeAccessChallenge)
	username := rfc2865.UserName_GetString(r.Packet)
	stateInOurPool, _ := inMemoryPool.Lookup(username)
	//stateInOurPool, groupInOurpoll := inMemoryPool.Lookup(username)

	//pool.Insert(username)
	//log.Println(username, " state ", stateInOurPool, " user groups is ", groupInOurpoll)

	rfc2865.State_SetString(paket, stateInOurPool)
	rfc2865.ReplyMessage_Set(paket, []byte("ENTER OTP CODE"))

	w.Write(paket)
	go sendToMonitoring(rfc2865.UserName_GetString(r.Packet), stage, "Challenge")

	return radius.CodeAccessChallenge

}

func AddVendorGroup(p *radius.Packet, r *radius.Request) {
	_, usergroup := inMemoryPool.Lookup(rfc2865.UserName_GetString(r.Packet))

	if RadiusConfigs.Enable_Fortinet_Group_Name {

		vendors.FortinetGroupName_SetString(p, usergroup)
	}

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
