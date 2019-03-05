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

	model "./model"
	util "./util"
)

var scanner = bufio.NewScanner(os.Stdin)

var apiRoot = "https://musicbrainz.org/ws/2"
var artistRequest = func(artist string) string { return apiRoot + "/artist/?query=" + url.QueryEscape(artist) + "&fmt=json" }

func main() {
	util.Logo()

	artist, album := getInput()

	fmt.Println(artist)
	fmt.Println(album)

	fmt.Println(artistRequest(artist))
	bestArtist, bestArtistID, bestArtistScore := getBestArtist(artist)

	fmt.Printf("\nBest artist: %s (%s) (Score: %d)\n", bestArtist, bestArtistID, bestArtistScore)
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

func getBestArtist(artist string) (name string, mbid string, score int) {
	fmt.Printf("\nSearching for artists matching '%s'...\n\n", artist)

	j, err := httpGet(artistRequest(artist))

	if err != nil {
		log.Fatal(err)
	}

	response := model.ArtistResponse{}

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

	return response.Artists[0].DisambiguatedName(), response.Artists[0].ID, response.Artists[0].Score
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
