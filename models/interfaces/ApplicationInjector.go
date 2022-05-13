package interfaces

type ApplicationInjector interface {
	BuildApp() (app any, err error)
}
