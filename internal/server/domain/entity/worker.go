package entity

import "github.com/limes-cloud/kratosx/model"

type WorkerGroup struct {
	Name        string  `json:"name" gorm:"column:name"`
	Description *string `json:"description" gorm:"column:description"`
	model.BaseModel
}

type Worker struct {
	Name        string       `json:"name" gorm:"column:name"`
	Apps        []string     `json:"apps" gorm:"-"`       // 可用应用白名单（worker_app 表）；空 = 所有应用可用
	Ip          string       `json:"ip" gorm:"column:ip"`
	Ak          string       `json:"ak" gorm:"column:ak"`
	Sk          string       `json:"sk" gorm:"column:sk"`
	GroupId     *uint32      `json:"groupId" gorm:"column:group_id"`
	Status      *bool        `json:"status" gorm:"column:status"`
	LastCheckAt *int64       `json:"lastCheckAt" gorm:"column:last_check_at"`       // 最后一次探测时间（含失败，Unix 秒）
	LastHeartbeat *int64     `json:"lastHeartbeat" gorm:"column:last_heartbeat"`    // 最后一次有效探测时间（成功，Unix 秒；NULL=从未上线）
	LastError   *string      `json:"lastError" gorm:"column:last_error"`            // 最后一次探测失败原因（成功时清空）
	Description *string      `json:"description" gorm:"column:description"`
	Group       *WorkerGroup `json:"group"`
	model.BaseModel
}

// WorkerApp 节点可用应用授权（worker_app 表；表达"该节点可供哪些 app 使用"）
type WorkerApp struct {
	Id        uint32 `json:"id"        gorm:"primaryKey;column:id"`
	WorkerId  uint32 `json:"workerId"  gorm:"column:worker_id"` // 节点ID
	App       string `json:"app"       gorm:"column:app"`       // 可用应用（X-App-Key）
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt int64  `json:"updatedAt" gorm:"column:updated_at"`
}
