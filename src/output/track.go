package output

type Track struct {
	Disc            int      `json:"disc"`
	Position        int      `json:"position"`
	Number          string   `json:"number"`
	Title           string   `json:"title"`
	Length          int      `json:"length"`
	Score           float64  `json:"score"`
	AlternateTitles []string `json:"alternate-titles,omitempty"`
}
