package storage

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

const (
	User_not_found string = "user not found"
)

var (
	MaxOpenConns,
	MaxIdleConns,
	ConnMaxLifetimeInMiuntes int
	db_otp *gorm.DB
	once   sync.Once
)

type otps struct {
	ID       uint   `gorm:"primarykey"`
	Username string `gorm:"index:idx_users,type:hash;unique"`
	Secret   string
}

type postgresOtp struct{}

func (p postgresOtp) Set(username string, secret string) error {
	username = strings.ToLower(username)

	if strings.Contains(username, "@") || strings.Contains(username, "\\") {
		splitChar := "@"
		if strings.Contains(username, "\\") {
			splitChar = "\\"
		}
		return fmt.Errorf("username is not valid. please insert a username without %v", splitChar)

	}

	otpUser := otps{
		Username: username,
		Secret:   aesEncrypt(secret, generateEncKey(username)),
	}

	tx := db_otp.Create(&otpUser)
	if tx.Error != nil {
		log.Println("*****db.go", tx.Error)
	} else {
		log.Println("otp code for user ", username, " saved to db")
	}
	return tx.Error
}

func (p postgresOtp) Update(username string, secret string) error {
	username = strings.ToLower(username)

	otpUser := otps{
		Username: username,
		Secret:   aesEncrypt(secret, generateEncKey(username)),
	}

	tx := db_otp.Model(&otpUser).Where("username = ?", username).Update("secret", otpUser.Secret)
	if tx.Error != nil {
		log.Println("*****db.go", tx.Error)
	}
	if tx.RowsAffected != 1 {
		tx.AddError(errors.New(User_not_found))
		log.Println(username, " not found")

	} else {
		log.Println("otp code for user ", username, " saved to db")
	}

	return tx.Error
}

func (p postgresOtp) Delete(username string) error {
	username = strings.ToLower(username)

	otpUser := otps{Username: username}
	tx := db_otp.Model(&otpUser).Where("username = ?", username).Delete(otpUser)
	if tx.Error != nil {
		log.Println("*****db.go", tx.Error)
	} else {
		log.Println("otp code for user ", username, " removed from db")
	}
	if tx.RowsAffected != 1 {
		tx.AddError(errors.New(User_not_found))
	}

	return tx.Error

}

func (p postgresOtp) Get(username string) (password string, err error) {
	username = strings.ToLower(username)
	winNTSplitChar := "\\"
	if strings.Contains(username, winNTSplitChar) && strings.Split(username, winNTSplitChar)[1] != "" {
		username = strings.Split(username, winNTSplitChar)[1]
	}
	otpUser := otps{Username: username}

	tx := db_otp.First(&otpUser, "Username = ?", username)
	if tx.Error != nil && strings.Contains(tx.Error.Error(), "record not found") {
		if strings.Contains(username, "@") {
			splitChar := "@"
			username = strings.Split(username, splitChar)[0]
			tx = db_otp.First(&otpUser, "Username = ?", username)
		}

	}
	if tx.Error != nil {

		return "", tx.Error
	}

	if otpUser.Secret != "" {

		password = aesDecrypt(otpUser.Secret, generateEncKey(username))
	}
	return password, nil

}

func (p postgresOtp) Connect() *gorm.DB {
	var err error
	once.Do(func() {

		db_otp, err = gorm.Open(postgres.Open(Dsn), &gorm.Config{PrepareStmt: true})
		if err != nil {
			panic("failed to connect database")
		}

	})

	sqlDB, err := db_otp.DB()
	sqlDB.SetMaxOpenConns(MaxOpenConns)
	sqlDB.SetMaxIdleConns(MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(ConnMaxLifetimeInMiuntes) * time.Minute)

	return db_otp
}

func (p postgresOtp) Migrate() {
	db := p.Connect()
	db.AutoMigrate(&otps{})

}
