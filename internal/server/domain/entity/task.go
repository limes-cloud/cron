package entity

import "github.com/limes-cloud/kratosx/types"

type TaskGroup struct {
	Name        string  `json:"name" gorm:"column:name"`
	Description *string `json:"description" gorm:"column:description"`
	types.BaseModel
}

type Task struct {
	GroupId       uint32       `json:"groupId" gorm:"column:group_id"`
	Name          string       `json:"name" gorm:"column:name"`
	Tag           string       `json:"tag" gorm:"column:tag"`
	Spec          string       `json:"spec" gorm:"column:spec"`
	Status        *bool        `json:"status" gorm:"column:status"`
	WorkerType    string       `json:"workerType" gorm:"column:worker_type"`
	WorkerGroupId *uint32      `json:"workerGroupId" gorm:"column:worker_group_id"`
	WorkerId      *uint32      `json:"workerId" gorm:"column:worker_id"`
	ExecType      string       `json:"execType" gorm:"column:exec_type"`
	ExecValue     string       `json:"execValue" gorm:"column:exec_value"`
	ExpectCode    uint32       `json:"expectCode" gorm:"column:expect_code"`
	RetryCount    uint32       `json:"retryCount" gorm:"column:retry_count"`
	RetryWaitTime uint32       `json:"retryWaitTime" gorm:"column:retry_wait_time"`
	MaxExecTime   uint32       `json:"maxExecTime" gorm:"column:max_exec_time"`
	Version       string       `json:"version" gorm:"column:version"`
	Description   *string      `json:"description" gorm:"column:description"`
	Start         *uint32      `json:"start" gorm:"column:start"`
	End           *uint32      `json:"end" gorm:"column:end"`
	Group         *TaskGroup   `json:"group"`
	Worker        *Worker      `json:"worker"`
	WorkerGroup   *WorkerGroup `json:"workerGroup"`
	types.BaseModel
}
