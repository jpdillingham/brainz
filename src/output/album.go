package output

type Album struct {
	Artist string  `json:"artist"`
	Album  string  `json:"album"`
	MBID   string  `json:"mbid"`
	Score  float64 `json:"score"`
	Tracks []Track `json:"tracks"`
}