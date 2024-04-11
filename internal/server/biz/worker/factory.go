package worker

import "github.com/limes-cloud/kratosx"

type Factory interface {
	CheckIP(ctx kratosx.Context, ip string) error
}
