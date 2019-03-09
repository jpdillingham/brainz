package model

type Alias struct {
	SortName  string `json:"sort-name"`
	Name      string `json:"name"`
	Locale    string `json:"locale"`
	Type      string `json:"type"`
	Primary   bool   `json:"primary"`
	BeginDate string `json:"begin-date"`
	EndDate   string `json:"end-date"`
}
