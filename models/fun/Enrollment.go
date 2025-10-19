package fun

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	EnrollmentStatusSeatAllocationInitiated = "SEAT_ALLOCATION_INITIATED"
	EnrollmentStatusSeatReserved            = "SEAT_RESERVED"
	EnrollmentStatusWaitlisted              = "WAITLISTED"
	EnrollmentStatusConfirmed               = "CONFIRMED"
	EnrollmentStatusCancelled               = "CANCELLED"
)

// EnrollmentRequest drives the enrollment orchestration using an existing person.
type EnrollmentRequest struct {
	PersonID string `json:"personId" binding:"required"`
	Grade    int    `json:"grade" binding:"required,min=1,max=12"`
}

type EnrollmentPath struct {
	PersonID string `uri:"personId" binding:"required"`
}

type Enrollment struct {
	ID        string    `gorm:"primaryKey" json:"enrollmentId"`
	PersonID  string    `gorm:"not null;uniqueIndex" json:"personId"`
	Grade     int       `gorm:"not null" json:"grade"`
	Status    string    `gorm:"not null" json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Person    Person    `gorm:"foreignKey:PersonID;references:Id;constraint:OnDelete:CASCADE" json:"-"`
}

func (e *Enrollment) BeforeCreate(_ *gorm.DB) (err error) {
	e.ID = uuid.NewString()[:8]
	return
}

// EnrollmentResponse summarizes the enrollment outcome.
// TODO: Extend with richer metadata once student model is introduced.
type EnrollmentResponse struct {
	EnrollmentID string `json:"enrollmentId"`
	PersonID     string `json:"personId"`
	Grade        int    `json:"grade"`
	Status       string `json:"status"`
}
