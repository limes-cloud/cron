package types

type ListWorkerGroupRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Name     *string `json:"name"`
}

type GetWorkerRequest struct {
	Id *uint32 `json:"id"`
	Ip *string `json:"ip"`
}

type ListWorkerRequest struct {
	Page     uint32  `json:"page"`
	PageSize uint32  `json:"pageSize"`
	Name     *string `json:"name"`
	Ip       *string `json:"ip"`
	GroupId  *uint32 `json:"groupId"`
	Status   *bool   `json:"status"`
}

type CheckWorkerRequest struct {
	Uuid string
	IP   string
	Ak   string
	Sk   string
}
