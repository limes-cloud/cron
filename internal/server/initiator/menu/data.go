package menu

import (
	"github.com/limes-cloud/kratosx/pkg/value"
	managerMenu "github.com/limes-cloud/manager/api/menu"
)

const (
	appKeyword = "cron"

	menuTypeApi   = "A"
	menuTypePage  = "M"
	menuTypeGroup = "G"
)

// menus cron 定时任务的菜单树：页面（M）+ 接口（A）两级结构，与 admin-web 的 views/apps/cron 对应
var menus = []*managerMenu.RegisterMenuNode{
	{
		Key:       value.Pointer("CronTask"),
		Title:     "定时任务",
		Type:      menuTypePage,
		Icon:      value.Pointer("iconsys-clock-in"),
		Path:      value.Pointer("/cron/task"),
		Component: value.Pointer("/cron/task"),
		Weight:    value.Pointer[int32](40),
		IsHidden:  value.Pointer(false),
		IsCache:   value.Pointer(true),
		IsAffix:   value.Pointer(true),
		Children: []*managerMenu.RegisterMenuNode{
			group("任务分组", "cron:task_group:group", []*managerMenu.RegisterMenuNode{
				api("查询任务分组", "/cron/api/v1/task_groups", "GET", "cron:task_group:query"),
				api("新增任务分组", "/cron/api/v1/task_group", "POST", "cron:task_group:add"),
				api("修改任务分组", "/cron/api/v1/task_group", "PUT", "cron:task_group:update"),
				api("删除任务分组", "/cron/api/v1/task_group", "DELETE", "cron:task_group:delete"),
			}),
			group("任务", "cron:task:group", []*managerMenu.RegisterMenuNode{
				api("查询任务", "/cron/api/v1/tasks", "GET", "cron:task:query"),
				api("新增任务", "/cron/api/v1/task", "POST", "cron:task:add"),
				api("修改任务", "/cron/api/v1/task", "PUT", "cron:task:update"),
				api("删除任务", "/cron/api/v1/task", "DELETE", "cron:task:delete"),
				api("启用/禁用任务", "/cron/api/v1/task/status", "PUT", "cron:task:update-status"),
				api("立即执行任务", "/cron/api/v1/task/exec", "POST", "cron:task:exec"),
				api("取消执行任务", "/cron/api/v1/task/cancel", "POST", "cron:task:cancel"),
			}),
			group("执行日志", "cron:log:group", []*managerMenu.RegisterMenuNode{
				api("查询执行日志", "/cron/api/v1/logs", "GET", "cron:log:query"),
				api("查询日志详情", "/cron/api/v1/log", "GET", "cron:log:detail"),
			}),
		},
	},
	{
		// 独立顶级菜单（隐藏）：/cron/task/group 作为一级路由渲染，
		// 若作为 CronTask 的子路由，父组件（index.vue）无 <router-view> 会导致子路由渲染时停留在父组件
		Key:       value.Pointer("CronTaskGroup"),
		Title:     "任务分组",
		Type:      menuTypePage,
		Path:      value.Pointer("/cron/task/group"),
		Component: value.Pointer("/cron/task/group"),
		IsHidden:  value.Pointer(true),
		Weight:    value.Pointer[int32](39),
	},
	{
		Key:       value.Pointer("CronWorker"),
		Title:     "节点管理",
		Type:      menuTypePage,
		Icon:      value.Pointer("iconsys-shebei"),
		Path:      value.Pointer("/cron/worker"),
		Component: value.Pointer("/cron/worker"),
		Weight:    value.Pointer[int32](50),
		IsHidden:  value.Pointer(false),
		IsCache:   value.Pointer(true),
		IsAffix:   value.Pointer(true),
		Children: []*managerMenu.RegisterMenuNode{
			group("节点分组", "cron:worker_group:group", []*managerMenu.RegisterMenuNode{
				api("查询节点分组", "/cron/api/v1/worker_groups", "GET", "cron:worker_group:query"),
				api("新增节点分组", "/cron/api/v1/worker_group", "POST", "cron:worker_group:add"),
				api("修改节点分组", "/cron/api/v1/worker_group", "PUT", "cron:worker_group:update"),
				api("删除节点分组", "/cron/api/v1/worker_group", "DELETE", "cron:worker_group:delete"),
			}),
			group("节点", "cron:worker:group", []*managerMenu.RegisterMenuNode{
				api("查询节点", "/cron/api/v1/workers", "GET", "cron:worker:query"),
				api("新增节点", "/cron/api/v1/worker", "POST", "cron:worker:add"),
				api("修改节点", "/cron/api/v1/worker", "PUT", "cron:worker:update"),
				api("删除节点", "/cron/api/v1/worker", "DELETE", "cron:worker:delete"),
				api("启用/禁用节点", "/cron/api/v1/worker/status", "PUT", "cron:worker:update-status"),
			}),
		},
	},
	{
		// 独立顶级菜单（隐藏）：/cron/worker/group 作为一级路由渲染
		Key:       value.Pointer("CronWorkerGroup"),
		Title:     "节点分组",
		Type:      menuTypePage,
		Path:      value.Pointer("/cron/worker/group"),
		Component: value.Pointer("/cron/worker/group"),
		IsHidden:  value.Pointer(true),
		Weight:    value.Pointer[int32](49),
	},
}

func group(title string, permission string, children []*managerMenu.RegisterMenuNode) *managerMenu.RegisterMenuNode {
	return &managerMenu.RegisterMenuNode{
		Title:      title,
		Type:       menuTypeGroup,
		Permission: value.Pointer(permission),
		Weight:     value.Pointer[int32](0),
		Children:   children,
	}
}

func api(title string, path string, method string, permission string) *managerMenu.RegisterMenuNode {
	node := &managerMenu.RegisterMenuNode{
		Title:  title,
		Type:   menuTypeApi,
		Api:    value.Pointer(path),
		Method: value.Pointer(method),
		Weight: value.Pointer[int32](0),
	}
	if permission != "" {
		node.Permission = value.Pointer(permission)
	}
	return node
}