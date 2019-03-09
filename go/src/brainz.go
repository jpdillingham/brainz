package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"sort"

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

	trackList := getTrackList(bestReleaseGroup.ID)

	for _, track := range trackList {
		fmt.Printf("%s. %s\n", track.Number, track.Title)
	}
}

func getInput() (string, string) {
	artistPtr := flag.String("artist", "", "The artist for which to search")
	albumPtr := flag.String("album", "", "The album for which to search")
	flag.Parse()

	if *artistPtr == "" {
		artistInput := util.PromptForInput("Enter artist: ")
		artistPtr = &artistInput
	}

	if *albumPtr == "" {
		albumInput := util.PromptForInput("Enter album: ")
		albumPtr = &albumInput
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

func getTrackList(mbid string) (tracks []model.Track) {
	fmt.Printf("Fetching releases...\n\n")

	response := responses.ReleaseResponse{}
	var releases = []model.Release{}

	for true {
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

	// count the number of releases for each combination of media and track counts
	mediaCounts := make(map[string]int)

	for _, release := range releases {
		media, tracks := release.MediaInfo()
		mediaCounts[media+", "+tracks]++
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

	// display the found formats
	for index, kv := range mediaCountSlice {
		prefix := "   "

		if index == 0 {
			prefix = "-->"
		}

		fmt.Printf("%s %s %d\n", prefix, kv.Key, kv.Value)
	}

	return releases[0].Media[0].Tracks
}
