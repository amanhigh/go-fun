package play_fast

import (
	"github.com/amanhigh/go-fun/models/learn"
	"github.com/facebookgo/inject"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

var _ = Describe("Inject", func() {
	var (
		myApp     learn.MyApplication
		component learn.MyComponent
		err       error

		//Names
		redisName       = "RedisClient"
		dbName          = "DatabaseClient"
		appDBName       = "AppDatabaseClient"
		customRedisName = "MyRedisClient"
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
			graph inject.Graph
		)

		BeforeEach(func() {
			//Create Fresh Graph
			graph = inject.Graph{}

			// Provide Components & Build App (myApp)
			err = graph.Provide(
				&inject.Object{Value: learn.NewRedisClient(redisName)},
				&inject.Object{Value: learn.NewDatabaseClient(dbName)},
				&inject.Object{Value: learn.NewDatabaseClient(appDBName), Name: "appdb"},
				&inject.Object{Value: &myApp},
			)
			Expect(err).To(BeNil())
		})

		It("should build App", func() {
			// Initiate Populate
			err = graph.Populate()
			Expect(err).To(BeNil())

			Expect(myApp.AppDB).To(Not(BeNil()), "Inject Fields")
			Expect(myApp.AppDB.Name).To(Equal(appDBName))

			Expect(myApp.Container.Db).To(Not(BeNil()), "Inject Nested Fields")
			Expect(myApp.Container.Db.Name).To(Equal(dbName))
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

	Context("Uber Fx", func() {
		var uberApp *learn.MyApplication
		BeforeEach(func() {
			app := fx.New(
				fx.Provide(
					func() *learn.RedisClient {
						return learn.NewRedisClient(redisName)
					},
					func() *learn.DatabaseClient {
						return learn.NewDatabaseClient(dbName)
					},
					fx.Annotate(func() *learn.DatabaseClient {
						return learn.NewDatabaseClient(appDBName)
					}, fx.ResultTags(`name:"appdb"`)),
					learn.NewMyComponent,
					fx.Annotate(
						learn.NewMyApplication,
						fx.ParamTags(`name:""`, `name:"appdb"`),
					),
				),
				fx.Populate(&uberApp),
			)

			Expect(app.Err()).ShouldNot(HaveOccurred())
		})

		It("should build App", func() {
			Expect(uberApp).ShouldNot(BeNil())

			Expect(uberApp.AppDB).To(Not(BeNil()), "Inject Fields")
			Expect(uberApp.AppDB.Name).To(Equal(appDBName))

			Expect(uberApp.Container.Db).To(Not(BeNil()), "Inject Nested Fields")
			Expect(uberApp.Container.Db.Name).To(Equal(dbName))
			Expect(uberApp.Container.Redis).To(Not(BeNil()))
			Expect(uberApp.Container.Redis.Name).To(Equal(redisName))

			Expect(uberApp.NonInjectedField).To(Equal(""), "Leave Non Tagged Field")
		})
	})

})
