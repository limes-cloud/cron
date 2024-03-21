package task

import "github.com/limes-cloud/kratosx/types"

type Group struct {
	types.BaseModel
	Name    string `json:"name" gorm:"not null;size:128;comment:分组名称"`
	Keyword string `json:"keyword" gorm:"not null;size:32;comment:分组标识"`
}

type Task struct {
	Name          string `json:"name" gorm:"not null;size:128;comment:任务名称"`
	GroupID       uint32 `json:"group_id" gorm:"not null;comment:分组id"`
	OwnerID       uint32 `json:"owner_id" gorm:"not null;comment:管理员"`
	Spec          string `json:"spec" gorm:"not null;size:32;comment:定时表达式"`
	Type          string `json:"type" gorm:"not null;size:32;comment:任务类型"`
	Content       string `json:"content" gorm:"not null;type:text;comment:任务内容"`
	RetryCount    int    `json:"retry_count" gorm:"not null;default:0;comment:重试次数"`
	RetryWaitTime int    `json:"retry_wait_time" gorm:"not null;default:0;comment:重试等待时间"`
}

type Tasking struct {
	TaskID uint32 `json:"task_id" gorm:"not null;comment:任务id"`
}
