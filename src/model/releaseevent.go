package model

type ReleaseEvent struct {
	Date string `json:"date"`
	Area Area   `json:"area"`
}
