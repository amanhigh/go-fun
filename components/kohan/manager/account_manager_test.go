package manager_test

import (
	"context"
	"net/http"

	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository/mocks"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AccountManager", func() {
	var (
		ctx            = context.Background()
		mockRepo       *mocks.AccountRepository
		accountManager manager.AccountManager
		testAccount    tax.Account
		accountDir     = "/tmp"
		testYear       = 2024
	)

	BeforeEach(func() {
		mockRepo = mocks.NewAccountRepository(GinkgoT())
		accountManager = manager.NewAccountManager(mockRepo, accountDir)

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
					GetAllRecordsForYear(ctx, testYear).
					Return([]tax.Account{testAccount}, nil)
			})

			It("returns account details", func() {
				account, err := accountManager.GetRecord(ctx, testAccount.Symbol, testYear)
				Expect(err).ToNot(HaveOccurred())
				Expect(account).To(Equal(testAccount))
			})
		})

		Context("when no account found", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().
					GetAllRecordsForYear(ctx, testYear).
					Return([]tax.Account{}, nil)
			})

			It("returns not found error", func() {
				_, err := accountManager.GetRecord(ctx, testAccount.Symbol, testYear)
				Expect(err).To(Equal(common.ErrNotFound))
			})
		})

		Context("when multiple accounts found", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().
					GetAllRecordsForYear(ctx, testYear).
					Return([]tax.Account{testAccount, testAccount}, nil)
			})

			It("returns bad request error", func() {
				_, err := accountManager.GetRecord(ctx, testAccount.Symbol, testYear)
				Expect(err).To(HaveOccurred())
				Expect(err.Code()).To(Equal(http.StatusBadRequest))
				Expect(err.Error()).To(ContainSubstring("multiple accounts found"))
			})
		})

		Context("when repository error occurs", func() {
			BeforeEach(func() {
				mockRepo.EXPECT().
					GetAllRecordsForYear(ctx, testYear).
					Return(nil, common.ErrInternalServerError)
			})

			It("returns error from repository", func() {
				_, err := accountManager.GetRecord(ctx, testAccount.Symbol, testYear)
				Expect(err).To(Equal(common.ErrInternalServerError))
			})
		})
	})

	Context("Smart Account File Detection", func() {
		Context("when accounts_2023.csv exists for 2024 tax computation", func() {
			BeforeEach(func() {
				// Mock repository returns test data from accounts_2023.csv
				mockRepo.EXPECT().
					GetAllRecordsForYear(ctx, testYear).
					Return([]tax.Account{testAccount}, nil)
			})

			It("should auto-use accounts_2023.csv", func() {
				account, err := accountManager.GetRecord(ctx, testAccount.Symbol, testYear)
				Expect(err).ToNot(HaveOccurred())
				Expect(account.Symbol).To(Equal(testAccount.Symbol))
				Expect(account.Quantity).To(Equal(testAccount.Quantity))
			})
		})

		Context("when no previous year file exists", func() {
			BeforeEach(func() {
				// Mock repository returns not found error
				mockRepo.EXPECT().
					GetAllRecordsForYear(ctx, testYear).
					Return(nil, common.ErrNotFound)
			})

			It("should handle fresh start scenario", func() {
				_, err := accountManager.GetRecord(ctx, "NEWSTOCK", testYear)
				Expect(err).To(Equal(common.ErrNotFound))
			})
		})

		Context("when ticker not found in previous year file", func() {
			BeforeEach(func() {
				// Mock repository returns different account
				differentAccount := tax.Account{
					Symbol:      "MSFT",
					Quantity:    50,
					Cost:        5000,
					MarketValue: 5500,
				}
				mockRepo.EXPECT().
					GetAllRecordsForYear(ctx, testYear).
					Return([]tax.Account{differentAccount}, nil)
			})

			It("should return fresh start for that ticker", func() {
				_, err := accountManager.GetRecord(ctx, "NEWSTOCK", testYear)
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

			// Set up mock expectation for SaveYearEndAccounts
			expectedAccounts := []tax.Account{
				{
					Symbol:      "GOOG",
					Quantity:    10,
					Cost:        1500, // 10 * 150
					MarketValue: 1500, // 10 * 150
				},
			}
			mockRepo.EXPECT().
				SaveYearEndAccounts(ctx, 2023, expectedAccounts).
				Return(nil)

			err := accountManager.GenerateYearEndAccounts(ctx, 2023, valuations)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
