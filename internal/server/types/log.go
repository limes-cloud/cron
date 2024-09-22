package types

type ListLogRequest struct {
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"pageSize"`
	TaskId   uint32 `json:"taskId"`
}

type ExecTaskLog struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}
