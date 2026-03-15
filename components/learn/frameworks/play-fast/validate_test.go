package play_fast

import (
	"errors"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type ValidateTestStruct struct {
	// Basic types with common validations
	Name      string    `validate:"required,min=3,max=50"`
	Age       int       `validate:"required,min=18,max=120"`
	Price     float64   `validate:"required,gt=0"`
	Email     string    `validate:"required,email"`
	IsActive  bool      `validate:"required"`
	Tags      []string  `validate:"required,min=1,max=3"`
	BirthTime time.Time `validate:"required"`

	// DateTime format field (string for datetime validation)
	CreatedAt string `validate:"required,datetime=2006-01-02"`

	// DateTime format field (nullable pointer for datetime validation)
	ReviewedAt *string `validate:"omitempty,datetime=2006-01-02,not_future"`

	// Custom validator field
	SecretKey string `validate:"required,custom=secret"`

	// Custom regex field
	ProductCode string `validate:"required,regexp=^[A-Z]{2}-\\d{4}$"`
}

var _ = Describe("Validate", func() {
	var (
		validate  *validator.Validate
		err       error
		validData ValidateTestStruct
	)

	BeforeEach(func() {
		validate = validator.New()

		// Register custom validators
		err = validate.RegisterValidation("custom", func(fl validator.FieldLevel) bool {
			param := fl.Param()
			return fl.Field().String() == param
		})
		Expect(err).ToNot(HaveOccurred())

		// Register regex validator
		err = validate.RegisterValidation("regexp", func(fl validator.FieldLevel) bool {
			regex := fl.Param()
			matched, err := regexp.MatchString(regex, fl.Field().String())
			if err != nil {
				return false
			}
			return matched
		})
		Expect(err).ToNot(HaveOccurred())

		// Register not_future validator
		err = validate.RegisterValidation("not_future", func(fl validator.FieldLevel) bool {
			dateStr := fl.Field().String()
			if dateStr == "" {
				return true // Empty string is handled by omitempty
			}
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return false // Invalid format
			}
			now := time.Now()
			return parsedDate.Before(now) || parsedDate.Equal(now)
		})
		Expect(err).ToNot(HaveOccurred())

		// Initialize valid data struct
		validData = ValidateTestStruct{
			Name:        "John Doe",
			Age:         25,
			Price:       99.99,
			Email:       "john@example.com",
			IsActive:    true,
			Tags:        []string{"tech", "golang"},
			BirthTime:   time.Date(1990, 1, 1, 14, 30, 0, 0, time.UTC),
			CreatedAt:   "2023-12-25",
			ReviewedAt:  func() *string { s := "2023-12-25"; return &s }(),
			SecretKey:   "secret",
			ProductCode: "AB-1234",
		}
	})

	It("should build", func() {
		Expect(validate).To(Not(BeNil()))
	})

	Context("Valid Data", func() {
		It("should validate all basic types successfully", func() {
			err = validate.Struct(validData)
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("String Validations", func() {
		It("should fail required string", func() {
			validData.Name = "" // Empty to trigger required validation
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("Name"))
			Expect(fieldError.Tag()).To(Equal("required"))
		})

		It("should fail min length string", func() {
			validData.Name = "Jo" // Too short
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("Name"))
			Expect(fieldError.Tag()).To(Equal("min"))
			Expect(fieldError.Param()).To(Equal("3"))
		})

		It("should fail email format", func() {
			validData.Email = "invalid-email" // Invalid format
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("Email"))
			Expect(fieldError.Tag()).To(Equal("email"))
		})
	})

	Context("Number Validations", func() {
		It("should fail min int", func() {
			validData.Age = 17 // Below min 18
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("Age"))
			Expect(fieldError.Tag()).To(Equal("min"))
			Expect(fieldError.Param()).To(Equal("18"))
		})

		It("should fail max int", func() {
			validData.Age = 121 // Above max 120
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("Age"))
			Expect(fieldError.Tag()).To(Equal("max"))
			Expect(fieldError.Param()).To(Equal("120"))
		})

		It("should fail gt float64", func() {
			validData.Price = -1 // Not > 0
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("Price"))
			Expect(fieldError.Tag()).To(Equal("gt"))
			Expect(fieldError.Param()).To(Equal("0"))
		})
	})

	Context("Bool Validations", func() {
		It("should fail required bool", func() {
			validData.IsActive = false
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("IsActive"))
			Expect(fieldError.Tag()).To(Equal("required"))
		})
	})

	Context("Slice Validations", func() {
		It("should fail empty slice", func() {
			validData.Tags = []string{} // Empty slice
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("Tags"))
			Expect(fieldError.Tag()).To(Equal("min"))
			Expect(fieldError.Param()).To(Equal("1"))
		})

		It("should fail max slice length", func() {
			validData.Tags = []string{"one", "two", "three", "four"} // Too many items
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("Tags"))
			Expect(fieldError.Tag()).To(Equal("max"))
			Expect(fieldError.Param()).To(Equal("3"))
		})
	})

	Context("DateTime Validations (BirthTime - time.Time)", func() {
		It("should fail required datetime", func() {
			validData.BirthTime = time.Time{} // Zero time to trigger required validation
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("BirthTime"))
			Expect(fieldError.Tag()).To(Equal("required"))
		})
	})

	Context("DateTime Format Validations (CreatedAt - string)", func() {
		It("should pass datetime format", func() {
			// validData already has correct datetime format (2023-12-25)
			err = validate.Struct(validData)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail required datetime", func() {
			validData.CreatedAt = "" // Empty string to trigger required validation
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("CreatedAt"))
			Expect(fieldError.Tag()).To(Equal("required"))
		})

		It("should fail datetime format - includes time", func() {
			validData.CreatedAt = "2023-12-25T14:30:00Z" // ISO 8601 format with time
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("CreatedAt"))
			Expect(fieldError.Tag()).To(Equal("datetime"))
			Expect(fieldError.Param()).To(Equal("2006-01-02"))
		})

		It("should fail datetime format - wrong format", func() {
			validData.CreatedAt = "25/12/2023" // European date format
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("CreatedAt"))
			Expect(fieldError.Tag()).To(Equal("datetime"))
			Expect(fieldError.Param()).To(Equal("2006-01-02"))
		})

		It("should fail datetime format - invalid date", func() {
			validData.CreatedAt = "2023-13-45" // Invalid month and day
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("CreatedAt"))
			Expect(fieldError.Tag()).To(Equal("datetime"))
			Expect(fieldError.Param()).To(Equal("2006-01-02"))
		})
	})

	Context("DateTime Format Validations (ReviewedAt - *string)", func() {
		It("should pass with valid date", func() {
			// validData already has correct date format (2023-12-25)
			err = validate.Struct(validData)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should pass with nil pointer", func() {
			validData.ReviewedAt = nil // Nil pointer should be valid with omitempty
			err = validate.Struct(validData)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail with invalid date format", func() {
			invalidDate := "25/12/2023" // European date format
			validData.ReviewedAt = &invalidDate
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("ReviewedAt"))
			Expect(fieldError.Tag()).To(Equal("datetime"))
			Expect(fieldError.Param()).To(Equal("2006-01-02"))
		})

		It("should fail with datetime instead of date", func() {
			invalidDate := "2023-12-25T14:30:00Z" // DateTime format, should be date only
			validData.ReviewedAt = &invalidDate
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("ReviewedAt"))
			Expect(fieldError.Tag()).To(Equal("datetime"))
			Expect(fieldError.Param()).To(Equal("2006-01-02"))
		})

		It("should fail with invalid date", func() {
			invalidDate := "2023-13-45" // Invalid month and day
			validData.ReviewedAt = &invalidDate
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("ReviewedAt"))
			Expect(fieldError.Tag()).To(Equal("datetime"))
			Expect(fieldError.Param()).To(Equal("2006-01-02"))
		})

		It("should fail with future date", func() {
			futureDate := "2099-12-25" // Future date
			validData.ReviewedAt = &futureDate
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("ReviewedAt"))
			Expect(fieldError.Tag()).To(Equal("not_future"))
		})
	})

	Context("Custom Validations", func() {
		It("should pass custom validation", func() {
			// validData already has correct values
			err = validate.Struct(validData)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail custom validation", func() {
			validData.SecretKey = "wrong" // Doesn't match custom=secret
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("SecretKey"))
			Expect(fieldError.Tag()).To(Equal("custom"))
			Expect(fieldError.Param()).To(Equal("secret"))
		})
	})

	Context("Regex Validations", func() {
		It("should pass regex validation", func() {
			validData.ProductCode = "CD-9999" // Matches regex: 2 uppercase letters, dash, 4 digits
			err = validate.Struct(validData)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should fail regex validation - wrong format", func() {
			validData.ProductCode = "AB123" // Missing dash, wrong format
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("ProductCode"))
			Expect(fieldError.Tag()).To(Equal("regexp"))
			Expect(fieldError.Param()).To(Equal("^[A-Z]{2}-\\d{4}$"))
		})

		It("should fail regex validation - lowercase letters", func() {
			validData.ProductCode = "ab-1234" // Lowercase letters, should be uppercase
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("ProductCode"))
			Expect(fieldError.Tag()).To(Equal("regexp"))
			Expect(fieldError.Param()).To(Equal("^[A-Z]{2}-\\d{4}$"))
		})

		It("should fail regex validation - not enough digits", func() {
			validData.ProductCode = "AB-123" // Only 3 digits, need 4
			err = validate.Struct(validData)
			Expect(err).To(HaveOccurred())

			// Typecast to ValidationErrors and check field name
			var validationErrors validator.ValidationErrors
			Expect(errors.As(err, &validationErrors)).To(BeTrue())
			Expect(validationErrors).To(HaveLen(1))

			fieldError := validationErrors[0]
			Expect(fieldError.Field()).To(Equal("ProductCode"))
			Expect(fieldError.Tag()).To(Equal("regexp"))
			Expect(fieldError.Param()).To(Equal("^[A-Z]{2}-\\d{4}$"))
		})
	})

})
