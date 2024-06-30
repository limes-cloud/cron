package model

import (
	"github.com/limes-cloud/kratosx/types"
)

type Worker struct {
	Name        string       `json:"name" gorm:"column:name"`
	Ip          string       `json:"ip" gorm:"column:ip"`
	GroupId     *uint32      `json:"groupId" gorm:"column:group_id"`
	Status      *bool        `json:"status" gorm:"column:status"`
	Description *string      `json:"description" gorm:"column:description"`
	Group       *WorkerGroup `json:"group"`
	types.BaseModel
}
