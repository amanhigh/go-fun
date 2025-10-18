package fun

// EnrollmentRequest drives the enrollment orchestration using an existing person.
type EnrollmentRequest struct {
	PersonID string `json:"personId" binding:"required"`
	Grade    int    `json:"grade" binding:"required,min=1,max=12"`
}

// EnrollmentResponse summarizes the enrollment outcome.
// TODO: Extend with richer metadata once student model is introduced.
type EnrollmentResponse struct {
	PersonID string `json:"personId"`
	Grade    int    `json:"grade"`
	Status   string `json:"status"`
}
