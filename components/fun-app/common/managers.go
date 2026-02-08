package common

import (
	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"github.com/golobby/container/v3"
	"go.opentelemetry.io/otel/trace"
)

// Manager providers return interfaces while delegating to pointer-returning constructors.

func (fi *FunAppInjector) registerManager() {
	container.MustSingleton(fi.di, fi.providePersonManager)
	container.MustSingleton(fi.di, fi.provideSeatManager)
	container.MustSingleton(fi.di, fi.provideEnrollmentManager)
}

func (fi *FunAppInjector) providePersonManager(personDao dao.PersonDaoInterface, tracer trace.Tracer) manager.PersonManagerInterface {
	return manager.NewPersonManager(personDao, tracer)
}

func (fi *FunAppInjector) provideSeatManager(seatPublisher publisher.SeatAllocationPublisher) manager.SeatManagerInterface {
	return manager.NewSeatManager(seatPublisher)
}

func (fi *FunAppInjector) provideEnrollmentManager(
	personManager manager.PersonManagerInterface,
	enrollmentDao dao.EnrollmentDaoInterface,
	enrollmentPublisher publisher.EnrollmentPublisher,
	seatManager manager.SeatManagerInterface,
) manager.EnrollmentManagerInterface {
	return manager.NewEnrollmentManager(personManager, enrollmentDao, enrollmentPublisher, seatManager)
}
