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

func (release Release) DisambiguatedTitle() string {
	disambiguation := ""

	if release.Disambiguation != "" {
		disambiguation = fmt.Sprintf(" (%s)", release.Disambiguation)
	}

	return fmt.Sprintf("%s%s", release.Title, disambiguation)
}

func (release Release) MediaInfo() (string, string) {
	mediastr := ""
	trackstr := ""

	for index, media := range release.Media {
		sep := ""

		if index > 0 {
			sep = " + "
		}

		mediastr = fmt.Sprintf("%s%s%s", mediastr, sep, media.Format)
		trackstr = fmt.Sprintf("%s%s%d", trackstr, sep, media.TrackCount)
	}

	return mediastr, trackstr
}
