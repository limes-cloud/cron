package types

type ListTaskGroupRequest struct {
	App      string  `json:"-"` // 中转参数（app 层根据 X-App-Key 推导设置，不入 proto）
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Order    *string `json:"order"`
	OrderBy  *string `json:"orderBy"`
	Name     *string `json:"name"`
}

type ListTaskRequest struct {
	App      string  `json:"-"` // 中转参数（app 层根据 X-App-Key 推导设置，不入 proto）
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	GroupId  *uint32 `json:"groupId"`
	Name     *string `json:"name"`
	Tag      *string `json:"tag"`
	Status   *bool   `json:"status"`
}

type ExecTaskRequest struct {
	Id            uint32
	Uuid          string
	IP            string
	Ak            string
	Sk            string
	ExecType      string
	ExecValue     string
	ExpectCode    uint32
	RetryCount    uint32
	RetryWaitTime uint32
	MaxExecTime   uint32
}

type CancelExecRequest struct {
	Uuid string
	IP   string
	Ak   string
	Sk   string
}

// ListExecutionRequest 执行历史查询请求
type ListExecutionRequest struct {
	App         string  `json:"-"` // 中转参数（app 层根据 X-App-Key 推导设置，不入 proto）
	Page        uint32  `json:"page"`
	PageSize    uint32  `json:"pageSize"`
	TaskId      uint32  `json:"taskId"`
	Status      *string `json:"status"`      // 可选：按状态过滤（pending/running/success/fail/cancel/skip）
	TriggerType *string `json:"triggerType"` // 可选：按触发方式过滤（cron/manual/misfire）
}
