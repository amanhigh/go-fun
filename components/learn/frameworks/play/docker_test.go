package play_test

import (
	"context"
	"database/sql"

	util2 "github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/go-redis/redis/v8"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
)

var _ = Describe("Docker", Label(models.GINKGO_SETUP), func() {
	var (
		pool *dockertest.Pool
		err  error
	)
	BeforeEach(func() {
		// Reads Endpoint, CERT etc from ENV Variables
		pool, err = dockertest.NewPool("")
		Expect(err).To(BeNil())
	})

	It("should build", func() {
		Expect(pool).To(Not(BeNil()))
	})

	Context("Mysql", func() {
		var (
			password = "root"
			db       *sql.DB
			port     int
			mysqlR   *dockertest.Resource
		)
		BeforeEach(func() {
			// pulls an image, creates a container based on it and runs it
			mysqlR, err = pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=" + password})
			Expect(err).To(BeNil())

			//Wait for Startup to Happen
			// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
			err = pool.Retry(func() (err error) {
				port, err = util2.ParseInt(mysqlR.GetPort("3306/tcp"))
				Expect(err).To(BeNil())

				if db, err = util2.CreateMysqlConnection("root", password, "docker", "mysql", port); err == nil {
					err = db.Ping()
				}
				return err
			})
			Expect(err).To(BeNil())

		})

		AfterEach(func() {
			err = pool.Purge(mysqlR)
			Expect(err).To(BeNil())
		})

		It("should reachable", func() {
			Expect(db.Ping()).To(BeNil())
		})

		Context("Container", func() {
			It("mysql exists", func() {
				name, ok := pool.ContainerByName("mysql")
				Expect(ok).To(BeTrue())
				Expect(name).To(Not(BeNil()))

			})

			It("random should not exists", func() {
				name, ok := pool.ContainerByName("random")
				Expect(ok).To(BeFalse())
				Expect(name).To(BeNil())
			})
		})

		Context("Redis", func() {
			var (
				redisR *dockertest.Resource
				rdb    *redis.Client
			)

			BeforeEach(func() {
				//FIXME:Connect Network of Mysql and Redis
				redisR, err = pool.Run("bitnami/redis", "latest", []string{"ALLOW_EMPTY_PASSWORD=yes"})
				Expect(err).To(BeNil())
				err = pool.Retry(func() error {
					rdb = redis.NewClient(&redis.Options{
						Addr:     "docker:" + redisR.GetPort("6379/tcp"),
						Password: "", // no password set
						DB:       0,  // use default DB
					})
					return rdb.Ping(context.Background()).Err()
				})
				Expect(err).To(BeNil())
			})

			AfterEach(func() {
				err = pool.Purge(redisR)
				Expect(err).To(BeNil())
			})

			It("should pong", func() {
				Expect(rdb.Ping(context.Background()).String()).To(Equal("ping: PONG"))
			})
		})
	})
})
