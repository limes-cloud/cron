package entity

import "github.com/limes-cloud/kratosx/types"

type WorkerGroup struct {
	Name        string  `json:"name" gorm:"column:name"`
	Description *string `json:"description" gorm:"column:description"`
	types.BaseModel
}

type Worker struct {
	Name        string       `json:"name" gorm:"column:name"`
	Ip          string       `json:"ip" gorm:"column:ip"`
	Ak          string       `json:"ak" gorm:"column:ak"`
	Sk          string       `json:"sk" gorm:"column:sk"`
	GroupId     *uint32      `json:"groupId" gorm:"column:group_id"`
	Status      *bool        `json:"status" gorm:"column:status"`
	Description *string      `json:"description" gorm:"column:description"`
	Group       *WorkerGroup `json:"group"`
	types.BaseModel
}
