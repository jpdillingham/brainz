package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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

	fmt.Println(httpGet(artistRequest(artist)))
}

func getInput() (string, string) {
	artistPtr := flag.String("artist", "", "The artist for which to search")
	albumPtr := flag.String("album", "", "The album for which to search")
	flag.Parse()

	return *artistPtr, *albumPtr
}

func httpGet(url string) (string, error) {
	var client http.Client

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("User-Agent", "brainz/1.0.0 (https://github.com/jpdillingham/brainz)")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		bodyString := string(bodyBytes)
		return bodyString, nil
	}

	return "", fmt.Errorf("MusicBrainz server returned status code %d", resp.StatusCode)
}
