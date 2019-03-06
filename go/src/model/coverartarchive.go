package model

type CoverArtArchive struct {
	Artwork  bool `json:"artwork"`
	Front    bool `json:"front"`
	Count    int  `json:"count"`
	Back     bool `json:"back"`
	Darkened bool `json:"darkened"`
}
