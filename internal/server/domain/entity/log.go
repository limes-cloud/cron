package entity

import "encoding/json"

type Log struct {
	Id             uint32 `json:"id" gorm:"primaryKey;column:id"`
	Uuid           string `json:"uuid" gorm:"column:uuid"`
	WorkerId       uint32 `json:"workerId" gorm:"column:worker_id"`
	WorkerSnapshot string `json:"workerSnapshot" gorm:"column:worker_snapshot"`
	TaskId         uint32 `json:"taskId" gorm:"column:task_id"`
	TaskSnapshot   string `json:"taskSnapshot" gorm:"column:task_snapshot"`
	StartAt        int64  `json:"startAt" gorm:"column:start_at"`
	EndAt          int64  `json:"endAt" gorm:"column:end_at"`
	Content        string `json:"content" gorm:"column:content"`
	Status         string `json:"status" gorm:"column:status"`
}

func (l *Log) GetTask() *Task {
	var task Task
	_ = json.Unmarshal([]byte(l.TaskSnapshot), &task)
	return &task
}

func (l *Log) GetWorker() *Worker {
	var worker Worker
	_ = json.Unmarshal([]byte(l.TaskSnapshot), &worker)
	return &worker
}
