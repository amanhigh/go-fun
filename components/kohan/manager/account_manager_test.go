package manager_test

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AccountManager", func() {
	var (
		ctx             = context.Background()
		mockRepo        *mocks.AccountRepository
		accountManager  manager.AccountManager
		testAccount     tax.Account
		accountFilePath = "/tmp/accounts.csv"
	)

	BeforeEach(func() {
		mockRepo = mocks.NewAccountRepository(GinkgoT())
		accountManager = manager.NewAccountManager(mockRepo, accountFilePath)

		// Setup test account
		testAccount = tax.Account{
			Symbol:      "AAPL",
			Quantity:    100,
			Cost:        3833,
			MarketValue: 4201,
		}
	})

	Context("GetRecord", func() {
		Context("when single account exists", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().
					GetRecordsForTicker(ctx, testAccount.Symbol).
					Return([]tax.Account{testAccount}, nil)
			})

			It("returns account details", func() {
				account, err := accountManager.GetRecord(ctx, testAccount.Symbol)
				Expect(err).ToNot(HaveOccurred())
				Expect(account).To(Equal(testAccount))
			})
		})

		Context("when no account found", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().
					GetRecordsForTicker(ctx, testAccount.Symbol).
					Return([]tax.Account{}, nil)
			})

			It("returns not found error", func() {
				_, err := accountManager.GetRecord(ctx, testAccount.Symbol)
				Expect(err).To(Equal(common.ErrNotFound))
			})
		})

		Context("when multiple accounts found", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().
					GetRecordsForTicker(ctx, testAccount.Symbol).
					Return([]tax.Account{testAccount, testAccount}, nil)
			})

			It("returns bad request error", func() {
				_, err := accountManager.GetRecord(ctx, testAccount.Symbol)
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusBadRequest))
				Expect(err.Error()).To(ContainSubstring("multiple accounts found"))
			})
		})

		Context("when repository error occurs", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().
					GetRecordsForTicker(ctx, testAccount.Symbol).
					Return(nil, common.ErrInternalServerError)
			})

			It("returns error from repository", func() {
				_, err := accountManager.GetRecord(ctx, testAccount.Symbol)
				Expect(err).To(Equal(common.ErrInternalServerError))
			})
		})
	})

	Context("GenerateYearEndAccounts", func() {
		It("should generate year end accounts", func() {
			valuations := []tax.Valuation{
				{
					Ticker: "GOOG",
					YearEndPosition: tax.Position{
						Quantity: 10,
						USDPrice: 150,
					},
				},
			}
			err := accountManager.GenerateYearEndAccounts(ctx, 2023, valuations)
			Expect(err).ToNot(HaveOccurred())

			// Verify file was created
			fileName := "accounts_2023.csv"
			filePath := filepath.Join(filepath.Dir(accountFilePath), fileName)
			_, statErr := os.Stat(filePath)
			Expect(statErr).ToNot(HaveOccurred())

			// Clean up
			os.Remove(filePath)
		})
	})
})
