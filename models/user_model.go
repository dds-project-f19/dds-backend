package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Name     string
	Surname  string
	Phone    string `gorm:"unique"`
	Address  string
	Claim    int `gorm:"not null;default:1"`
}

func (u *User) ToMap() map[string]interface{} {
	result := make(map[string]interface{})
	result["username"] = u.Username
	result["name"] = u.Name
	result["surname"] = u.Surname
	result["phone"] = u.Phone
	result["address"] = u.Address
	return result
}
