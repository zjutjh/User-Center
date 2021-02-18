package database

import (
	"auth-server/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func init() {
	cfg := config.Get()
	if cfg.Db.Type != "postgres" {
		log.Fatalln("unsupported database type")
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", cfg.Db.Host, cfg.Db.User,
		cfg.Db.Password, cfg.Db.DbName, cfg.Db.Port)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&User{}, &Scope{})
	if err != nil {
		log.Fatalln(err)
	}
}