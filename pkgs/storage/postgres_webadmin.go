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

type postgresWebAdmin struct{}

type webadmin struct {
	Uname    string
	Password string
}

func (p postgresWebAdmin) GetAdminPassword(uname string) (password string) {
	var resualt webadmin

	db_web.First(&resualt, "Uname = ?", uname)
	return resualt.Password

}

func (p postgresWebAdmin) SetAdminPassword(password string) {
	hash_password := ShaGenerator(password)

	tx := db_web.Model(&webadmin{}).Where("uname = ?", "admin").Update("password", hash_password)
	if tx.Error != nil {
		log.Println("*****db.go", tx.Error)
	} else {
		log.Println("admin password is updated")
	}
	if tx.RowsAffected != 1 {
		//create admin user with default password
		tx.AddError(errors.New("user not found "))
		log.Print(" create admin/admin username for web access")
		db_web.Create(&webadmin{Uname: "admin", Password: hash_password})
	}

}

func (p postgresWebAdmin) Migrate() {

	db_web := p.Connection()
	db_web.AutoMigrate(&webadmin{})
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
