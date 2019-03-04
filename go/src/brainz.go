package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	util "./util"
)

var stdin = bufio.NewReader(os.Stdin)

func main() {
	util.Logo()

	artist, album := getInput()

	fmt.Println(artist)
	fmt.Println(album)
}

func getInput() (string, string) {
	artistPtr := flag.String("artist", "", "The artist for which to search")
	albumPtr := flag.String("album", "", "The album for which to search")
	flag.Parse()

	return *artistPtr, *albumPtr
}
