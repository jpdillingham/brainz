package output

type Album struct {
	Title  string  `json:"title"`
	MBID   string  `json:"mbid"`
	Score  float64 `json:"score"`
	Tracks []Track `json:"tracks"`
}
