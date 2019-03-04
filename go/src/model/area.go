package model

type Area struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	TypeID   string   `json:"type-id"`
	Name     string   `json:"name"`
	SortName string   `json:"sort-name"`
	LifeSpan Lifespan `json:"life-span"`
}
