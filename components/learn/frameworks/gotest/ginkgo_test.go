package gotest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gmeasure"
	gomock "go.uber.org/mock/gomock"
)

var _ = Describe("Json Encode/Decode", func() {
	var (
		originalPerson Person
		personJson     string
	)
	BeforeEach(func() {
		originalPerson = Person{"Zoye", 44, 8983333}
		personJson = fmt.Sprintf(`{"name":"%s","Age":%d,"MobileNumber":%d}`, originalPerson.Name, originalPerson.Age, originalPerson.MobileNumber)
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

	It("should do things efficiently", FlakeAttempts(3), func() {
		action := "Encode"
		experiment := gmeasure.NewExperiment("Json Handling")
		AddReportEntry(experiment.Name, experiment)

		experiment.SampleDuration(action, func(_ int) {
			output, _ := encodePerson(originalPerson)
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
				encodeCall = mockEncoder.EXPECT().EncodePerson(gomock.Eq(originalPerson)).Return(personJson, nil)
			})

			It("should return mocked json", func() {
				json, err := mockEncoder.EncodePerson(originalPerson)
				Expect(err).To(BeNil())
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
					decodeCall = mockEncoder.EXPECT().DecodePerson(personJson).Return(per, nil)
					encodeCall.After(decodeCall)

					//gomock.InOrder(
					//	mockEncoder.EXPECT().DecodePerson(personJson).Return(per, nil),
					//	mockEncoder.EXPECT().EncodePerson(gomock.Eq(per)).Return(personJson, nil),
					//)
				})

				It("should decode", func() {
					decodedPerson, err := mockEncoder.DecodePerson(personJson)
					Expect(err).To(BeNil())
					Expect(decodedPerson).To(Equal(originalPerson))
					mockEncoder.EncodePerson(originalPerson)

				})
			})

			// TODO: Custom Matcher https://medium.com/modanisa-engineering/writing-a-custom-matcher-for-testing-with-gomock-d3ef5f13db82
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

		It("should build", func() {
			Expect(mockEncoder).To(Not(BeNil()))
		})

		It("should return mocked json", func() {
			mockEncoder.EXPECT().EncodePerson(originalPerson).Return(personJson, nil)

			json, err := mockEncoder.EncodePerson(originalPerson)
			Expect(err).To(BeNil())
			Expect(json).To(Equal(personJson))
		})

		Context("Match Field", func() {

			It("should match age", func() {
			})

		})

	})
})
