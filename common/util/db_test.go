// Test file with similar patterns for different DB types
//
//nolint:dupl // False positives: Similar test patterns for MySQL/PostgreSQL containers
package util_test

import (
	"context"
	"embed"
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/glebarez/sqlite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed testdata/migrations/*.sql
var migrationFS embed.FS

var _ = Describe("DBUtil", Ordered, Label(models.GINKGO_SLOW), func() {
	var (
		ctx               = context.Background()
		db                *gorm.DB
		err               error
		mysqlContainer    testcontainers.Container
		postgresContainer testcontainers.Container
	)

	Context("RunMigrations", func() {
		Context("SQLite", func() {
			const migrationDir = "testdata/migrations"

			BeforeEach(func() {
				db, err = gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			})

			It("should run migrations successfully", func() {
				err = util.RunMigrations(db, migrationFS, migrationDir)
				Expect(err).ToNot(HaveOccurred())

				// Verify table was created
				Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
			})

			It("should be idempotent", func() {
				// Run migrations twice
				err = util.RunMigrations(db, migrationFS, migrationDir)
				Expect(err).ToNot(HaveOccurred())
				err = util.RunMigrations(db, migrationFS, migrationDir)
				Expect(err).ToNot(HaveOccurred())

				// Verify table still exists
				Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
			})
		})

		Context("MySQL", func() {
			const migrationDir = "testdata/migrations"

			BeforeEach(func() {
				mysqlContainer, err = util.MysqlTestContainer(ctx)
				Expect(err).ToNot(HaveOccurred())
				DeferCleanup(func() {
					Expect(testcontainers.TerminateContainer(mysqlContainer)).To(Succeed())
				})

				mysqlHost, hostErr := mysqlContainer.PortEndpoint(ctx, "3306/tcp", "")
				Expect(hostErr).ToNot(HaveOccurred())
				dbUrl := fmt.Sprintf("aman:aman@tcp(%s)/compute", mysqlHost)

				db, err = gorm.Open(mysql.Open(dbUrl), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			})

			It("should run migrations successfully", func() {
				err = util.RunMigrations(db, migrationFS, migrationDir)
				Expect(err).ToNot(HaveOccurred())

				// Verify table was created
				Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
			})

			It("should be idempotent", func() {
				// Run migrations twice
				err = util.RunMigrations(db, migrationFS, migrationDir)
				Expect(err).ToNot(HaveOccurred())
				err = util.RunMigrations(db, migrationFS, migrationDir)
				Expect(err).ToNot(HaveOccurred())

				// Verify table still exists
				Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
			})
		})

		Context("PostgreSQL", func() {
			const migrationDir = "testdata/migrations"

			BeforeEach(func() {
				postgresContainer, err = util.PostgresTestContainer(ctx)
				Expect(err).ToNot(HaveOccurred())
				DeferCleanup(func() {
					Expect(testcontainers.TerminateContainer(postgresContainer)).To(Succeed())
				})

				postgresHost, hostErr := postgresContainer.PortEndpoint(ctx, "5432/tcp", "")
				Expect(hostErr).ToNot(HaveOccurred())
				dbUrl := fmt.Sprintf("postgres://test:test@%s/testdb?sslmode=disable", postgresHost)

				db, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{
					Logger: logger.Default.LogMode(logger.Silent),
				})
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				sqlDB, _ := db.DB()
				sqlDB.Close()
			})

			It("should run migrations successfully", func() {
				err = util.RunMigrations(db, migrationFS, migrationDir)
				Expect(err).ToNot(HaveOccurred())

				// Verify table was created
				Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
			})

			It("should be idempotent", func() {
				// Run migrations twice
				err = util.RunMigrations(db, migrationFS, migrationDir)
				Expect(err).ToNot(HaveOccurred())
				err = util.RunMigrations(db, migrationFS, migrationDir)
				Expect(err).ToNot(HaveOccurred())

				// Verify table still exists
				Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
			})
		})

		// Error cases
		Context("Error handling", func() {
			const migrationDir = "testdata/migrations/sqlite"
			It("should panic with nil db", func() {
				Expect(func() {
					err := util.RunMigrations(nil, migrationFS, migrationDir)
					Expect(err).To(HaveOccurred())
				}).To(Panic())
			})
		})
	})

	Context("GormErrorMapper", func() {
		It("should return nil for nil error", func() {
			httpErr := util.GormErrorMapper(nil)
			Expect(httpErr).ToNot(HaveOccurred())
		})

		It("should return ErrNotFound for gorm.ErrRecordNotFound", func() {
			httpErr := util.GormErrorMapper(gorm.ErrRecordNotFound)
			Expect(httpErr).To(HaveOccurred())
			Expect(httpErr.Code()).To(Equal(404))
			Expect(httpErr.Error()).To(Equal("NotFound"))
		})

		It("should return ErrEntityExists for SQLite constraint failed", func() {
			sqliteErr := fmt.Errorf("constraint failed: UNIQUE constraint failed: test_table.test_column")
			httpErr := util.GormErrorMapper(sqliteErr)
			Expect(httpErr).To(HaveOccurred())
			Expect(httpErr.Code()).To(Equal(409))
			Expect(httpErr.Error()).To(Equal("EntityExists"))
		})

		It("should return ErrEntityExists for MySQL UNIQUE constraint", func() {
			mysqlErr := fmt.Errorf("Error 1062: Duplicate entry 'test' for key 'uq_test_column' - UNIQUE constraint failed")
			httpErr := util.GormErrorMapper(mysqlErr)
			Expect(httpErr).To(HaveOccurred())
			Expect(httpErr.Code()).To(Equal(409))
			Expect(httpErr.Error()).To(Equal("EntityExists"))
		})

		It("should return ErrEntityExists for PostgreSQL unique constraint", func() {
			postgresErr := fmt.Errorf("pq: duplicate key value violates unique constraint \"test_table_pkey\" - UNIQUE constraint failed")
			httpErr := util.GormErrorMapper(postgresErr)
			Expect(httpErr).To(HaveOccurred())
			Expect(httpErr.Code()).To(Equal(409))
			Expect(httpErr.Error()).To(Equal("EntityExists"))
		})

		It("should return server error for unknown errors", func() {
			unknownErr := fmt.Errorf("some random error")
			httpErr := util.GormErrorMapper(unknownErr)
			Expect(httpErr).To(HaveOccurred())
			Expect(httpErr.Code()).To(Equal(500))
			Expect(httpErr.Error()).To(Equal("some random error"))
		})
	})
})
