package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"

	model "./model"
	responses "./responses"
	util "./util"
)

var scanner = bufio.NewScanner(os.Stdin)

var apiRoot = "https://musicbrainz.org/ws/2"
var artistRequest = func(artist string) string { return apiRoot + "/artist/?query=" + url.QueryEscape(artist) + "&fmt=json" }
var releaseGroupRequest = func(mbid string, offset int, limit int) string {
	return fmt.Sprintf("%s/release-group?artist=%s&offset=%d&limit=%d&fmt=json", apiRoot, mbid, offset, limit)
}
var releaseRequest = func(mbid string, offset int, limit int) string {
	return fmt.Sprintf("%s/release?release-group=%s&offset=%d&limit=%d&inc=media+recordings&fmt=json", apiRoot, mbid, offset, limit)
}

var distance = func(source string, target string) float64 {
	return levenshtein.RatioForStrings([]rune(strings.ToLower(source)), []rune(strings.ToLower(target)), levenshtein.DefaultOptions)
}

func main() {
	util.Logo()

	artist, album := getInput()

	bestArtist := getBestArtist(artist)

	fmt.Printf("\nBest artist: %s (%s) (Score: %d%%)\n\n", bestArtist.DisambiguatedName(), bestArtist.ID, bestArtist.Score)

	bestReleaseGroup := getBestReleaseGroup(album, bestArtist.ID)

	fmt.Printf("\nBest release group: %s (%s) (Score: %.0f%%)\n\n", bestReleaseGroup.Title, bestReleaseGroup.ID, distance(bestReleaseGroup.Title, album)*100)

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
		artistInput := promptForInput("Enter artist: ")
		artistPtr = &artistInput
	}

	if *albumPtr == "" {
		albumInput := promptForInput("Enter album: ")
		albumPtr = &albumInput
	}

	return *artistPtr, *albumPtr
}

func promptForInput(prompt string) string {
	fmt.Print(prompt)

	input := ""

	for scanner.Scan() {
		input = scanner.Text()
		break
	}

	return input
}

func getBestArtist(artist string) (bestArtist model.Artist) {
	fmt.Printf("Searching for artists matching '%s'...\n\n", artist)

	j, err := httpGet(artistRequest(artist))

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
		j, err := httpGet(releaseGroupRequest(mbid, len(releaseGroups), 100))

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
		return distance(releaseGroups[i].Title, album) > distance(releaseGroups[j].Title, album)
	})

	for index, releaseGroup := range releaseGroups {
		prefix := "   "

		if index == 0 {
			prefix = "-->"
		}

		if index < 5 {
			fmt.Printf("%s %3.0f%%\t%s\n", prefix, distance(releaseGroup.Title, album)*100, releaseGroup.DisambiguatedName())
		} else {
			break
		}
	}

	return releaseGroups[0]
}

func getTrackList(mbid string) (tracks []model.Track) {
	fmt.Printf("Selecting best release...\n\n")
	fmt.Printf("\n%s\n", releaseRequest(mbid, 0, 100))

	response := responses.ReleaseResponse{}
	var releases = []model.Release{}

	for true {
		j, err := httpGet(releaseRequest(mbid, len(releases), 100))

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

	for index, release := range releases {
		prefix := "   "

		if index == 0 {
			prefix = "-->"
		}

		if index < 5 {
			fmt.Printf("%s %3.0f%%\t%s\n", prefix, 0, release.DisambiguatedName())
		} else {
			break
		}
	}

	return releases[0].Media[0].Tracks
}

func httpGet(url string) ([]byte, error) {
	var client http.Client

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "brainz/1.0.0 (https://github.com/jpdillingham/brainz)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		//fmt.Println(string(bodyBytes))
		return bodyBytes, nil
	}

	return nil, fmt.Errorf("MusicBrainz server returned status code %d", resp.StatusCode)
}
