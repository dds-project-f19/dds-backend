package database

import (
	"dds-backend/config"
	"dds-backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	conf := config.Get()
	db, err := gorm.Open("mysql", conf.DSN)

	if err == nil {
		db.DB().SetMaxIdleConns(conf.MaxIdleConn)
		DB = db
		db.AutoMigrate(&models.User{})
		return db, err
	}
	return nil, err
}
