package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strings"

	model "./model"
	output "./output"
	responses "./responses"
	util "./util"
)

var apiRoot = "https://musicbrainz.org/ws/2"

var artistRequest = func(artist string) string { return apiRoot + "/artist/?query=" + url.QueryEscape(artist) + "&fmt=json" }
var releaseGroupRequest = func(mbid string, offset int, limit int) string {
	return fmt.Sprintf("%s/release-group?artist=%s&type=album|ep&offset=%d&limit=%d&fmt=json", apiRoot, mbid, offset, limit)
}
var releaseRequest = func(mbid string, offset int, limit int) string {
	return fmt.Sprintf("%s/release?release-group=%s&offset=%d&limit=%d&inc=media+recordings&fmt=json", apiRoot, mbid, offset, limit)
}

var out = func(msg string) {}

func main() {
	artist, album, output := setup()

	out(fmt.Sprintf("Searching for artists matching '%s'...\n\n", artist))
	bestArtist := getBestArtist(artist)
	out(fmt.Sprintf("\nBest artist: %s (%s) (Score: %d%%)\n\n", bestArtist.DisambiguatedName(), bestArtist.ID, bestArtist.Score))

	out(fmt.Sprintf("Searching for release group matching '%s'...\n\n", album))
	bestReleaseGroup := getBestReleaseGroup(album, bestArtist.ID)
	out(fmt.Sprintf("\nBest release group: %s (%s) (Score: %.0f%%)\n\n", bestReleaseGroup.Title, bestReleaseGroup.ID, util.Distance(bestReleaseGroup.Title, album)*100))

	out(fmt.Sprintf("Fetching releases...\n\n"))
	releases := getAllReleases(bestReleaseGroup.ID)

	out(fmt.Sprintf("%d releases fetched. Computing canonical track count(s)...\n\n", len(releases)))
	tracks, err := getCanonicalFormat(releases)

	if err != nil {
		out(fmt.Sprintf("\nInconclusive.  Assuming the earliest release is canonical.\n"))
		releases = []model.Release{releases[0]}
	} else {
		out(fmt.Sprintf("\nCanonical track count(s): %s\n\n", tracks))
		releases = filterNonCanonicalReleases(releases, tracks)
	}

	canonicalRelease := getCanonicalRelease(releases)
	canonicalReleaseFormat, _ := canonicalRelease.MediaInfo()
	canonicalReleaseDate, _ := canonicalRelease.FuzzyDate()
	out(fmt.Sprintf("\nBest probable canonical release: %s (%s, %s, %s) (Score: %.0f%%):\n\n", canonicalRelease.DisambiguatedTitle(), canonicalReleaseDate.Format("2006-01-02"), canonicalReleaseFormat, canonicalRelease.ID, canonicalRelease.Score*100))

	maxLen := 0
	for _, media := range canonicalRelease.Media {
		for _, track := range media.Tracks {
			len := len(track.Title)
			if len > maxLen {
				maxLen = len
			}
		}
	}

	for _, media := range canonicalRelease.Media {
		for _, track := range media.Tracks {
			out(fmt.Sprintf("%5s\t%-*s\t%3.0f%%\t%s\n", track.Number, maxLen, track.Title, track.Score*100, strings.Join(track.AlternateTitles, ", ")))
		}
	}

	if output == "json" {
		printJSON(bestArtist, canonicalRelease)
	}
}

func printJSON(artist model.Artist, release model.Release) {
	a := output.Artist{Artist: artist.Name, MBID: artist.ID, Score: artist.Score}
	album := output.Album{Title: release.Title, MBID: release.ID, Score: release.Score}

	for _, media := range release.Media {
		for _, track := range media.Tracks {
			album.Tracks = append(album.Tracks, output.Track{
				Disc:            media.Position,
				Position:        track.Position,
				Number:          track.Number,
				Title:           track.Title,
				MBID:            track.ID,
				Length:          track.Length,
				Score:           track.Score,
				AlternateTitles: track.AlternateTitles,
			})
		}
	}

	a.Albums = append(a.Albums, album)

	bytes, _ := json.Marshal(a)
	fmt.Print(string(bytes))
}

func setup() (artist string, album string, output string) {
	artistPtr := flag.String("artist", "", "The artist for which to search")
	albumPtr := flag.String("album", "", "The album for which to search")
	outputPtr := flag.String("output", "", "Output mode; json or leave blank for interactive")

	flag.Parse()

	if *outputPtr == "" {
		out = func(msg string) {
			fmt.Print(msg)
		}

		util.Logo(out)
	}

	prompted := false

	if *artistPtr == "" {
		artistInput := util.PromptForInput("Enter artist: ")
		artistPtr = &artistInput
		prompted = true
	}

	if *albumPtr == "" {
		albumInput := util.PromptForInput("Enter album: ")
		albumPtr = &albumInput
		prompted = true
	}

	if prompted {
		fmt.Println()
	}

	return *artistPtr, *albumPtr, *outputPtr
}

func getBestArtist(artist string) (bestArtist model.Artist) {
	j, err := util.HttpGet(artistRequest(artist))

	if err != nil {
		log.Fatal(err)
	}

	response := responses.ArtistResponse{}

	err = json.Unmarshal(j, &response)

	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(response.Artists[:], func(i, j int) bool {
		return response.Artists[i].Score > response.Artists[j].Score
	})

	for index, artist := range response.Artists {
		prefix := "   "

		if index == 0 {
			prefix = "-->"
		}

		if index < 5 {
			out(fmt.Sprintf("%s %3d%%\t%s\n", prefix, artist.Score, artist.DisambiguatedName()))
		} else {
			break
		}
	}

	return response.Artists[0]
}

