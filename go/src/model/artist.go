package model

type Artist struct {
	ID             string   `json:"id"`
	Type           string   `json:"type,omitempty"`
	TypeID         string   `json:"type-id,omitempty"`
	Score          int      `json:"score"`
	Name           string   `json:"name"`
	SortName       string   `json:"sort-name"`
	Country        string   `json:"country,omitempty"`
	Area           Area     `json:"area,omitempty"`
	BeginArea      Area     `json:"begin-area,omitempty"`
	Disambiguation string   `json:"disambiguation,omitempty"`
	LifeSpan       Lifespan `json:"life-span"`
	Tags           []Tag    `json:"tags,omitempty"`
	Aliases        []Alias  `json:"aliases,omitempty"`
	Gender         string   `json:"gender,omitempty"`
}
