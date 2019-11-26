package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type TelegramChat struct {
	gorm.Model
	Username          string `gorm:"unique_index;not null"`
	ChatID            int64  `gorm:"unique"`
	RegistrationToken string `gorm:"unique"`
	TokenExpiration   time.Time
}
