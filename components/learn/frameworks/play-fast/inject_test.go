package play_fast

import (
	"errors"

	"github.com/amanhigh/go-fun/models/learn"
	"github.com/facebookgo/inject"
	"github.com/golobby/container/v3"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/fx"
)

var _ = Describe("Inject", func() {
	var (
		app       learn.InjectApp
		component learn.InjectComponent
		redis     learn.Redis
		err       error

		//Names
		redisName     = "RedisClient"
		dbName        = "DatabaseClient"
		appDBName     = "AppDatabaseClient"
		mockRedisName = "IAmMockRedis"
	)
	BeforeEach(func() {
		app = learn.InjectApp{}
		component = learn.InjectComponent{}
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
				&inject.Object{Value: &app},
			)
			Expect(err).To(BeNil())
		})

		It("should build App", func() {
			// Initiate Populate
			err = graph.Populate()
			Expect(err).To(BeNil())

			Expect(app.AppDB).To(Not(BeNil()), "Inject Fields")
			Expect(app.AppDB.GetDatabaseName()).To(Equal(appDBName))

			Expect(app.Component.Db).To(Not(BeNil()), "Inject Nested Fields")
			Expect(app.Component.Db.GetDatabaseName()).To(Equal(dbName))
			Expect(app.Component.Redis).To(Not(BeNil()))
			Expect(app.Component.Redis.GetRedisName()).To(Equal(redisName))

			Expect(app.NonInjectedField).To(Equal(""), "Leave Non Tagged Field")

		})

		It("should build Component", func() {
			// Build Redis Client directly
			component.Redis = &learn.RedisClient{Name: mockRedisName}

			err = graph.Provide(
				&inject.Object{Value: &component},
			)
			Expect(err).To(BeNil())

			err = graph.Populate()
			Expect(err).To(BeNil())

			Expect(component.Db).To(Not(BeNil()), "Inject Fields")
			Expect(component.Redis).To(Not(BeNil()))
			Expect(component.Db.GetDatabaseName()).To(Equal(dbName), "Inject Graph Component")
			Expect(component.Redis.GetRedisName(), "").To(Equal(mockRedisName), "Leave out Custom Component")

		})
	})

	// https://uber-go.github.io/fx/get-started/
	Context("Uber Fx", func() {
		var (
			//Only Accepts Pointer to Pointer
			uberApp = &app
			module  fx.Option
			app     *fx.App
		)
		BeforeEach(func() {
			// Create Module
			module = fx.Module("inject",
				fx.Provide(
					func() learn.Redis {
						return learn.NewRedisClient(redisName)
					},
					func() learn.Database {
						return learn.NewDatabaseClient(dbName)
					},
					fx.Annotate(func() *learn.DatabaseClient {
						return learn.NewDatabaseClient(appDBName)
					},
						fx.ResultTags(`name:"appdb"`),
						fx.As(new(learn.Database)),
					),
					learn.NewInjectComponent,
					fx.Annotate(
						learn.NewInjectApp,
						fx.ParamTags(`name:""`, `name:"appdb"`),
					),
				),
			)

			//Use Module and Generate App
			app = fx.New(
				module,
				fx.Populate(&uberApp),
			)

			Expect(app.Err()).ShouldNot(HaveOccurred())
		})

		It("should build App", func() {
			Expect(uberApp).ShouldNot(BeNil())

			Expect(uberApp.AppDB).To(Not(BeNil()), "Inject Fields")
			Expect(uberApp.AppDB.GetDatabaseName()).To(Equal(appDBName))

			Expect(uberApp.Component.Db).To(Not(BeNil()), "Inject Nested Fields")
			Expect(uberApp.Component.Db.GetDatabaseName()).To(Equal(dbName))
			Expect(uberApp.Component.Redis).To(Not(BeNil()))
			Expect(uberApp.Component.Redis.GetRedisName()).To(Equal(redisName))

			Expect(uberApp.NonInjectedField).To(Equal(""), "Leave Non Tagged Field")
		})

		It("should resolve", func() {
			err = fx.New(
				module,
				fx.Populate(&redis),
			).Err()

			Expect(err).To(BeNil())
			Expect(redis.GetRedisName()).To(Equal(redisName))
		})

		It("should override", func() {
			err = fx.New(
				module,
				fx.Decorate(
					func() learn.Redis {
						return learn.NewRedisClient(mockRedisName)
					},
				),
				fx.Populate(&redis),
			).Err()

			Expect(err).To(BeNil())
			Expect(redis.GetRedisName()).To(Equal(mockRedisName))
		})
	})

	// https://github.com/golobby/container
	Context("Golobby Container", func() {
		var (
			c = container.New()
		)

		BeforeEach(func() {
			container.MustSingleton(c, func() learn.Redis {
				return learn.NewRedisClient(redisName)
			})
			container.MustSingleton(c, func() learn.Database {
				return learn.NewDatabaseClient(dbName)
			})
			container.MustNamedSingleton(c, "AppDB", func() learn.Database {
				return learn.NewDatabaseClient(appDBName)
			})
			container.MustSingleton(c, learn.NewInjectComponent)

			//Build App
			err = c.Fill(&app)
			Expect(err).To(BeNil())
		})

		It("should build App", func() {
			Expect(app).ShouldNot(BeNil())

			Expect(app.AppDB).To(Not(BeNil()), "Inject Fields")
			Expect(app.AppDB.GetDatabaseName()).To(Equal(appDBName))

			Expect(app.Component.Db).To(Not(BeNil()), "Inject Nested Fields")
			Expect(app.Component.Db.GetDatabaseName()).To(Equal(dbName))
			Expect(app.Component.Redis).To(Not(BeNil()))
			Expect(app.Component.Redis.GetRedisName()).To(Equal(redisName))

			Expect(app.NonInjectedField).To(Equal(""), "Leave Non Tagged Field")
		})

		It("should resolve", func() {
			err = c.Resolve(&redis)
			Expect(err).To(BeNil())

			Expect(redis.GetRedisName()).To(Equal(redisName))
		})

		It("should call", func() {
			err = c.Call(func(r learn.Redis) {
				Expect(r.GetRedisName()).To(Equal(redisName))
			})
			Expect(err).To(BeNil())
		})

		It("should override", func() {
			container.MustSingleton(c, func() learn.Redis {
				return learn.NewRedisClient(mockRedisName)
			})
			err = c.Resolve(&redis)
			Expect(err).To(BeNil())

			Expect(redis.GetRedisName()).To(Equal(mockRedisName))
		})

		It("should pass error", func() {
			err = c.Singleton(func() (learn.Redis, error) {
				return nil, errors.New("oops")
			})
			Expect(err).Should(HaveOccurred())
		})

		It("should overwrite Existing Field", func() {
			var comp *learn.InjectComponent
			comp = &learn.InjectComponent{
				Redis: learn.NewRedisClient("randomRedisClient"),
			}
			err = c.Resolve(&comp)
			Expect(err).To(BeNil())
			Expect(comp.Redis.GetRedisName()).To(Equal(redisName))
		})
	})
})
