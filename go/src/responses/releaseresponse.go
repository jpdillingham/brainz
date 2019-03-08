package responses

import model "../model"

type ReleaseResponse struct {
	ReleaseOffset int             `json:"release-offset"`
	Releases      []model.Release `json:"releases"`
	ReleaseCount  int             `json:"release-count"`
}
