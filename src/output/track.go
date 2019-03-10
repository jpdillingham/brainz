package output

type Track struct {
	Number string  `json:"number"`
	Title  string  `json:"title"`
	Length int     `json:"length"`
	Score  float32 `json:"score"`
}
