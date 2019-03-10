package model

type Track struct {
	Position        int       `json:"position"`
	ID              string    `json:"id"`
	Length          int       `json:"length"`
	Title           string    `json:"title"`
	Number          string    `json:"number"`
	Recording       Recording `json:"recording"`
	Score           float64   `json:"score"`
	AlternateTitles []string  `json:"alternate-titles"`
}
