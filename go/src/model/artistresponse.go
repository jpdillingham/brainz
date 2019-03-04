package model

import "time"

// ArtistResponse encapsulates the response received from an artist search.
type ArtistResponse struct {
	Created time.Time `json:"created"`
	Count   int       `json:"count"`
	Offset  int       `json:"offset"`
	Artists []Artist  `json:"artists"`
}
