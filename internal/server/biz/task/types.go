package task

type GetTaskGroupRequest struct {
	Id *uint32 `json:"id"`
}

type ListTaskGroupRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Order    *string `json:"order"`
	OrderBy  *string `json:"orderBy"`
	Name     *string `json:"name"`
}

type GetTaskRequest struct {
	Id *uint32 `json:"id"`
}

type ListTaskRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Order    *string `json:"order"`
	OrderBy  *string `json:"orderBy"`
	GroupId  *uint32 `json:"groupId"`
	Name     *string `json:"name"`
	Tag      *string `json:"tag"`
	Status   *bool   `json:"status"`
}