func getBestReleaseGroup(album string, mbid string) (bestReleaseGroup model.ReleaseGroup) {
	response := responses.ReleaseGroupResponse{}
	var releaseGroups = []model.ReleaseGroup{}

	for true {
		j, err := util.HttpGet(releaseGroupRequest(mbid, len(releaseGroups), 100))

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(j, &response)

		if err != nil {
			log.Fatal(err)
		}

		releaseGroups = append(releaseGroups, response.ReleaseGroups[:]...)

		if len(releaseGroups) >= response.ReleaseGroupCount {
			break
		}
	}

	sort.Slice(releaseGroups[:], func(i, j int) bool {
		return util.Distance(releaseGroups[i].Title, album) > util.Distance(releaseGroups[j].Title, album)
	})

	top5ReleaseGroups := []model.ReleaseGroup{}
	maxTitleLength := 0

	for index, releaseGroup := range releaseGroups {
		if index < 5 {
			len := len(releaseGroup.DisambiguatedTitle())
			if len > maxTitleLength {
				maxTitleLength = len
			}

			top5ReleaseGroups = append(top5ReleaseGroups, releaseGroup)
		} else {
			break
		}
	}

	for index, releaseGroup := range top5ReleaseGroups {
		prefix := "   "

		if index == 0 {
			prefix = "-->"
		}

		out(fmt.Sprintf("%s %3.0f%%\t%-*s\t%s\n", prefix, util.Distance(releaseGroup.Title, album)*100, maxTitleLength, releaseGroup.DisambiguatedTitle(), releaseGroup.Types()))
	}

	return releaseGroups[0]
}

func getAllReleases(mbid string) (releases []model.Release) {
	response := responses.ReleaseResponse{}
	releases = []model.Release{}

	for true {
		//fmt.Println(releaseRequest(mbid, len(releases), 100))
		j, err := util.HttpGet(releaseRequest(mbid, len(releases), 100))

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(j, &response)

		if err != nil {
			log.Fatal(err)
		}

		releases = append(releases, response.Releases[:]...)

		if len(releases) >= response.ReleaseCount {
			break
		}
	}

	// sort by release date asc
	sort.Slice(releases, func(i, j int) bool {
		idate, _ := releases[i].FuzzyDate()
		jdate, _ := releases[j].FuzzyDate()
		return jdate.After(idate)
	})

	return releases
}

func getCanonicalRelease(releases []model.Release) (canonicalRelease model.Release) {
	bestReleaseIndex := 0
	bestScore := 0.0

	for releaseIndex, release := range releases {
		tempScore := 0.0
		tempTracks := 0

		for mediaIndex, media := range release.Media {
			for trackIndex, track := range media.Tracks {
				match := 0.0

				for i := range releases {
					otherTrack := releases[i].Media[mediaIndex].Tracks[trackIndex]
					if otherTrack.Title == track.Title {
						match++
					} else {
						// prevent adding it if it's already added
						exists := false
						for _, alt := range track.AlternateTitles {
							if alt == otherTrack.Title {
								exists = true
							}
						}

						if !exists {
							track.AlternateTitles = append(track.AlternateTitles, otherTrack.Title)
						}
					}
				}

				// fix the track number for multidisc albums
				if len(releases[0].Media) > 1 {
					track.Number = fmt.Sprintf("%d%02d", media.Position, track.Position)
				} else {
					track.Number = fmt.Sprintf("%02d", track.Position)
				}

				// todo: compare length?
				track.Score = match / float64(len(releases))
				tempScore += track.Score
				tempTracks++

				release.Media[mediaIndex].Tracks[trackIndex] = track
			}
		}

		release.Score = tempScore / float64(tempTracks)
		releases[releaseIndex] = release

		if release.Score > bestScore {
			bestScore = release.Score
			bestReleaseIndex = releaseIndex
		}
	}

	return releases[bestReleaseIndex]
}

func filterNonCanonicalReleases(releases []model.Release, tracks string) (canonicalReleases []model.Release) {
	filteredReleases := []model.Release{}

	for _, release := range releases {
		_, releaseTracks := release.MediaInfo()

		if releaseTracks == tracks {
			filteredReleases = append(filteredReleases, release)
		}
	}

	releases = filteredReleases

	out(fmt.Sprintf("Probable canonical releases:\n\n"))

	maxLen := 0
	for _, release := range releases {
		len := len(release.DisambiguatedTitle())
		if len > maxLen {
			maxLen = len
		}
	}

	for _, release := range releases {
		format, _ := release.MediaInfo()
		date, _ := release.FuzzyDate()
		out(fmt.Sprintf("%s\t%-*s\t%s\n", date.Format("2006-01-02"), maxLen, release.DisambiguatedTitle(), format))
	}

	out(fmt.Sprintf(""))

	return releases
}

func getCanonicalFormat(releases []model.Release) (tracks string, err error) {
	mediaCounts := make(map[string]int)
	trackCounts := make(map[string]int)

	for _, release := range releases {
		_, tracks := release.MediaInfo()
		mediaCounts[tracks]++
		trackCounts[tracks]++
	}

	slice := util.ToSortedKeyValueSlice(trackCounts)
	inconclusive := false

	if len(slice) > 1 && slice[0].Value == slice[1].Value {
		inconclusive = true
	}

	for index, kv := range slice {
		prefix := "   "

		if !inconclusive && index == 0 {
			prefix = "-->"
		}

		out(fmt.Sprintf("%s %3dx\t%s\n", prefix, kv.Value, kv.Key))
	}

	if inconclusive {
		return "", errors.New("unable to determine canonical format")
	}

	return slice[0].Key, nil
}
