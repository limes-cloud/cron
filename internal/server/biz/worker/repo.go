package worker

import "github.com/limes-cloud/kratosx"

type Repo interface {
	AddWorker(ctx kratosx.Context, in *Worker) (uint32, error)
	GetWorker(ctx kratosx.Context, id uint32) (*Worker, error)
	GetWorkersByTag(ctx kratosx.Context, tag string) ([]*Worker, error)
	PageWorker(ctx kratosx.Context, req *PageWorkerRequest) ([]*Worker, uint32, error)
	UpdateWorker(ctx kratosx.Context, c *Worker) error
	DeleteWorker(ctx kratosx.Context, id uint32) error
}
