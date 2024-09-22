package types

type ListTaskGroupRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Order    *string `json:"order"`
	OrderBy  *string `json:"orderBy"`
	Name     *string `json:"name"`
}

type ListTaskRequest struct {
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
