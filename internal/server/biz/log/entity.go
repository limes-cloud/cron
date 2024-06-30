package log

type Log struct {
	Id             uint32 `json:"id"`
	Uuid           string `json:"uuid"`
	WorkerId       uint32 `json:"workerId"`
	WorkerSnapshot string `json:"workerSnapshot"`
	TaskId         uint32 `json:"taskId"`
	TaskSnapshot   string `json:"taskSnapshot"`
	StartAt        int64  `json:"startAt"`
	EndAt          int64  `json:"endAt"`
	Content        string `json:"content"`
	Status         string `json:"status"`
}
