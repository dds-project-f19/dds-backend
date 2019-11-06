package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Auth struct {
	gorm.Model
	Username   string    `gorm:"unique;not null"`
	Claim      int       `gorm:"not null"`
	Token      string    `gorm:"unique;not null"`
	Expiration time.Time `gorm:"not null"`
}
