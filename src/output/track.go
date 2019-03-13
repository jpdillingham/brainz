package output

type Track struct {
	Title           string   `json:"title"`
	MBID            string   `json:"mbid"`
	Score           float64  `json:"score"`
	Disc            int      `json:"disc"`
	Position        int      `json:"position"`
	Number          string   `json:"number"`
	Length          int      `json:"length"`
	AlternateTitles []string `json:"alternate-titles,omitempty"`
}
