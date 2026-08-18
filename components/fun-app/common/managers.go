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
	container.MustSingleton(fi.di, fi.provideStudentManager)
	container.MustSingleton(fi.di, fi.provideSeatManager)
	container.MustSingleton(fi.di, fi.provideEnrollmentManager)
}

func (fi *FunAppInjector) provideStudentManager(studentDao dao.StudentDaoInterface, tracer trace.Tracer) manager.StudentManagerInterface {
	return manager.NewStudentManager(studentDao, tracer)
}

func (fi *FunAppInjector) provideSeatManager(seatPublisher publisher.SeatAllocationPublisher) manager.SeatManagerInterface {
	return manager.NewSeatManager(seatPublisher)
}

func (fi *FunAppInjector) provideEnrollmentManager(
	studentManager manager.StudentManagerInterface,
	enrollmentDao dao.EnrollmentDaoInterface,
	enrollmentPublisher publisher.EnrollmentPublisher,
	seatManager manager.SeatManagerInterface,
) manager.EnrollmentManagerInterface {
	return manager.NewEnrollmentManager(studentManager, enrollmentDao, enrollmentPublisher, seatManager)
}
