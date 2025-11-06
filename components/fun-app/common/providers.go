package common

import (
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/fun-app/dao"
	"github.com/amanhigh/go-fun/components/fun-app/manager"
	"github.com/amanhigh/go-fun/components/fun-app/publisher"
	"go.opentelemetry.io/otel/trace"
)

// Provider helpers returning interfaces while delegating to pointer-returning constructors.

func ProvidePersonDao(base util.BaseDao) dao.PersonDaoInterface {
	return dao.NewPersonDao(base)
}

func ProvideEnrollmentDao(base util.BaseDao) dao.EnrollmentDaoInterface {
	return dao.NewEnrollmentDao(base)
}

func ProvidePersonManager(dao dao.PersonDaoInterface, tracer trace.Tracer) manager.PersonManagerInterface {
	return manager.NewPersonManager(dao, tracer)
}

func ProvideSeatManager(seatPublisher publisher.SeatAllocationPublisher) manager.SeatManagerInterface {
	return manager.NewSeatManager(seatPublisher)
}

func ProvideEnrollmentManager(
	personManager manager.PersonManagerInterface,
	enrollmentDao dao.EnrollmentDaoInterface,
	enrollmentPublisher publisher.EnrollmentPublisher,
	seatManager manager.SeatManagerInterface,
) manager.EnrollmentManagerInterface {
	return manager.NewEnrollmentManager(personManager, enrollmentDao, enrollmentPublisher, seatManager)
}
