package util_test

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	util2 "github.com/amanhigh/go-fun/common/util"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net"
	"time"
)

var _ = Describe("DbResolver", func() {
	var (
		interval   = time.Millisecond * 50
		pingTable  = "test"
		err        error
		connErr    = net.InvalidAddrError("Test")
		policy     *util2.FallBackPolicy
		db         *sql.DB
		gormDB     *gorm.DB
		mock       sqlmock.Sqlmock
		queryRegex = fmt.Sprintf("SELECT count.*%s.*", pingTable)
	)
	BeforeEach(func() {
		/* Mock DB */
		db, mock, err = sqlmock.New()
		Expect(err).To(BeNil())

		/* Gorm With Gorm DB */
		gormDB, err = gorm.Open(mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
		Expect(err).To(BeNil())

		policy = util2.NewFallBackPolicy(gormDB, interval, pingTable)
	})

	AfterEach(func() {
		defer db.Close()
	})

	It("should build", func() {
		Expect(policy).To(Not(BeNil()))
	})

	Context("Default Pool", func() {
		It("should be PRIMARY", func() {
			Expect(policy.GetPool()).To(Equal(util2.POOL_PRIMARY))
		})

		Context("On Valid Error", func() {
			BeforeEach(func() {
				policy.ReportError(connErr)
				mock.ExpectQuery(queryRegex).WillReturnError(connErr)
			})

			It("should be FALLBACK", func() {
				Consistently(policy.GetPool).Should(Equal(util2.POOL_FALLBACK))
			})

			Context("Post Recover", func() {
				BeforeEach(func() {
					mock.ExpectQuery(queryRegex).WillReturnRows(sqlmock.NewRows([]string{"5"}))
				})

				It("should be PRIMARY", func() {
					Eventually(policy.GetPool).Should(Equal(util2.POOL_PRIMARY))
				})
			})
		})

		Context("Should Remain Primary", func() {
			AfterEach(func() {
				Consistently(policy.GetPool).Should(Equal(util2.POOL_PRIMARY))
			})

			It("on nil error", func() {
				policy.ReportError(nil)
			})

			It("on non net error", func() {
				policy.ReportError(errors.New("Some Error"))
			})
		})
	})
})
