//nolint:dupl // Test file with similar patterns for different DB types
package util_test

import (
	"context"
	"embed"
	"fmt"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/config"
	"github.com/glebarez/sqlite"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//go:embed testdata/migrations/*/*.sql
var migrationFS embed.FS

var _ = Describe("RunMigrations", Ordered, Label(models.GINKGO_SLOW), func() {
	var (
		ctx = context.Background()
		db  *gorm.DB
		err error
	)

	// SQLite tests (fast, in-memory)
	Context("SQLite", func() {
		const migrationDir = "testdata/migrations/sqlite"
		BeforeEach(func() {
			db, err = gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			sqlDB, _ := db.DB()
			if sqlDB != nil {
				sqlDB.Close()
			}
		})

		It("should apply migrations successfully", func() {
			err = util.RunMigrations(db, migrationFS, migrationDir)
			Expect(err).ToNot(HaveOccurred())

			// Verify table was created
			Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
		})

		It("should be idempotent - running migrations twice should not error", func() {
			err = util.RunMigrations(db, migrationFS, migrationDir)
			Expect(err).ToNot(HaveOccurred())

			// Run again - should not error
			err = util.RunMigrations(db, migrationFS, migrationDir)
			Expect(err).ToNot(HaveOccurred())

			// Verify table still exists
			Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
		})

		It("should fail with empty migration directory", func() {
			err = util.RunMigrations(db, migrationFS, "")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("migration directory cannot be empty"))
		})
	})

	// MySQL tests (testcontainers)
	Context("MySQL", func() {
		var mysqlContainer testcontainers.Container
		const migrationDir = "testdata/migrations/mysql"

		BeforeAll(func() {
			mysqlContainer, err = util.MysqlTestContainer(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterAll(func() {
			if mysqlContainer != nil {
				err = mysqlContainer.Terminate(ctx)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		BeforeEach(func() {
			host, err := mysqlContainer.Host(ctx)
			Expect(err).ToNot(HaveOccurred())

			port, err := mysqlContainer.MappedPort(ctx, "3306")
			Expect(err).ToNot(HaveOccurred())

			dsn := fmt.Sprintf("aman:aman@tcp(%s:%s)/compute?charset=utf8mb4&parseTime=True&loc=Local", host, port.Port())
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			if db != nil {
				sqlDB, _ := db.DB()
				if sqlDB != nil {
					sqlDB.Close()
				}
			}
		})

		It("should apply migrations successfully", func() {
			err = util.RunMigrations(db, migrationFS, migrationDir)
			Expect(err).ToNot(HaveOccurred())

			// Verify table was created
			Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
		})

		It("should be idempotent - running migrations twice should not error", func() {
			err = util.RunMigrations(db, migrationFS, migrationDir)
			Expect(err).ToNot(HaveOccurred())

			// Run again - should not error
			err = util.RunMigrations(db, migrationFS, migrationDir)
			Expect(err).ToNot(HaveOccurred())

			// Verify table still exists
			Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
		})
	})

	// PostgreSQL tests (testcontainers)
	Context("PostgreSQL", func() {
		var postgresContainer testcontainers.Container
		const migrationDir = "testdata/migrations/postgres"

		BeforeAll(func() {
			postgresContainer, err = util.PostgresTestContainer(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		AfterAll(func() {
			if postgresContainer != nil {
				err = postgresContainer.Terminate(ctx)
				Expect(err).ToNot(HaveOccurred())
			}
		})

		BeforeEach(func() {
			host, err := postgresContainer.Host(ctx)
			Expect(err).ToNot(HaveOccurred())

			port, err := postgresContainer.MappedPort(ctx, "5432")
			Expect(err).ToNot(HaveOccurred())

			dsn := fmt.Sprintf("host=%s user=test password=test dbname=testdb port=%s sslmode=disable", host, port.Port())
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			if db != nil {
				sqlDB, _ := db.DB()
				if sqlDB != nil {
					sqlDB.Close()
				}
			}
		})

		It("should apply migrations successfully", func() {
			err = util.RunMigrations(db, migrationFS, migrationDir)
			Expect(err).ToNot(HaveOccurred())

			// Verify table was created
			Expect(db.Migrator().HasTable("test_users")).To(BeTrue())
		})

		It("should be idempotent - running migrations twice should not error", func() {
			err = util.RunMigrations(db, migrationFS, migrationDir)
			Expect(err).ToNot(HaveOccurred())

			// Run again - should not error
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

// Helper to connect to test DBs
func connectToTestDB(cfg config.Db) (*gorm.DB, error) {
	switch cfg.DbType {
	case models.SQLITE:
		return gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	case models.MYSQL:
		return gorm.Open(mysql.Open(cfg.Url), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	case models.POSTGRES:
		return gorm.Open(postgres.Open(cfg.Url), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
	default:
		return nil, fmt.Errorf("unsupported db type: %s", cfg.DbType)
	}
}
