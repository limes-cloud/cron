package task

type PageTaskRequest struct {
	Page     uint32
	PageSize uint32
	Tag      *string
	Status   *string
	Name     *string
}

type PageTaskGroupRequest struct {
	Page     uint32
	PageSize uint32
	Tag      *string
	Status   *string
	Name     *string
}
