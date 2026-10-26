package types

type ListLogRequest struct {
	App      string `json:"-"` // 中转参数（app 层根据 X-App-Key 推导设置，不入 proto）
	Page     uint32 `json:"page"`
	PageSize uint32 `json:"pageSize"`
	TaskId   uint32 `json:"taskId"`
}

type ExecTaskLog struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    int64  `json:"time"`
}
