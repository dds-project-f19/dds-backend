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
	Claim    int `gorm:"not null"`
}

//type Worker struct {
//	User
//	GameType string   `form:"game_type" binding:"required" gorm:"type:varchar(64);not null;default:''" json:"game_type,omitempty"`
//	Cells    []string `form:"cells" gorm:"type:varchar(64);not null" json:"cells,omitempty"`
//}
//
//type InventoryItem struct {
//	gorm.Model
//	GameType string `form:"game_type" binding:"required" gorm:"type:varchar(64);not null;default:''" json:"password,omitempty"`
//	ItemId   string `form:"item_id" binding:"required" gorm:"type:varchar(64);not null;default:''" json:"password,omitempty"`
//	Count    int    `form:"password" binding:"required" gorm:"type:varchar(64);not null;default:''" json:"password,omitempty"`
//}
