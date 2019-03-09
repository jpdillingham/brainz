package model

import "fmt"

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

func (releaseGroup ReleaseGroup) DisambiguatedTitle() string {
	disambiguation := ""

	if releaseGroup.Disambiguation != "" {
		disambiguation = fmt.Sprintf(" (%s)", releaseGroup.Disambiguation)
	}

	return fmt.Sprintf("%s%s", releaseGroup.Title, disambiguation)
}
