package fun

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	EnrollmentStatusSeatAllocationInitiated = "SEAT_ALLOCATION_INITIATED"
	EnrollmentStatusWaitlisted              = "WAITLISTED"
	EnrollmentStatusConfirmed               = "CONFIRMED"
	EnrollmentStatusCancelled               = "CANCELLED"
)

// EnrollmentRequest drives the enrollment orchestration using an existing student.
type EnrollmentRequest struct {
	StudentID string `json:"studentId" binding:"required"`
	Grade    int    `json:"grade" binding:"required,min=1,max=12"`
}

type EnrollmentPath struct {
	StudentID string `uri:"studentId" binding:"required"`
}

type Enrollment struct {
	ID        string    `gorm:"primaryKey" json:"enrollmentId"`
	StudentID  string    `gorm:"not null;uniqueIndex" json:"studentId"`
	Grade     int       `gorm:"not null" json:"grade"`
	Status    string    `gorm:"not null" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Student    Student    `gorm:"foreignKey:StudentID;references:Id;constraint:OnDelete:CASCADE" json:"-"`
}

func (e *Enrollment) BeforeCreate(_ *gorm.DB) (err error) {
	e.ID = uuid.NewString()[:8]
	return
}

// EnrollCmdV1 triggers the enrollment saga flow.
type EnrollCmdV1 struct {
	EnrollmentID string    `json:"enrollmentId" validate:"required"`
	StudentID    string    `json:"studentId" validate:"required"`
	Grade        int       `json:"grade" validate:"required,min=1,max=12"`
	RequestedAt  time.Time `json:"requestedAt" validate:"required"`
}

// EnrollmentEvent is the shared base for all enrollment domain events.
type EnrollmentEvent struct {
	EnrollmentID string `json:"enrollmentId" validate:"required"`
	StudentID    string `json:"studentId" validate:"required"`
}
