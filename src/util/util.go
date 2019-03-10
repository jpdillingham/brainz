package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

type KeyValuePair struct {
	Key   string
	Value int
}

func ToSortedKeyValueSlice(kvmap map[string]int) (slice []KeyValuePair) {
	for k, v := range kvmap {
		slice = append(slice, KeyValuePair{k, v})
	}

	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Value > slice[j].Value
	})

	return slice
}

// Logo prints the application logo to stdout.
func Logo(out func(string)) {
	out(fmt.Sprintf("\n"))
	out(fmt.Sprintf("       ,. brainz\n"))
	out(fmt.Sprintf(" (¬º-°)¬        \n"))
	out(fmt.Sprintf("\n"))
}

func PromptForInput(prompt string) string {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print(prompt)

	input := ""

	for scanner.Scan() {
		input = scanner.Text()
		break
	}

	return input
}

func HttpGet(url string) ([]byte, error) {
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

		return bodyBytes, nil
	}

	return nil, fmt.Errorf("MusicBrainz server returned status code %d", resp.StatusCode)
}

func Distance(source string, target string) float64 {
	return levenshtein.RatioForStrings([]rune(strings.ToLower(source)), []rune(strings.ToLower(target)), levenshtein.DefaultOptions)
}
