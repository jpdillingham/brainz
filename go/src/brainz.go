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

func main() {
	util.Logo()

	artist, album := getInput()

	bestArtist := getBestArtist(artist)

	fmt.Printf("\nBest artist: %s (%s) (Score: %d%%)\n\n", bestArtist.DisambiguatedName(), bestArtist.ID, bestArtist.Score)

	bestReleaseGroup := getBestReleaseGroup(album, bestArtist.ID)

	fmt.Printf("\nBest release group: %s (%s) (Score: %.0f%%)\n\n", bestReleaseGroup.Title, bestReleaseGroup.ID, util.Distance(bestReleaseGroup.Title, album)*100)

	trackList := getCanonicalTrackList(bestReleaseGroup.ID)

	// todo: final output
	for _, track := range trackList {
		fmt.Printf("   %3.0f%%\t%s\t%s\n", track.Score*100, track.Number, track.Title)
	}
}

func getInput() (string, string) {
	artistPtr := flag.String("artist", "", "The artist for which to search")
	albumPtr := flag.String("album", "", "The album for which to search")
	flag.Parse()

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

	return *artistPtr, *albumPtr
}

func getBestArtist(artist string) (bestArtist model.Artist) {
	fmt.Printf("Searching for artists matching '%s'...\n\n", artist)

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
			fmt.Printf("%s %3d%%\t%s\n", prefix, artist.Score, artist.DisambiguatedName())
		} else {
			break
		}
	}

	return response.Artists[0]
}

func getBestReleaseGroup(album string, mbid string) (bestReleaseGroup model.ReleaseGroup) {
	fmt.Printf("Searching for release group matching '%s'...\n\n", album)

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

		fmt.Printf("%s %3.0f%%\t%-*s\t%s\n", prefix, util.Distance(releaseGroup.Title, album)*100, maxTitleLength, releaseGroup.DisambiguatedTitle(), releaseGroup.Types())
	}

	return releaseGroups[0]
}

func getCanonicalTrackList(mbid string) (canonicalTracks []model.Track) {
	fmt.Printf("Determining canonical track listing...\n\n")

	response := responses.ReleaseResponse{}
	var releases = []model.Release{}

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

	releases = filterNonCanonicalReleases(releases)

	fmt.Printf("\nProbable canonical track listing:\n\n")

	var tracks = []model.Track{}
	bestScore := 0.0

	for _, release := range releases {
		tempScore := 0.0
		var tempTracks = []model.Track{}

		for mediaIndex, media := range release.Media {
			for trackIndex, track := range media.Tracks {
				match := 0.0

				for i := range releases {
					otherTrack := releases[i].Media[mediaIndex].Tracks[trackIndex]
					if otherTrack.Title == track.Title {
						match++
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
				tempTracks = append(tempTracks, track)
				tempScore += track.Score
			}
		}

		releaseScore := tempScore / float64(len(tempTracks))

		if releaseScore > bestScore {
			bestScore = releaseScore
			tracks = tempTracks
		}
	}

	return tracks
}

func filterNonCanonicalReleases(releases []model.Release) (canonicalReleases []model.Release) {
	canonicalFormat, canonicalTracks, formatErr := getCanonicalFormat(releases)

	if formatErr == nil {
		// filter non-canonical releases
		fmt.Printf("\nCanonical format: %s, canonical track count: %s\n\n", canonicalFormat, canonicalTracks)

		filteredReleases := []model.Release{}

		for _, release := range releases {
			releaseFormat, releaseTracks := release.MediaInfo()

			if releaseFormat == canonicalFormat && releaseTracks == canonicalTracks {
				filteredReleases = append(filteredReleases, release)
			}
		}

		releases = filteredReleases

		fmt.Printf("Probable canonical releases:\n\n")
	} else {
		fmt.Printf("\nCanonical format is inconclusive.\n\n")
		fmt.Printf("Unfiltered releases:\n\n")
	}

	sort.Slice(releases, func(i, j int) bool {
		idate, _ := releases[i].FuzzyDate()
		jdate, _ := releases[j].FuzzyDate()
		return jdate.After(idate)
	})

	maxLen := 0
	for _, release := range releases {
		len := len(release.DisambiguatedTitle())
		if len > maxLen {
			maxLen = len
		}
	}

	for _, release := range releases {
		date, _ := release.FuzzyDate()
		media, _ := release.MediaInfo()
		fmt.Printf("   %s\t%-*s\t%s\n", date.Format("2006-01-02"), maxLen, release.DisambiguatedTitle(), media)
	}

	return releases
}

func getCanonicalFormat(releases []model.Release) (format string, tracks string, err error) {
	fmt.Printf("Compiling formats..\n\n")

	// count the number of releases for each combination of media and track counts
	mediaCounts := make(map[string]int)

	for _, release := range releases {
		media, tracks := release.MediaInfo()
		mediaCounts[media+":"+tracks]++
	}

	// sort the map by descending number of occurances
	type KeyValuePair struct {
		Key   string
		Value int
	}

	var mediaCountSlice []KeyValuePair

	for k, v := range mediaCounts {
		mediaCountSlice = append(mediaCountSlice, KeyValuePair{k, v})
	}

	sort.Slice(mediaCountSlice, func(i, j int) bool {
		return mediaCountSlice[i].Value > mediaCountSlice[j].Value
	})

	// determine the max format string length for spacing purposes
	// initialize to the length of "format"
	maxLen := 6
	for _, format := range mediaCountSlice {
		len := len(strings.Split(format.Key, ":")[0])
		if len > maxLen {
			maxLen = len
		}
	}

	// check to see if the count of the top format matches the next, indicating
	// an inconclusive match
	inconclusive := false

	if len(mediaCountSlice) > 1 && mediaCountSlice[0].Value == mediaCountSlice[1].Value {
		inconclusive = true
	}

	for index, kv := range mediaCountSlice {
		prefix := "   "

		if !inconclusive && index == 0 {
			prefix = "-->"
		}

		parts := strings.Split(kv.Key, ":")

		fmt.Printf("%s %3dx\t%-*s\t%s\n", prefix, kv.Value, maxLen, parts[0], parts[1])
	}

	if inconclusive {
		return "", "", errors.New("unable to determine canonical format")
	}

	bestFormat := mediaCountSlice[0].Key
	bestFormatParts := strings.Split(bestFormat, ":")

	return bestFormatParts[0], bestFormatParts[1], nil
}
