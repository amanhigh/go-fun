package common_test

import (
	"net/http"

	"github.com/amanhigh/go-fun/models/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HttpError", func() {
	Context("Standard Http Errors", func() {
		It("should have correct error codes and messages", func() {
			Expect(common.ErrBadRequest.Code()).To(Equal(http.StatusBadRequest))
			Expect(common.ErrBadRequest.Error()).To(Equal("BadRequest"))

			Expect(common.ErrNotFound.Code()).To(Equal(http.StatusNotFound))
			Expect(common.ErrNotFound.Error()).To(Equal("NotFound"))

			Expect(common.ErrNotAuthorized.Code()).To(Equal(http.StatusUnauthorized))
			Expect(common.ErrNotAuthorized.Error()).To(Equal("NotAuthorized"))

			Expect(common.ErrNotAuthenticated.Code()).To(Equal(http.StatusForbidden))
			Expect(common.ErrNotAuthenticated.Error()).To(Equal("NotAuthenticated"))

			Expect(common.ErrEntityExists.Code()).To(Equal(http.StatusConflict))
			Expect(common.ErrEntityExists.Error()).To(Equal("EntityExists"))

			Expect(common.ErrInternalServerError.Code()).To(Equal(http.StatusInternalServerError))
			Expect(common.ErrInternalServerError.Error()).To(Equal("InternalServerError"))
		})
	})

	Context("NewHttpError", func() {
		It("should create HttpError with custom message and code", func() {
			customErr := common.NewHttpError("Custom Error", http.StatusTeapot)

			Expect(customErr.Error()).To(Equal("Custom Error"))
			Expect(customErr.Code()).To(Equal(http.StatusTeapot))
		})

		It("should implement error interface", func() {
			var err error = common.NewHttpError("Test", http.StatusBadRequest)
			Expect(err.Error()).To(Equal("Test"))
		})
	})

	Context("NewServerError", func() {
		It("should create HttpError from standard error", func() {
			originalErr := &customError{message: "database connection failed"}
			httpErr := common.NewServerError(originalErr)

			Expect(httpErr.Error()).To(Equal("database connection failed"))
			Expect(httpErr.Code()).To(Equal(http.StatusInternalServerError))
		})
	})

	Context("HttpErrorImpl", func() {
		It("should implement HttpError interface", func() {
			var httpErr common.HttpError = &common.HttpErrorImpl{
				Msg:     "Test Message",
				ErrCode: http.StatusBadRequest,
			}

			Expect(httpErr.Error()).To(Equal("Test Message"))
			Expect(httpErr.Code()).To(Equal(http.StatusBadRequest))
		})
	})
})

var _ = Describe("Pagination", func() {
	Context("Struct Fields", func() {
		It("should have correct field types and tags", func() {
			pagination := common.Pagination{
				Offset: 10,
				Limit:  5,
			}

			Expect(pagination.Offset).To(Equal(10))
			Expect(pagination.Limit).To(Equal(5))
		})

		It("should work with zero values", func() {
			pagination := common.Pagination{}

			Expect(pagination.Offset).To(Equal(0))
			Expect(pagination.Limit).To(Equal(0))
		})
	})
})

var _ = Describe("Sort", func() {
	Context("Struct Fields", func() {
		It("should have correct field types", func() {
			sort := common.Sort{
				SortBy: "name",
				Order:  "asc",
			}

			Expect(sort.SortBy).To(Equal("name"))
			Expect(sort.Order).To(Equal("asc"))
		})

		It("should work with empty values", func() {
			sort := common.Sort{}

			Expect(sort.SortBy).To(Equal(""))
			Expect(sort.Order).To(Equal(""))
		})
	})
})

var _ = Describe("PaginatedResponse", func() {
	Context("Struct Fields", func() {
		It("should have correct field types", func() {
			response := common.PaginatedResponse{
				Total: 100,
			}

			Expect(response.Total).To(Equal(int64(100)))
		})

		It("should work with zero value", func() {
			response := common.PaginatedResponse{}

			Expect(response.Total).To(Equal(int64(0)))
		})
	})
})

// Test helper
type customError struct {
	message string
}

func (e *customError) Error() string {
	return e.message
}
