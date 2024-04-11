package task

type Factory interface {
	DrySpec(s string) bool
	AddCron(id uint32, spec string) error
	UpdateCron(id uint32, spec string) error
	DeleteCron(id uint32) error
}
