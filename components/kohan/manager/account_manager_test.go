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
	"github.com/gocarina/gocsv"
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
				account, err := accountManager.GetRecord(ctx, testAccount.Symbol, 2024)
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
				_, err := accountManager.GetRecord(ctx, testAccount.Symbol, 2024)
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
				_, err := accountManager.GetRecord(ctx, testAccount.Symbol, 2024)
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
				_, err := accountManager.GetRecord(ctx, testAccount.Symbol, 2024)
				Expect(err).To(Equal(common.ErrInternalServerError))
			})
		})
	})

	Context("Smart Account File Detection", func() {
		Context("when accounts_2023.csv exists for 2024 tax computation", func() {
			var accounts2023Path string

			BeforeEach(func() {
				// Create accounts_2023.csv with test data
				accounts2023Path = filepath.Join(accountFilePath, "../accounts_2023.csv")
				accounts := []tax.Account{testAccount}
				file, err := os.Create(accounts2023Path)
				Expect(err).ToNot(HaveOccurred())
				defer file.Close()
				err = gocsv.MarshalFile(&accounts, file)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				os.Remove(accounts2023Path)
			})

			It("should auto-use accounts_2023.csv", func() {
				account, err := accountManager.GetRecord(ctx, testAccount.Symbol, 2024)
				Expect(err).ToNot(HaveOccurred())
				Expect(account.Symbol).To(Equal(testAccount.Symbol))
				Expect(account.Quantity).To(Equal(testAccount.Quantity))
			})
		})

		Context("when no previous year file exists", func() {
			It("should handle fresh start scenario", func() {
				_, err := accountManager.GetRecord(ctx, "NEWSTOCK", 2024)
				Expect(err).To(Equal(common.ErrNotFound))
			})
		})

		Context("when ticker not found in previous year file", func() {
			var accounts2023Path string

			BeforeEach(func() {
				// Create accounts_2023.csv with different stock
				accounts2023Path = filepath.Join(filepath.Dir(accountFilePath), "accounts_2023.csv")
				differentAccount := tax.Account{
					Symbol:      "MSFT",
					Quantity:    50,
					Cost:        5000,
					MarketValue: 5500,
				}
				accounts := []tax.Account{differentAccount}
				file, err := os.Create(accounts2023Path)
				Expect(err).ToNot(HaveOccurred())
				defer file.Close()
				err = gocsv.MarshalFile(&accounts, file)
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				os.Remove(accounts2023Path)
			})

			It("should return fresh start for that ticker", func() {
				_, err := accountManager.GetRecord(ctx, "NEWSTOCK", 2024)
				Expect(err).To(Equal(common.ErrNotFound))
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
