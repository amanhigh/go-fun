package play_test

import (
	"database/sql"
	"github.com/amanhigh/go-fun/apps/common/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
)

var _ = Describe("Docker", func() {
	var (
		pool     *dockertest.Pool
		resource *dockertest.Resource
		err      error
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
		)
		BeforeEach(func() {
			// pulls an image, creates a container based on it and runs it
			resource, err = pool.Run("mysql", "5.7", []string{"MYSQL_ROOT_PASSWORD=" + password})
			Expect(err).To(BeNil())

		})

		AfterEach(func() {
			err = pool.Purge(resource)
			Expect(err).To(BeNil())
		})

		It("should run", func() {
			// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
			err = pool.Retry(func() error {
				port, err = util.ParseInt(resource.GetPort("3306/tcp"))
				Expect(err).To(BeNil())

				db, err = util.CreateMysqlConnection("root", password, "docker", "mysql", port)
				if err != nil {
					return err
				}
				return db.Ping()
			})
			Expect(err).To(BeNil())
		})
	})
})
