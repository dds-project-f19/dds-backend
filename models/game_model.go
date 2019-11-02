package models

type ToolType struct {
	ModelBase
	ToolId      int    `gorm:"type:integer;unique_index" json:"tool_id,omitempty"`
	Description string `gorm:"type:varchar(30);unique_index" json:"tool_id,omitempty"`
}

type GameState struct {
	ModelBase
	Username string `gorm:"type:varchar(30);unique_index" json:"username,omitempty"`
	GameType string `gorm:"type:varchar(30);" json:"game_state,omitempty"`
}

type Game struct {
}
