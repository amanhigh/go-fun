package gotest

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
	"github.com/stretchr/testify/mock"
	gomock "go.uber.org/mock/gomock"
)

type personHasName struct{}

// Matches checks if the given interface{} is of type Person and has a non-empty Name field.
func (m *personHasName) Matches(x interface{}) bool {
	p, ok := x.(Person)
	return ok && p.Name != ""
}

// String returns a description of the matcher.
func (m *personHasName) String() string {
	return "is a person with a non-empty name"
}

var _ = Describe("Json Encode/Decode", func() {
	var (
		originalPerson Person
		personJson     string
		personEncoder  PersonEncoder
	)
	BeforeEach(func() {
		originalPerson = Person{"Zoye", 44, 8983333}
		personJson = fmt.Sprintf(`{"name":"%s","Age":%d,"MobileNumber":%d}`, originalPerson.Name, originalPerson.Age, originalPerson.MobileNumber)
		personEncoder = &PersonEncoderImpl{}
	})

	It("should build", func() {
		Expect(personEncoder).To(Not(BeNil()))
	})

	Context("Success", func() {
		It("should encode Properly", func() {
			jsonString, err := personEncoder.EncodePerson(originalPerson)
			Expect(err).ToNot(HaveOccurred())
			Expect(jsonString).To(Equal(personJson))
		})

		It("should decode Properly", func() {
			decodedPerson, err := personEncoder.DecodePerson(personJson)
			Expect(err).ToNot(HaveOccurred())
			Expect(decodedPerson).To(Equal(originalPerson))
		})
	})

	Context("Fail", func() {
		BeforeEach(func() {
			originalPerson.Name = "Bob"
		})

		It("should throw error on invalid json", func() {
			_, err := personEncoder.DecodePerson("abcd")
			Expect(err).To(HaveOccurred())
		})

		It("should not match original person", func() {
			jsonString, err := personEncoder.EncodePerson(originalPerson)
			Expect(err).ToNot(HaveOccurred())
			Expect(jsonString).To(Not(Equal(personJson)))
		})

		It("should fail for negative age", func() {
			originalPerson.Age = -44
			_, err := personEncoder.EncodePerson(originalPerson)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Encode", func() {
		var (
			jsonString string
			err        error
		)

		AfterEach(func() {
			Expect(err).ToNot(HaveOccurred())

		})

		JustAfterEach(func() {
			// Creation
			jsonString, err = personEncoder.EncodePerson(originalPerson)
		})

		Context("Success", func() {
			AfterEach(func() {
				Expect(jsonString).To(Equal(personJson))
			})

			It("should encode Properly", func() {
			})
		})

		Context("Fail", func() {
			// Assertions
			AfterEach(func() {
				Expect(jsonString).To(Not(Equal(personJson)))
			})

			// Configuration
			It("with changed age", func() {
				originalPerson.Age = 88
			})

			It("with changed name", func() {
				originalPerson.Name = "Bob"
			})
		})
	})

	It("should do things efficiently", FlakeAttempts(3), func() {
		action := "Encode"
		experiment := gmeasure.NewExperiment("Json Handling")
		AddReportEntry(experiment.Name, experiment)

		experiment.SampleDuration(action, func(_ int) {
			output, _ := personEncoder.EncodePerson(originalPerson)
			Expect(output).To(Equal(personJson))
		}, gmeasure.SamplingConfig{N: 1000})
		AddReportEntry(action, experiment.GetStats(action))

		Expect(experiment.GetStats(action).DurationFor(gmeasure.StatMax)).To(BeNumerically("<", 5*time.Millisecond))
	})

	Context("Interesting Assertions", func() {
		var (
			err error
		)

		It("should match", func() {
			// Symbol Equivalent to Expect
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("should async", func() {
			Eventually(func() []int {
				time.Sleep(time.Millisecond * 10)
				return []int{2, 3}
			}, time.Millisecond*100, time.Millisecond*2).Should(HaveLen(2))
		})

		It("should deep equal", func() {
			pizza := "Cheeseboard Pizza"
			type FoodSrce string

			Ω(FoodSrce(pizza)).ShouldNot(Equal(pizza))       // will fail
			Ω(FoodSrce(pizza)).Should(BeEquivalentTo(pizza)) // will pass
		})

		It("should match collection", func() {
			// Array
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf("FooBar", "Foo"))
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf(ContainSubstring("Bar"), "Foo"))
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf(ContainSubstring("Foo"), ContainSubstring("Foo")))
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf([]string{"FooBar", "Foo"}))

			// Map
			Ω(map[string]string{"Foo": "Bar", "BazFoo": "Duck"}).Should(HaveKey(MatchRegexp(`.+Foo$`)))
			Ω(map[string]int{"Foo": 3, "BazFoo": 4}).Should(HaveKeyWithValue(MatchRegexp(`.+Foo$`), BeNumerically(">", 3)))
		})

		It("should panic", func() {
			Ω(func() { panic("FooBarBaz") }).Should(Panic())
		})

		Context("Channel", func() {
			var (
				c chan string
			)

			BeforeEach(func() {
				c = make(chan string)
			})

			Context("With Close", func() {
				AfterEach(func() {
					Eventually(c).Should(BeClosed())
				})

				It("should receive", func() {
					go DoSomething(c, true)
				})

				It("should receive", func() {
					go DoSomething(c, true)
					Eventually(c).Should(Receive(Equal("Done!")))
				})
			})

			Context("No Close", func() {
				AfterEach(func() {
					Eventually(c).ShouldNot(BeClosed())
				})

				It("Channel Check Content", func() {
					go DoSomething(c, false)
					Expect(<-c).To(ContainSubstring("Done!"))
				})

				It("should not receive", func() {
					Consistently(c, time.Millisecond*100).ShouldNot(Receive())
				})
			})
		})
	})

	Context("Http Test Server", func() {
		var (
			server       *httptest.Server
			getResponse  = "Hello"
			postResponse = "World"

			err      error
			response *http.Response
		)

		BeforeEach(func() {
			server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					w.Write([]byte(getResponse))
				} else {
					w.Write([]byte(postResponse))
				}
			}))
		})

		It("should build", func() {
			Expect(server).To(Not(BeNil()))
		})

		Context("Calls", func() {
			var (
				expectedResponse string
			)

			AfterEach(func() {
				// verify response
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				actualResponse, err := io.ReadAll(response.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(actualResponse).To(BeEquivalentTo(expectedResponse))

				server.Close()
			})

			It("should do Get", func() {
				response, err = http.Get(server.URL)
				expectedResponse = getResponse
			})

			It("should do Post", func() {
				response, err = http.Post(server.URL, "", nil)
				expectedResponse = postResponse

			})
		})

	})

	Context("GoMock", func() {
		var (
			ctrl        *gomock.Controller
			mockEncoder *MockPersonEncoder
		)
		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())
			mockEncoder = NewMockPersonEncoder(ctrl)
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		Context("Mock", func() {
			var (
				encodeCall *gomock.Call
			)
			BeforeEach(func() {
				encodeCall = mockEncoder.EXPECT().EncodePerson(gomock.Eq(originalPerson)).Return(personJson, nil)
			})

			It("should return mocked json", func() {
				json, err := mockEncoder.EncodePerson(originalPerson)
				Expect(err).ToNot(HaveOccurred())
				Expect(json).To(Equal(personJson))
			})

			Context("Do", func() {
				var (
					copiedPerson Person
				)
				BeforeEach(func() {
					encodeCall.DoAndReturn(func(per Person) { copiedPerson = per }).Return(personJson, nil)
				})

				It("should Do Something", func() {
					mockEncoder.EncodePerson(originalPerson)
					Expect(copiedPerson).To(Equal(originalPerson))
				})

			})

			Context("Order", func() {
				var (
					decodeCall *gomock.Call
				)
				BeforeEach(func() {
					decodeCall = mockEncoder.EXPECT().DecodePerson(personJson).Return(person, nil)
					encodeCall.After(decodeCall)

					// gomock.InOrder(
					//	mockEncoder.EXPECT().DecodePerson(personJson).Return(per, nil),
					//	mockEncoder.EXPECT().EncodePerson(gomock.Eq(per)).Return(personJson, nil),
					//)
				})

				It("should decode", func() {
					decodedPerson, err := mockEncoder.DecodePerson(personJson)
					Expect(err).ToNot(HaveOccurred())
					Expect(decodedPerson).To(Equal(originalPerson))
					mockEncoder.EncodePerson(originalPerson)

				})
			})

			Context("Times", func() {
				It("should runs 3 times", func() {
					mockEncoder.EXPECT().EncodePerson(gomock.Eq(originalPerson)).Return(personJson, nil).Times(2)
					// Initial Expect + 2 in Times
					mockEncoder.EncodePerson(originalPerson)
					mockEncoder.EncodePerson(originalPerson)
					mockEncoder.EncodePerson(originalPerson)
				})

				It("should support zero times", func() {
					mockEncoder.EncodePerson(originalPerson)
					zeroPerson := Person{Name: "", Age: 0, MobileNumber: 0}
					mockEncoder.EXPECT().EncodePerson(gomock.Eq(zeroPerson)).Return(personJson, nil).Times(0)
				})
			})
		})

		Context("Custom Matcher", func() {
			It("should match person with non-empty name", func() {
				mockEncoder.EXPECT().EncodePerson(&personHasName{}).Return(personJson, nil)

				result, err := mockEncoder.EncodePerson(Person{Name: "Zoye", Age: 44, MobileNumber: 8983333})
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(personJson))
			})

			It("should not match person with empty name", func() {
				mockEncoder.EXPECT().EncodePerson(gomock.AssignableToTypeOf(Person{})).Return(personJson, nil)
				mockEncoder.EXPECT().EncodePerson(&personHasName{}).Return(personJson, errors.New("Empty Name")).Times(0)

				_, err := mockEncoder.EncodePerson(Person{Name: "", Age: 44, MobileNumber: 8983333})
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		AfterEach(func() {
			ctrl.Finish()
		})
	})

	//https://vektra.github.io/mockery/latest/examples/#simple-case
	Context("Mockery", func() {
		var (
			mockEncoder *MockEncoder
		)
		BeforeEach(func() {
			mockEncoder = NewMockEncoder(GinkgoT())
		})

		AfterEach(func() {
			mockEncoder.AssertExpectations(GinkgoT())
		})

		It("should build", func() {
			Expect(mockEncoder).To(Not(BeNil()))
		})

		It("should return mocked json", func() {
			mockEncoder.EXPECT().EncodePerson(originalPerson).Return(personJson, nil).Once()

			json, err := mockEncoder.EncodePerson(originalPerson)
			Expect(err).ToNot(HaveOccurred())
			Expect(json).To(Equal(personJson))
		})

		It("should mock in order", func() {
			encodeCall := mockEncoder.EXPECT().EncodePerson(originalPerson).Return(personJson, nil).Once()
			mockEncoder.EXPECT().DecodePerson(personJson).Return(originalPerson, nil).Times(1).NotBefore(encodeCall)

			mockEncoder.EncodePerson(originalPerson)
			mockEncoder.DecodePerson(personJson)
		})
		// https://pkg.go.dev/github.com/stretchr/testify/mock#pkg-index
		Context("Match Field", func() {
			It("should match given name", func() {
				mockEncoder.EXPECT().EncodePerson(mock.MatchedBy(func(inputPerson Person) bool {
					return inputPerson.Name == "Zoye"
				})).Return(personJson, nil)

				result, _ := mockEncoder.EncodePerson(originalPerson)
				Expect(result).To(Equal(personJson))
			})

			It("should match age > 50", func() {
				mockEncoder.EXPECT().EncodePerson(mock.AnythingOfType("Person")).RunAndReturn(func(inputPerson Person) (result string, err error) {
					if inputPerson.Age > 50 {
						return personJson, nil
					}
					return "", errors.New("Invalid Age")
				})

				_, err := mockEncoder.EncodePerson(originalPerson)
				Expect(err).Should(HaveOccurred())

				originalPerson.Age = 60
				result, err := mockEncoder.EncodePerson(originalPerson)
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(personJson))
			})
		})

		Context("Times", func() {
			It("should run 3 times", func() {
				mockEncoder.EXPECT().EncodePerson(mock.AnythingOfType("Person")).Return(personJson, nil).Times(3)
				mockEncoder.EncodePerson(originalPerson)
				mockEncoder.EncodePerson(originalPerson)
				mockEncoder.EncodePerson(originalPerson)
			})

			It("should support zero times", func() {
				// TASK: Doesn't Suport Times 0, https://github.com/stretchr/testify/issues/566
				Expect(mockEncoder.AssertNotCalled(GinkgoT(), "EncodePerson")).To(BeTrue())
			})
		})

		It("should do operation and return", func() {
			var copiedPerson Person

			mockEncoder.EXPECT().EncodePerson(mock.Anything).RunAndReturn(func(inputPerson Person) (result string, err error) {
				copiedPerson = inputPerson
				return "Aman", nil
			})

			result, _ := mockEncoder.EncodePerson(originalPerson)
			Expect(copiedPerson).To(Equal(originalPerson))
			Expect(result).To(Equal("Aman"))
		})

		It("should capture arguments", func() {
			var name string
			mockEncoder.EXPECT().EncodePerson(mock.AnythingOfType("Person")).Run(func(inputPerson Person) {
				name = inputPerson.Name
			}).Return("", nil)

			mockEncoder.EncodePerson(originalPerson)
			Expect(name).To(Equal(originalPerson.Name))
		})
	})
})
