package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
	Score              float64            `json:"score"`
}

func (release Release) DisambiguatedTitle() string {
	disambiguation := ""

	if release.Disambiguation != "" {
		disambiguation = fmt.Sprintf(" (%s)", release.Disambiguation)
	}

	return fmt.Sprintf("%s%s", release.Title, disambiguation)
}

func (release Release) FuzzyDate() (date time.Time, err error) {
	format := "2006-01-02"

	if len(release.Date) == 0 {
		return time.Now(), nil
	}

	if len(release.Date) == 2 {
		y, _ := strconv.Atoi(release.Date)
		if y < 30 {
			return time.Parse(format, fmt.Sprintf("20%s-12-31", release.Date))
		}

		return time.Parse(format, fmt.Sprintf("19%s-12-31", release.Date))
	}

	if len(release.Date) == 4 {
		return time.Parse(format, fmt.Sprintf("%s-12-31", release.Date))
	}

	return time.Parse(format, release.Date)
}

func (release Release) MediaInfo() (formats string, tracks string) {
	mediaArray := []string{}
	trackArray := []string{}

	for _, media := range release.Media {
		mediaArray = append(mediaArray, media.Format)
		trackArray = append(trackArray, strconv.Itoa(media.TrackCount))
	}

	if len(release.Media) == 1 {
		return strings.Join(mediaArray, " + "), strings.Join(trackArray, " + ")
	}

	// if there's more than one type of media, iterate over the array and build a
	// string in the format Nx<Media> where N is the number of contiguously repeated
	// media of that type
	// e.g. CD CD CD DVD DVD = 3xCD + 2xDVD
	dedupedMediaArray := []string{}
	currentstr := release.Media[0].Format
	currentCount := 0

	for _, media := range release.Media {
		if media.Format != currentstr {
			dedupedMediaArray = append(dedupedMediaArray, fmt.Sprintf("%dx%s", currentCount, currentstr))
			currentstr = media.Format
			currentCount = 1
		} else {
			currentCount++
		}
	}

	dedupedMediaArray = append(dedupedMediaArray, fmt.Sprintf("%dx%s", currentCount, currentstr))

	return strings.Join(dedupedMediaArray, " + "), strings.Join(trackArray, " + ")
}
