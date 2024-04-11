package log

import "github.com/limes-cloud/kratosx"

type Repo interface {
	AddLog(ctx kratosx.Context, in *Log) (uint32, error)
	GetLog(ctx kratosx.Context, id uint32) (*Log, error)
	PageLog(ctx kratosx.Context, req *PageLogRequest) ([]*Log, uint32, error)
	AppendLogContent(ctx kratosx.Context, c string) error
	UpdateLogStatus(ctx kratosx.Context, err error) error
}
