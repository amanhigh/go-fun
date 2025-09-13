package repository_test

import (
	"context"
	"os"
	"path/filepath"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/common"
	"github.com/amanhigh/go-fun/models/tax"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testAccountsCSV = `Symbol,Quantity,Cost,MarketValue
AAPL,100,15050.00,16000.00
GOOGL,50,125000.00,130000.00
MSFT,75,22500.00,24000.00`

var _ = Describe("AccountRepository", func() {
	var (
		accountRepo  repository.AccountRepository
		testDir      string
		ctx          = context.Background()
		testAccounts []tax.Account
		err          error
	)

	BeforeEach(func() {
		// Create temp directory
		testDir, err = os.MkdirTemp("", "account-test-*")
		Expect(err).NotTo(HaveOccurred())

		accountRepo = repository.NewAccountRepository(testDir)

		// Setup test data
		testAccounts = []tax.Account{
			{Symbol: "AAPL", Quantity: 100, Cost: 15050.00, MarketValue: 16000.00},
			{Symbol: "GOOGL", Quantity: 50, Cost: 125000.00, MarketValue: 130000.00},
			{Symbol: "MSFT", Quantity: 75, Cost: 22500.00, MarketValue: 24000.00},
		}
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
	})

	Context("SaveYearEndAccounts", func() {
		It("should save accounts successfully", func() {
			httpErr := accountRepo.SaveYearEndAccounts(ctx, 2024, testAccounts)
			Expect(httpErr).ToNot(HaveOccurred())

			// Verify file was created
			expectedPath := filepath.Join(testDir, "accounts_2024.csv")
			Expect(expectedPath).To(BeAnExistingFile())

			// Verify file content
			content, readErr := os.ReadFile(expectedPath)
			Expect(readErr).ToNot(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("AAPL,100,15050,16000"))
			Expect(string(content)).To(ContainSubstring("GOOGL,50,125000,130000"))
			Expect(string(content)).To(ContainSubstring("MSFT,75,22500,24000"))
		})

		It("should handle empty accounts slice", func() {
			httpErr := accountRepo.SaveYearEndAccounts(ctx, 2024, []tax.Account{})
			Expect(httpErr).ToNot(HaveOccurred())

			expectedPath := filepath.Join(testDir, "accounts_2024.csv")
			Expect(expectedPath).To(BeAnExistingFile())
		})

		It("should overwrite existing file", func() {
			// First save
			httpErr := accountRepo.SaveYearEndAccounts(ctx, 2024, testAccounts)
			Expect(httpErr).ToNot(HaveOccurred())

			// Second save with different data
			newAccounts := []tax.Account{
				{Symbol: "TSLA", Quantity: 25, Cost: 50000.00, MarketValue: 55000.00},
			}
			httpErr = accountRepo.SaveYearEndAccounts(ctx, 2024, newAccounts)
			Expect(httpErr).ToNot(HaveOccurred())

			// Verify only new data exists
			expectedPath := filepath.Join(testDir, "accounts_2024.csv")
			content, readErr := os.ReadFile(expectedPath)
			Expect(readErr).ToNot(HaveOccurred())
			Expect(string(content)).To(ContainSubstring("TSLA,25,50000,55000"))
			Expect(string(content)).ToNot(ContainSubstring("AAPL"))
		})

		Context("Error Cases", func() {
			It("should handle invalid directory", func() {
				invalidRepo := repository.NewAccountRepository("/invalid/path")
				httpErr := invalidRepo.SaveYearEndAccounts(ctx, 2024, testAccounts)
				Expect(httpErr).To(HaveOccurred())
				Expect(httpErr.Code()).To(Equal(500))
			})

			It("should handle read-only directory", func() {
				// Make directory read-only
				chmodErr := os.Chmod(testDir, 0444)
				Expect(chmodErr).ToNot(HaveOccurred())

				httpErr := accountRepo.SaveYearEndAccounts(ctx, 2024, testAccounts)
				Expect(httpErr).To(HaveOccurred())
				Expect(httpErr.Code()).To(Equal(500))

				// Restore permissions for cleanup
				restoreErr := os.Chmod(testDir, 0755)
				Expect(restoreErr).ToNot(HaveOccurred())
			})
		})
	})

	Context("GetAllRecordsForYear", func() {
		BeforeEach(func() {
			// Create test file
			accountsFile := filepath.Join(testDir, "accounts_2024.csv")
			err = os.WriteFile(accountsFile, []byte(testAccountsCSV), util.DEFAULT_PERM)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should read accounts successfully", func() {
			accounts, httpErr := accountRepo.GetAllRecordsForYear(ctx, 2024)
			Expect(httpErr).ToNot(HaveOccurred())
			Expect(accounts).To(HaveLen(3))

			// Verify first account
			first := accounts[0]
			Expect(first.Symbol).To(Equal("AAPL"))
			Expect(first.Quantity).To(Equal(100.0))
			Expect(first.Cost).To(Equal(15050.00))
			Expect(first.MarketValue).To(Equal(16000.00))

			// Verify second account
			second := accounts[1]
			Expect(second.Symbol).To(Equal("GOOGL"))
			Expect(second.Quantity).To(Equal(50.0))
			Expect(second.Cost).To(Equal(125000.00))
			Expect(second.MarketValue).To(Equal(130000.00))
		})

		Context("Error Cases", func() {
			It("should return NotFound for missing file", func() {
				accounts, httpErr := accountRepo.GetAllRecordsForYear(ctx, 2023)
				Expect(accounts).To(BeNil())
				Expect(httpErr).To(Equal(common.ErrNotFound))
			})

			It("should handle malformed CSV", func() {
				malformedFile := filepath.Join(testDir, "accounts_2025.csv")
				err := os.WriteFile(malformedFile, []byte("invalid,csv,data"), util.DEFAULT_PERM)
				Expect(err).ToNot(HaveOccurred())

				accounts, httpErr := accountRepo.GetAllRecordsForYear(ctx, 2025)
				Expect(accounts).To(BeNil())
				Expect(httpErr).To(HaveOccurred())
				// gocsv might parse this as empty, returning ErrNotFound instead of server error
				Expect(httpErr.Code()).To(BeElementOf(404, 500))
			})

			It("should handle empty CSV file", func() {
				emptyFile := filepath.Join(testDir, "accounts_2026.csv")
				err := os.WriteFile(emptyFile, []byte(""), util.DEFAULT_PERM)
				Expect(err).ToNot(HaveOccurred())

				accounts, httpErr := accountRepo.GetAllRecordsForYear(ctx, 2026)
				Expect(accounts).To(BeNil())
				Expect(httpErr).To(HaveOccurred())
				// Empty file causes parsing error, returns server error not NotFound
				Expect(httpErr.Code()).To(Equal(500))
			})

			It("should handle CSV with only headers", func() {
				headerOnlyFile := filepath.Join(testDir, "accounts_2027.csv")
				err := os.WriteFile(headerOnlyFile, []byte("Symbol,Quantity,Cost,MarketValue"), util.DEFAULT_PERM)
				Expect(err).ToNot(HaveOccurred())

				accounts, httpErr := accountRepo.GetAllRecordsForYear(ctx, 2027)
				Expect(accounts).To(BeNil())
				Expect(httpErr).To(Equal(common.ErrNotFound))
			})

			It("should handle file access permission issues", func() {
				restrictedFile := filepath.Join(testDir, "accounts_2028.csv")
				err := os.WriteFile(restrictedFile, []byte(testAccountsCSV), util.DEFAULT_PERM)
				Expect(err).ToNot(HaveOccurred())

				// Remove read permissions
				err = os.Chmod(restrictedFile, 0000)
				Expect(err).ToNot(HaveOccurred())

				accounts, httpErr := accountRepo.GetAllRecordsForYear(ctx, 2028)
				Expect(accounts).To(BeNil())
				Expect(httpErr).To(HaveOccurred())
				Expect(httpErr.Code()).To(Equal(500))

				// Restore permissions for cleanup
				restoreErr := os.Chmod(restrictedFile, util.DEFAULT_PERM)
				Expect(restoreErr).ToNot(HaveOccurred())
			})
		})
	})

	Context("Integration - Save and Read", func() {
		It("should save and read back same accounts", func() {
			// Save accounts
			httpErr := accountRepo.SaveYearEndAccounts(ctx, 2024, testAccounts)
			Expect(httpErr).ToNot(HaveOccurred())

			// Read accounts back
			readAccounts, httpErr := accountRepo.GetAllRecordsForYear(ctx, 2024)
			Expect(httpErr).ToNot(HaveOccurred())
			Expect(readAccounts).To(HaveLen(len(testAccounts)))

			// Verify data integrity
			for i, expected := range testAccounts {
				actual := readAccounts[i]
				Expect(actual.Symbol).To(Equal(expected.Symbol))
				Expect(actual.Quantity).To(Equal(expected.Quantity))
				Expect(actual.Cost).To(Equal(expected.Cost))
				Expect(actual.MarketValue).To(Equal(expected.MarketValue))
			}
		})

		It("should handle multiple years independently", func() {
			accounts2023 := []tax.Account{
				{Symbol: "NVDA", Quantity: 30, Cost: 45000.00, MarketValue: 50000.00},
			}
			accounts2024 := []tax.Account{
				{Symbol: "AMD", Quantity: 200, Cost: 30000.00, MarketValue: 32000.00},
			}

			// Save different data for different years
			httpErr := accountRepo.SaveYearEndAccounts(ctx, 2023, accounts2023)
			Expect(httpErr).ToNot(HaveOccurred())

			httpErr = accountRepo.SaveYearEndAccounts(ctx, 2024, accounts2024)
			Expect(httpErr).ToNot(HaveOccurred())

			// Verify each year has correct data
			read2023, httpErr := accountRepo.GetAllRecordsForYear(ctx, 2023)
			Expect(httpErr).ToNot(HaveOccurred())
			Expect(read2023).To(HaveLen(1))
			Expect(read2023[0].Symbol).To(Equal("NVDA"))

			read2024, httpErr := accountRepo.GetAllRecordsForYear(ctx, 2024)
			Expect(httpErr).ToNot(HaveOccurred())
			Expect(read2024).To(HaveLen(1))
			Expect(read2024[0].Symbol).To(Equal("AMD"))
		})
	})
})
