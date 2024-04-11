package worker

import ktypes "github.com/limes-cloud/kratosx/types"

type Worker struct {
	ktypes.BaseModel
	Name        string `json:"name"`
	IP          string `json:"ip"`
	Tag         string `json:"tag"`
	Status      string `json:"status"`
	StopDesc    string `json:"stop_desc"`
	Description string `json:"description"`
}
