package storage

import (
	"errors"

	"gorm.io/gorm"
)

type otpTable interface {
	Set(username string, secret string) error
	Update(username string, secret string) error
	Delete(username string) error
	Get(username string) (password string, err error)
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
func Get(username string) (password string, err error) {
	password, err = otpDatabaseEngine.Get(username)
	if err != nil {
		return password, err
	}
	if username == "" || password == "" {
		return password, errors.New("user or password is empty")
	}
	return password, nil
}

func GetAdminPassword(uname string) (password string) {
	return webAdminDatabseEngine.GetAdminPassword(uname)
}
func SetAdminPassword(password string) {
	webAdminDatabseEngine.SetAdminPassword(password)

}
