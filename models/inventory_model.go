package models

import "github.com/jinzhu/gorm"

type AvailableItem struct {
	gorm.Model
	ItemType string
	Count    int
}

type TakenItem struct {
	gorm.Model
	TakenBy        string
	ItemType       string
	AssignedToSlot string
}
