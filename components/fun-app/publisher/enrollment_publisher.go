package publisher

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

// EnrollmentEventPublisher defines publishing APIs for enrollment saga commands and events.
type EnrollmentPublisher interface {
	Enroll(ctx context.Context, enrollment fun.Enrollment) common.HttpError
	EnrollmentConfirmedEvt(ctx context.Context, enrollment fun.Enrollment) common.HttpError
}

type enrollmentPublisher struct {
	base BasePublisher
}

// NewEnrollmentPublisher builds a new EnrollmentPublisher backed by the provided base publisher.
func NewEnrollmentPublisher(base BasePublisher) EnrollmentPublisher {
	return &enrollmentPublisher{base: base}
}

func (ep *enrollmentPublisher) Enroll(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	payload := fun.EnrollCmdV1{
		EnrollmentID: enrollment.ID,
		PersonID:     enrollment.PersonID,
		Grade:        enrollment.Grade,
		Status:       enrollment.Status,
		RequestedAt:  time.Now().UTC(),
	}

	extras := map[string]string{
		fun.MetadataEnrollmentID: enrollment.ID,
		fun.MetadataPersonID:     enrollment.PersonID,
	}

	return ep.base.PublishWithExtras(ctx, fun.TopicEnrollCmd, payload, extras)
}

func (ep *enrollmentPublisher) EnrollmentConfirmedEvt(ctx context.Context, enrollment fun.Enrollment) common.HttpError {
	payload := fun.EnrollmentConfirmedEvtV1{
		EnrollmentID: enrollment.ID,
		PersonID:     enrollment.PersonID,
		ConfirmedAt:  time.Now().UTC(),
	}

	extras := map[string]string{
		fun.MetadataEnrollmentID: enrollment.ID,
		fun.MetadataPersonID:     enrollment.PersonID,
	}

	return ep.base.PublishWithExtras(ctx, fun.TopicEnrollmentConfirmedEvt, payload, extras)
}
