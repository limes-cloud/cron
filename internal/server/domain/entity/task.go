package entity

import "github.com/limes-cloud/kratosx/model"

// 错过补偿策略（misfire）：执行时长 ≥ 周期导致帧被跳过时的处理方式
const (
	MisfireSkip      = "skip"      // 丢弃本帧，记日志
	MisfireImmediate = "immediate" // 当前执行结束后立即补执行一次
	MisfireLatest    = "latest"    // 只补最近一帧（默认）
)

// RouteType 路由方式（P2 从 exec_type 拆分）
const (
	RouteGroup  = "group"  // 按节点分组随机
	RouteSingle = "single" // 指定单节点
)

type TaskGroup struct {
	App         string  `json:"app" gorm:"column:app"`             // 所属应用（X-App-Key）
	Name        string  `json:"name" gorm:"column:name"`
	Description *string `json:"description" gorm:"column:description"`
	model.BaseModel
}

type Task struct {
	App           string       `json:"app" gorm:"column:app"`              // 所属应用（X-App-Key）
	GroupId       uint32       `json:"groupId" gorm:"column:group_id"`
	Name          string       `json:"name" gorm:"column:name"`
	Tag           string       `json:"tag" gorm:"column:tag"`
	Spec          string       `json:"spec" gorm:"column:spec"`
	Status        *bool        `json:"status" gorm:"column:status"`
	WorkerType    string       `json:"workerType" gorm:"column:worker_type"`
	WorkerGroupId *uint32      `json:"workerGroupId" gorm:"column:worker_group_id"`
	WorkerId      *uint32      `json:"workerId" gorm:"column:worker_id"`
	ExecType      string       `json:"execType" gorm:"column:exec_type"`
	ExecValue     string       `json:"execValue" gorm:"column:exec_value"`
	ExpectCode    uint32       `json:"expectCode" gorm:"column:expect_code"`
	RetryCount    uint32       `json:"retryCount" gorm:"column:retry_count"`
	RetryWaitTime uint32       `json:"retryWaitTime" gorm:"column:retry_wait_time"`
	MaxExecTime   uint32       `json:"maxExecTime" gorm:"column:max_exec_time"`
	Version       string       `json:"version" gorm:"column:version"`
	Description   *string      `json:"description" gorm:"column:description"`
	Start         *uint32      `json:"start" gorm:"column:start"`
	End           *uint32      `json:"end" gorm:"column:end"`
	RouteType     string       `json:"routeType" gorm:"column:route_type"`     // 路由方式：group | single（P2 从 exec_type 拆分，空值兼容旧数据）
	Misfire       string       `json:"misfire" gorm:"column:misfire"`          // 错过补偿策略：skip | immediate | latest（空=默认 latest）
	Group         *TaskGroup   `json:"group"`
	Worker        *Worker      `json:"worker"`
	WorkerGroup   *WorkerGroup `json:"workerGroup"`
	model.BaseModel
}
