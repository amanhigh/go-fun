package publisher

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/fun"
)

// EnrollmentPublisher publishes only the root enrollment command.
type EnrollmentPublisher interface {
	Enroll(ctx context.Context, enrollment fun.Enrollment) common.HttpError
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
		StudentID:    enrollment.StudentID,
		Grade:        enrollment.Grade,
		RequestedAt:  time.Now().UTC(),
	}

	return ep.base.PublishRoot(ctx, fun.TopicEnrollCmd, payload)
}
