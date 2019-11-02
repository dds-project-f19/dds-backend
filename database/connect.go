package database

import (
	"dds-backend/config"
	"dds-backend/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

func InitDB(dbConfig config.DBConfig, generalConfig config.GeneralConfig) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dbConfig.GetDSN())

	if err == nil {
		db.DB().SetMaxIdleConns(generalConfig.MaxIdleConn)
		DB = db
		db.AutoMigrate(&models.User{})
		return db, err
	}
	return nil, err
}
