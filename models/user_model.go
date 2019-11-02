package models

type User struct {
	ModelBase
	Username string `form:"username" binding:"required" gorm:"type:varchar(30);unique_index" json:"username,omitempty"`
	Name     string `form:"name" binding:"required" gorm:"type:varchar(30);not null;default:''" json:"name,omitempty"`
	Surname  string `form:"surname" binding:"required" gorm:"type:varchar(30);not null;default:''" json:"surname,omitempty"`
	Phone    string `form:"phone" binding:"required" gorm:"type:varchar(30);not null;default:''" json:"phone,omitempty"`
	Address  string `form:"address" binding:"required" gorm:"type:varchar(30);not null;default:''" json:"address,omitempty"`
	Password string `form:"password" binding:"required" gorm:"type:varchar(64);not null;default:''" json:"password,omitempty"`
}

type Worker struct {
	User
	GameType string   `form:"game_type" binding:"required" gorm:"type:varchar(64);not null;default:''" json:"game_type,omitempty"`
	Cells    []string `form:"cells" gorm:"type:varchar(64);not null" json:"cells,omitempty"`
}

type InventoryItem struct {
	ModelBase
	GameType string `form:"game_type" binding:"required" gorm:"type:varchar(64);not null;default:''" json:"password,omitempty"`
	ItemId   string `form:"item_id" binding:"required" gorm:"type:varchar(64);not null;default:''" json:"password,omitempty"`
	Count    int    `form:"password" binding:"required" gorm:"type:varchar(64);not null;default:''" json:"password,omitempty"`
}
