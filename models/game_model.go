package models

import "github.com/jinzhu/gorm"

type ToolType struct {
	gorm.Model
	ToolId      int
	Description string
}

type GameState struct {
	gorm.Model
	Username string
	GameType string
}

type Game struct {
}
