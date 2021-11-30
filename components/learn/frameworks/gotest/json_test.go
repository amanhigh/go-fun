package gotest

import (
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("Json Encode/Decode", func() {
	var (
		name           = "Zoye"
		age            = 44
		number         = int64(88983333)
		originalPerson person
		personJson     string
	)
	BeforeEach(func() {
		originalPerson = person{name, age, number}
		personJson = fmt.Sprintf(`{"name":"%s","Age":%d,"MobileNumber":%d}`, name, age, number)
	})
	Context("Success", func() {
		It("should encode Properly", func() {
			jsonString, err := encodePerson(originalPerson)
			Expect(err).To(BeNil())
			Expect(jsonString).To(Equal(personJson))
		})

		It("should decode Properly", func() {
			decodedPerson, err := decodePerson(personJson)
			Expect(err).To(BeNil())
			Expect(decodedPerson).To(Equal(originalPerson))
		})
	})

	Context("Fail", func() {
		BeforeEach(func() {
			originalPerson.Name = "Bob"
		})

		It("should throw error on invalid json", func() {
			_, err := decodePerson("abcd")
			Expect(err).To(Not(BeNil()))
		})

		It("should not match original person", func() {
			jsonString, err := encodePerson(originalPerson)
			Expect(err).To(BeNil())
			Expect(jsonString).To(Not(Equal(personJson)))
		})
	})

	Context("Encode", func() {
		var (
			jsonString string
			err        error
		)

		AfterEach(func() {
			Expect(err).To(BeNil())

		})

		JustAfterEach(func() {
			//Creation
			jsonString, err = encodePerson(originalPerson)
		})

		Context("Success", func() {
			AfterEach(func() {
				Expect(jsonString).To(Equal(personJson))
			})

			It("should encode Properly", func() {
			})
		})

		Context("Fail", func() {
			//Assertions
			AfterEach(func() {
				Expect(jsonString).To(Not(Equal(personJson)))
			})

			//Configuration
			It("with changed age", func() {
				originalPerson.Age = 88
			})

			It("with changed name", func() {
				originalPerson.Name = "Bob"
			})
		})
	})

	It("should do things efficiently", func() {
		action := "Encode"
		experiment := gmeasure.NewExperiment("Json Handling")
		AddReportEntry(experiment.Name, experiment)

		experiment.SampleDuration(action, func(_ int) {
			output, _ := encodePerson(originalPerson)
			Expect(output).To(Equal(personJson))
		}, gmeasure.SamplingConfig{N: 1000})
		AddReportEntry(action, experiment.GetStats(action))

		Expect(experiment.GetStats(action).DurationFor(gmeasure.StatMax)).To(BeNumerically("<", 3*time.Millisecond))
	})

	Context("Interesting Assertions", func() {
		var (
			err error
		)

		It("should match", func() {
			//Symbol Equivalent to Expect
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

			Ω(FoodSrce(pizza)).ShouldNot(Equal(pizza))       //will fail
			Ω(FoodSrce(pizza)).Should(BeEquivalentTo(pizza)) //will pass
		})

		It("should match collection", func() {
			//Array
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf("FooBar", "Foo"))
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf(ContainSubstring("Bar"), "Foo"))
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf(ContainSubstring("Foo"), ContainSubstring("Foo")))
			Ω([]string{"Foo", "FooBar"}).Should(ConsistOf([]string{"FooBar", "Foo"}))

			//Map
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
				c = make(chan string, 0)
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
				//verify response
				Expect(err).To(BeNil())
				Expect(response.StatusCode).To(Equal(http.StatusOK))
				actualResponse, err := ioutil.ReadAll(response.Body)
				Expect(err).To(BeNil())
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
		Context("Mock", func() {
			var (
				encodeCall *gomock.Call
			)
			BeforeEach(func() {
				encodeCall = mockEncoder.EXPECT().encodePerson(gomock.Eq(per)).Return(personJson, nil)
			})

			It("should return mocked json", func() {
				json, err := mockEncoder.encodePerson(per)
				Expect(err).To(BeNil())
				Expect(json).To(Equal(personJson))
			})

			Context("Do", func() {
				var (
					copiedPerson person
				)
				BeforeEach(func() {
					encodeCall.DoAndReturn(func(per person) { copiedPerson = per }).Return(personJson, nil)
				})

				It("should Do Something", func() {
					mockEncoder.encodePerson(per)
					Expect(copiedPerson).To(Equal(per))
				})

			})

			Context("Order", func() {
				var (
					decodeCall *gomock.Call
				)
				BeforeEach(func() {
					decodeCall = mockEncoder.EXPECT().decodePerson(personJson).Return(per, nil)
					encodeCall.After(decodeCall)

					//gomock.InOrder(
					//	mockEncoder.EXPECT().decodePerson(personJson).Return(per, nil),
					//	mockEncoder.EXPECT().encodePerson(gomock.Eq(per)).Return(personJson, nil),
					//)
				})

				It("should decode", func() {
					decodedPerson, err := mockEncoder.decodePerson(personJson)
					Expect(err).To(BeNil())
					Expect(decodedPerson).To(Equal(per))
					mockEncoder.encodePerson(per)

				})
			})
		})

		AfterEach(func() {
			ctrl.Finish()
		})
	})
})
