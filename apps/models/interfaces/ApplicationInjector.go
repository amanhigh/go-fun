package interfaces

type ApplicationInjector interface {
	BuildApp() (app interface{}, err error)
}
