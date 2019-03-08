package model

import "fmt"

type Release struct {
	PackagingID        string             `json:"packaging-id"`
	Asin               string             `json:"asin"`
	StatusID           string             `json:"status-id"`
	Disambiguation     string             `json:"disambiguation"`
	Date               string             `json:"date"`
	Packaging          string             `json:"packaging"`
	Status             string             `json:"status"`
	ReleaseEvents      []ReleaseEvent     `json:"release-events"`
	CoverArtArchive    CoverArtArchive    `json:"cover-art-archive"`
	TextRepresentation TextRepresentation `json:"text-representation"`
	Quality            string             `json:"quality"`
	Title              string             `json:"title"`
	Country            string             `json:"country"`
	ID                 string             `json:"id"`
	Media              []Media            `json:"media"`
	Barcode            string             `json:"barcode"`
}

func (release Release) DisambiguatedName() string {
	disambiguation := ""

	if release.Disambiguation != "" {
		disambiguation = fmt.Sprintf(" (%s)", release.Disambiguation)
	}

	return fmt.Sprintf("%s%s", release.Title, disambiguation)
}
