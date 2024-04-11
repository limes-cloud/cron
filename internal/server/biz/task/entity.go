package task

import ktypes "github.com/limes-cloud/kratosx/types"

type TaskGroup struct {
	ktypes.BaseModel
	Name        string
	Description string `json:"description"`
}

type Task struct {
	ktypes.BaseModel
	GroupId       uint32  `json:"group_id"`
	Name          string  `json:"name"`
	Tag           string  `json:"tag"`
	Spec          string  `json:"spec"`
	SelectType    string  `json:"select_type"`
	SelectValue   string  `json:"select_value"`
	ExecType      string  `json:"exec_type"`
	WorkerGroupId *uint32 `json:"worker_group_id"`
	WorkerId      *uint32 `json:"worker_id"`
	ExpectCode    uint32  `json:"expect_code"`
	RetryCount    uint32  `json:"retry_count"`
	RetryWaitTime uint32  `json:"retry_wait_time"`
	MaxExecTime   uint32  `json:"max_exec_time"`
	Status        string  `json:"status"`
	Description   string  `json:"description"`
	// 后续新增告警
}
