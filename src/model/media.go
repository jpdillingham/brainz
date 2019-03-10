package model

type Media struct {
	FormatID    string  `json:"format-id"`
	TrackCount  int     `json:"track-count"`
	Title       string  `json:"title"`
	Position    int     `json:"position"`
	Format      string  `json:"format"`
	Tracks      []Track `json:"tracks"`
	TrackOffset int     `json:"track-offset"`
}
