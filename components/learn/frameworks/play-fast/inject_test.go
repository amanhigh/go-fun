package play_fast

import (
	"github.com/amanhigh/go-fun/models/learn"
	"github.com/facebookgo/inject"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Inject", func() {
	var (
		myApp     learn.MyApplication
		component learn.MyComponent
		err       error
	)
	BeforeEach(func() {
		myApp = learn.MyApplication{}
		component = learn.MyComponent{}
	})

	/**
	https://scene-si.org/2016/06/16/dependency-injection-patterns-in-go/
	https://github.com/facebookgo/inject/blob/master/example_test.go
	*/
	Context("Facebook Inject", func() {
		var (
			graph           inject.Graph
			redisName       = "RedisClient"
			dbName          = "DatabaseClient"
			anonName        = "AnonDatabaseClient"
			customRedisName = "MyRedisClient"
		)

		BeforeEach(func() {
			//Create Fresh Graph
			graph = inject.Graph{}

			// Provide Components & Build App (myApp)
			err = graph.Provide(
				&inject.Object{Value: &learn.RedisClient{redisName}},
				&inject.Object{Value: &learn.DatabaseClient{dbName}},
				&inject.Object{Value: &learn.DatabaseClient{anonName}, Name: "anon"},
				&inject.Object{Value: &myApp},
			)
			Expect(err).To(BeNil())
		})

		It("should build App", func() {
			// Initiate Populate
			err = graph.Populate()
			Expect(err).To(BeNil())

			Expect(myApp.Anon).To(Not(BeNil()), "Inject Fields")
			Expect(myApp.Anon.Name).To(Equal(anonName))

			Expect(myApp.Container.Db).To(Not(BeNil()), "Inject Nested Fields")
			Expect(myApp.Container.Redis).To(Not(BeNil()))
			Expect(myApp.Container.Redis.Name).To(Equal(redisName))

			Expect(myApp.NonInjectedField).To(Equal(""), "Leave Non Tagged Field")

		})

		It("should build Component", func() {
			// Build Redis Client directly
			component.Redis = &learn.RedisClient{Name: customRedisName}

			err = graph.Provide(
				&inject.Object{Value: &component},
			)
			Expect(err).To(BeNil())

			err = graph.Populate()
			Expect(err).To(BeNil())

			Expect(component.Db).To(Not(BeNil()), "Inject Fields")
			Expect(component.Redis).To(Not(BeNil()))
			Expect(component.Db.Name).To(Equal(dbName), "Inject Graph Component")
			Expect(component.Redis.Name, "").To(Equal(customRedisName), "Leave out Cutom Component")

		})
	})

})
