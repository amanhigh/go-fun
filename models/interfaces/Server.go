package interfaces

import "context"

type ServerInterface interface {
	Start(c context.Context) (err error)
	Stop(c context.Context)
}
