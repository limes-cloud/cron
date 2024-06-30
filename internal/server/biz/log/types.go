package log

type GetLogRequest struct {
	Id *uint32 `json:"id"`
}

type ListLogRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	TaskId   uint32  `json:"taskId"`
	Order    *string `json:"order"`
	OrderBy  *string `json:"orderBy"`
}
