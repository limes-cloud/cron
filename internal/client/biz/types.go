package biz

type ExecTaskReply struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Time    uint32 `json:"time"`
}

type ExecTaskReplyFunc func(*ExecTaskReply) error
