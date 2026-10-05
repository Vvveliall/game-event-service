package model

type Event struct {
	PlayerID int    `json:"player_id"`
	Type     string `json:"type"`
	Payload  string `json:"payload"`
}
