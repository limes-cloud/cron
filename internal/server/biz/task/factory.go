package task

import "github.com/limes-cloud/kratosx"

type Factory interface {
	DrySpec(s string) bool
	AddCron(id uint32, spec string) error
	UpdateCron(id uint32, spec string) error
	RemoveCron(id uint32) error
	Scheduler(id uint32, force bool) error
	CancelExec(ctx kratosx.Context, uuid string) error
}
