package interfaces

type ServerInterface interface {
	Start() (err error)
	Stop()
}
