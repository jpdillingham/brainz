package model

type Recording struct {
	Disambiguation string `json:"disambiguation"`
	ID             string `json:"id"`
	Title          string `json:"title"`
	Video          bool   `json:"video"`
	Length         int    `json:"length"`
}
