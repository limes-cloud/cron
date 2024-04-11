package worker

type PageWorkerRequest struct {
	Page     uint32
	PageSize uint32
	Tag      *string
	Status   *string
	IP       *string
	Name     *string
}
