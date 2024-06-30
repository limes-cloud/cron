package task

type TaskGroup struct {
	Id          uint32  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	CreatedAt   int64   `json:"createdAt"`
	UpdatedAt   int64   `json:"updatedAt"`
}

type Worker struct {
	Name string `json:"name"`
}

type WorkerGroup struct {
	Name string `json:"name"`
}

type Task struct {
	Id            uint32       `json:"id"`
	GroupId       uint32       `json:"groupId"`
	Name          string       `json:"name"`
	Tag           string       `json:"tag"`
	Spec          string       `json:"spec"`
	Status        *bool        `json:"status"`
	WorkerType    string       `json:"workerType"`
	WorkerGroupId *uint32      `json:"workerGroupId"`
	WorkerId      *uint32      `json:"workerId"`
	ExecType      string       `json:"execType"`
	ExecValue     string       `json:"execValue"`
	ExpectCode    uint32       `json:"expectCode"`
	RetryCount    uint32       `json:"retryCount"`
	RetryWaitTime uint32       `json:"retryWaitTime"`
	MaxExecTime   uint32       `json:"maxExecTime"`
	Version       string       `json:"version"`
	Description   *string      `json:"description"`
	CreatedAt     int64        `json:"createdAt"`
	UpdatedAt     int64        `json:"updatedAt"`
	Group         *TaskGroup   `json:"group"`
	Worker        *Worker      `json:"worker"`
	WorkerGroup   *WorkerGroup `json:"workerGroup"`
}
