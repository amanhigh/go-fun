package play_fast

import (
	"reflect"

	"github.com/bxcodec/faker/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// SomeStructWithTags ...
type SomeStructWithTags struct {
	Latitude           float32 `faker:"lat"`
	Longitude          float32 `faker:"long"`
	CreditCardNumber   string  `faker:"cc_number"`
	CreditCardType     string  `faker:"cc_type"`
	Email              string  `faker:"email"`
	DomainName         string  `faker:"domain_name"`
	IPV4               string  `faker:"ipv4"`
	IPV6               string  `faker:"ipv6"`
	Password           string  `faker:"password"`
	PhoneNumber        string  `faker:"phone_number"`
	MacAddress         string  `faker:"mac_address"`
	URL                string  `faker:"url"`
	UserName           string  `faker:"username"`
	TollFreeNumber     string  `faker:"toll_free_number"`
	E164PhoneNumber    string  `faker:"e_164_phone_number"`
	TitleMale          string  `faker:"title_male"`
	TitleFemale        string  `faker:"title_female"`
	FirstName          string  `faker:"first_name"`
	FirstNameMale      string  `faker:"first_name_male"`
	FirstNameFemale    string  `faker:"first_name_female"`
	LastName           string  `faker:"last_name"`
	Name               string  `faker:"name"`
	UnixTime           int64   `faker:"unix_time"`
	Date               string  `faker:"date"`
	Time               string  `faker:"time"`
	MonthName          string  `faker:"month_name"`
	Year               string  `faker:"year"`
	DayOfWeek          string  `faker:"day_of_week"`
	DayOfMonth         string  `faker:"day_of_month"`
	Timestamp          string  `faker:"timestamp"`
	Century            string  `faker:"century"`
	TimeZone           string  `faker:"timezone"`
	TimePeriod         string  `faker:"time_period"`
	Word               string  `faker:"word"`
	Sentence           string  `faker:"sentence"`
	Paragraph          string  `faker:"paragraph"`
	Currency           string  `faker:"currency"`
	Amount             float64 `faker:"amount"`
	AmountWithCurrency string  `faker:"amount_with_currency"`
	UUIDHypenated      string  `faker:"uuid_hyphenated"`
	UUID               string  `faker:"uuid_digit"`
	Skip               string  `faker:"-"`

	Inta  int   `faker:"boundary_start=5, boundary_end=10"`
	Int8  int8  `faker:"boundary_start=100, boundary_end=1000"`
	Int16 int16 `faker:"boundary_start=123, boundary_end=1123"`
	Int32 int32 `faker:"boundary_start=-10, boundary_end=8123"`
	Int64 int64 `faker:"boundary_start=31, boundary_end=88"`

	UInta  uint   `faker:"boundary_start=35, boundary_end=152"`
	UInt8  uint8  `faker:"boundary_start=5, boundary_end=1425"`
	UInt16 uint16 `faker:"boundary_start=245, boundary_end=2125"`
	UInt32 uint32 `faker:"boundary_start=0, boundary_end=40"`
	UInt64 uint64 `faker:"boundary_start=14, boundary_end=50"`

	StringArray  []string          `faker:"len=50"`
	SingleString string            `faker:"len=25"`
	StringMap    map[string]string `faker:"len=30"`
	IntMap       map[int]int       `faker:"boundary_start=5, boundary_end=10"`

	UniqueWord    string `faker:"word,unique"`
	PaymentMethod string `faker:"oneof: cc, paypal, check, money order"`

	NestedStruct AStruct
	//Struct Array Not Supported
	//NestedStructs []AStruct `faker:"boundary_start=0, boundary_end=4"`
}

// AStruct ...
type AStruct struct {
	Number        int64
	Height        int64
	AnotherStruct BStruct
}

// BStruct ...
type BStruct struct {
	Image string
}

/*  CUSTOM GENERATION */

// Gondoruwo ...
type Gondoruwo struct {
	Name       string
	Locatadata int
}

// Sample ...
type Sample struct {
	ID        int64     `faker:"customIdFaker"`
	Gondoruwo Gondoruwo `faker:"gondoruwo"`
	Danger    string    `faker:"danger"`
}

// CustomGenerator ...
func CustomGenerator() {
	_ = faker.AddProvider("customIdFaker", func(v reflect.Value) (any, error) {
		return int64(43), nil
	})
	_ = faker.AddProvider("danger", func(v reflect.Value) (any, error) {
		return "danger-ranger", nil
	})

	_ = faker.AddProvider("gondoruwo", func(v reflect.Value) (any, error) {
		obj := Gondoruwo{
			Name:       "Power",
			Locatadata: 324,
		}
		return obj, nil
	})
}

var _ = Describe("Faker", func() {

	var (
		err error
	)

	Context("Default", func() {
		var data = SomeStructWithTags{}

		BeforeEach(func() {
			err = faker.FakeData(&data)
			Expect(err).To(BeNil())
		})

		It("should generate data", func() {
			Expect(data.Email).To(Not(BeNil()))
			Expect(data.Date).To(Not(BeNil()))
			Expect(data.Date).To(Not(BeNil()))
		})

		It("should generate nested data", func() {
			Expect(data.NestedStruct.Number).To(Not(BeNil()))
			Expect(data.NestedStruct.AnotherStruct.Image).To(Not(BeNil()))
		})
	})

	Context("Custom", func() {
		var data = Sample{}
		BeforeEach(func() {
			CustomGenerator()
			err = faker.FakeData(&data)
			Expect(err).To(BeNil())
		})

		It("should generate", func() {
			Expect(data.ID).To(Equal(int64(43)))
			Expect(data.Danger).To(Equal("danger-ranger"))
			Expect(data.Gondoruwo.Name).To(Equal("Power"))
			Expect(data.Gondoruwo.Locatadata).To(Equal(324))
		})
	})
})
