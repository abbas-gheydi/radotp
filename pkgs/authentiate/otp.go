package authentiate

import (
	"time"

	"github.com/Abbas-gheydi/radotp/pkgs/storage"

	"github.com/xlzd/gotp"
)

type otpProvider struct {
	OtpUser  *gotp.TOTP
	secret   string
	username string
	err      error
}

func (o *otpProvider) getSecret() {
	//o.secret := db.getsecret(o.username)
	//o.secret = "TTVJP2C4XPLJZDVSGORIFHE552"
	o.secret, o.err = storage.Get(o.username)

}
func (o *otpProvider) make() {
	o.getSecret()
	o.OtpUser = gotp.NewDefaultTOTP(o.secret)

}

func IsOtpCodeCurrect(username string, otpcode string) bool {
	var user otpProvider
	user.username = username
	user.make()
	if user.err != nil {
		return false
	}
	return user.OtpUser.Verify(otpcode, int(time.Now().Unix()))

}
func NewOtpUser(username string, issuerName string) (secret string, qrcode string) {
	secret = gotp.RandomSecret(26)
	otp := gotp.NewDefaultTOTP(secret)
	qrcode = otp.ProvisioningUri(username, issuerName)

	return

}
