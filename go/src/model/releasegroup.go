package model

type ReleaseGroup struct {
	SecondaryTypeIDs []string `json:"secondary-type-ids"`
	Disambiguation   string   `json:"disambiguation"`
	FirstReleaseDate string   `json:"first-release-date"`
	PrimaryTypeID    string   `json:"primary-type-id"`
	PrimaryType      string   `json:"primary-type"`
	ID               string   `json:"id"`
	Title            string   `json:"title"`
	SecondaryTypes   []string `json:"secondary-types"`
}
