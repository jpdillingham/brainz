package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"

	model "./model"
	util "./util"
)

var stdin = bufio.NewReader(os.Stdin)

var apiRoot = "https://musicbrainz.org/ws/2"
var artistRequest = func(artist string) string { return apiRoot + "/artist/?query=" + artist + "&fmt=json" }

func main() {
	util.Logo()

	artist, album := getInput()

	fmt.Println(artist)
	fmt.Println(album)

	fmt.Println(artistRequest(artist))
	bestArtist, bestArtistId, bestArtistScore := getBestArtist(artist)

	fmt.Printf("Best artist: %s (%s) (Score: %d)\n", bestArtist, bestArtistId, bestArtistScore)
}

func getBestArtist(search string) (name string, mbid string, score int) {

	j, _ := httpGet(artistRequest(search))
	response := model.ArtistResponse{}

	err := json.Unmarshal(j, &response)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(fmt.Sprintf("Count: %d", response.Count))

	sort.Slice(response.Artists[:], func(i, j int) bool {
		return response.Artists[i].Score > response.Artists[j].Score
	})

	for _, artist := range response.Artists {
		fmt.Printf("%d%%\t%s\n", artist.Score, artist.DisambiguatedName())
	}

	return response.Artists[0].DisambiguatedName(), response.Artists[0].ID, response.Artists[0].Score
}

func getDisambiguatedArtistName(artist model.Artist) string {
	disambiguation := ""

	if artist.Disambiguation != "" {
		disambiguation = fmt.Sprintf(" (%s)", artist.Disambiguation)
	}

	return fmt.Sprintf("%s%s", artist.Name, disambiguation)
}

func getInput() (string, string) {
	artistPtr := flag.String("artist", "", "The artist for which to search")
	albumPtr := flag.String("album", "", "The album for which to search")
	flag.Parse()

	return *artistPtr, *albumPtr
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

		//bodyString := string(bodyBytes)
		return bodyBytes, nil
	}

	return nil, fmt.Errorf("MusicBrainz server returned status code %d", resp.StatusCode)
}
