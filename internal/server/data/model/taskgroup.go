package model

import (
	"github.com/limes-cloud/kratosx/types"
)

type TaskGroup struct {
	Name        string  `json:"name" gorm:"column:name"`
	Description *string `json:"description" gorm:"column:description"`
	types.BaseModel
}
