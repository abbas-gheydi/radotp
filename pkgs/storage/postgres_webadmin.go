package storage

import (
	"errors"
	"log"
	"sync"

	"gorm.io/driver/postgres"
	_ "gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

var db_web *gorm.DB
var once_web sync.Once

type radotpWebAdmin struct {
	Uname    string
	Role     string
	Password string
}

type postgresWebAdmin struct{}

func (p postgresWebAdmin) GetAdminPassword(uname string) (password string) {
	var resualt radotpWebAdmin

	db_web.First(&resualt, "Uname = ?", uname)
	return resualt.Password

}

func (p postgresWebAdmin) SetAdminPassword(password string) {
	hash_password := ShaGenerator(password)
	tx := db_web.Model(&radotpWebAdmin{}).Where("uname = ?", "admin").Update("password", hash_password)
	if tx.Error != nil {
		log.Println("*****db.go", tx.Error)
	} else {
		log.Println("admin password is updated")
	}
	if tx.RowsAffected != 1 {
		//create admin user with default password
		tx.AddError(errors.New("user not found "))
		db_web.Create(&radotpWebAdmin{Uname: "admin", Role: "admin", Password: hash_password})
		log.Print(" create admin/admin username for web access")

	}

}

func (p postgresWebAdmin) Migrate() {

	db_web := p.Connection()
	db_web.AutoMigrate(&radotpWebAdmin{})
	if p.GetAdminPassword("admin") == "" {
		p.SetAdminPassword("admin")

	}

}

func (p postgresWebAdmin) Connection() *gorm.DB {
	var err error

	once_web.Do(func() {
		db_web, err = gorm.Open(postgres.Open(Dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}
	})
	return db_web

}
