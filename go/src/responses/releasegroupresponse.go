package responses

import model "../model"

type ReleaseGroupResponse struct {
	ReleaseGroups      []model.ReleaseGroup `json:"release-groups"`
	ReleaseGroupCount  int                  `json:"release-group-count"`
	ReleaseGroupOffset int                  `json:"release-group-offset"`
}
