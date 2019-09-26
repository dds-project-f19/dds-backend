package models

type User struct {
	ModelBase
	Username string `gorm:"type:varchar(30);unique_index" json:"username,omitempty"`
	Name     string `gorm:"type:varchar(30);not null;default:''" json:"name,omitempty"`
	Surname  string `gorm:"type:varchar(30);not null;default:''" json:"surname,omitempty"`
	Phone    string `gorm:"type:varchar(30);not null;default:''" json:"phone,omitempty"`
	Address  string `gorm:"type:varchar(30);not null;default:''" json:"address,omitempty"`
	Password string `gorm:"type:varchar(64);not null;default:''" json:"password,omitempty"`
}
