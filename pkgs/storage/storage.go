package storage

import "gorm.io/gorm"

type otpTable interface {
	Set(username string, secret string) error
	Update(username string, secret string) error
	Delete(username string) error
	Get(username string) (password string)
	Connect() *gorm.DB
	Migrate()
}

type webAdminTable interface {
	GetAdminPassword(uname string) (password string)
	SetAdminPassword(password string)
	Connection() *gorm.DB
	Migrate()
}

var (
	Dsn                   string
	otpDatabaseEngine     otpTable
	webAdminDatabseEngine webAdminTable
)

func SetOtpStorageType(database string) otpTable {
	return postgresOtp{}
}

func SetWebAminStorageType(database string) webAdminTable {
	return postgresWebAdmin{}
}

func Initialize() {
	otpDatabaseEngine = SetOtpStorageType("postgres")
	otpDatabaseEngine.Migrate()
	webAdminDatabseEngine = SetWebAminStorageType("postgres")
	webAdminDatabseEngine.Migrate()

}

func Set(username string, secret string) error {
	return otpDatabaseEngine.Set(username, secret)
}
func Update(username string, secret string) error {
	return otpDatabaseEngine.Update(username, secret)
}
func Delete(username string) error {
	return otpDatabaseEngine.Delete(username)
}
func Get(username string) (password string) {
	return otpDatabaseEngine.Get(username)
}

func GetAdminPassword(uname string) (password string) {
	return webAdminDatabseEngine.GetAdminPassword(uname)
}
func SetAdminPassword(password string) {
	webAdminDatabseEngine.SetAdminPassword(password)

}
