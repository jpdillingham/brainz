package output

type Artist struct {
	Artist string  `json:"artist"`
	MBID   string  `json:"mbid"`
	Score  int     `json:"score"`
	Albums []Album `json:"albums"`
}
