package log

import ktypes "github.com/limes-cloud/kratosx/types"

type Log struct {
	ktypes.BaseModel
	WorkerId       uint32 `json:"worker_id"`
	WorkerSnapshot string `json:"worker_snapshot"`
	TaskId         string `json:"task_id"`
	TaskSnapshot   string `json:"task_snapshot"`
	Start          uint32 `json:"start"`
	End            uint32 `json:"end"`
	ExecCount      uint32 `json:"exec_count"`
	Content        string `json:"content"`
	Status         bool   `json:"status"`
}
