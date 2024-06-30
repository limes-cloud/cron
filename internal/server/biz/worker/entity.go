package worker

type WorkerGroup struct {
	Id          uint32  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedAt   int64   `json:"createdAt"`
	UpdatedAt   int64   `json:"updatedAt"`
}

type Worker struct {
	Id          uint32       `json:"id"`
	Name        string       `json:"name"`
	Ip          string       `json:"ip"`
	GroupId     *uint32      `json:"groupId"`
	Status      *bool        `json:"status"`
	Description *string      `json:"description"`
	CreatedAt   int64        `json:"createdAt"`
	UpdatedAt   int64        `json:"updatedAt"`
	Group       *WorkerGroup `json:"group"`
}
