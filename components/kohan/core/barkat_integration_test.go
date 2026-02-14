package core_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/components/kohan/core"
	"github.com/amanhigh/go-fun/components/kohan/handler"
	"github.com/amanhigh/go-fun/components/kohan/manager"
	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const testPort = 19020

var baseURL string

var _ = Describe("Barkat Integration Test", func() {
	BeforeEach(func() {
		// Start server once for the suite
		if baseURL == "" {
			db, err := core.SetupBarkatDB("file::memory:?cache=shared")
			Expect(err).ToNot(HaveOccurred())
			Expect(db).ToNot(BeNil())

			repo := repository.NewJournalRepository(db)
			mgr := manager.NewJournalManager(repo)
			journalHandler := handler.NewJournalHandler(mgr)
			server := core.NewKohanServer("", nil, journalHandler)
			baseURL = fmt.Sprintf("http://localhost:%d", testPort)

			go func() {
				defer GinkgoRecover()
				_ = server.Start(testPort, util.NewGracefulShutdown())
			}()

			// Wait for server to be ready
			Eventually(func() error {
				_, err := http.Get(baseURL + "/api/v1/journal-entries")
				return err
			}, 5*time.Second, 100*time.Millisecond).Should(Succeed())
		}
	})

	Context("Create", func() {
		var (
			createdEntry barkat.Entry
		)

		BeforeEach(func() {
			entry := barkat.Entry{
				Ticker:   "RELIANCE",
				Sequence: "mwd",
				Type:     "rejected",
				Outcome:  "fail",
				Trend:    "trend",
				Images: []barkat.Image{
					{Position: 1, Path: "assets/trading/2024/01/RELIANCE.mwd.trend.rejected__20240115__100000.png"},
					{Position: 2, Path: "assets/trading/2024/01/RELIANCE.mwd.trend.rejected__20240115__100001.png"},
				},
			}

			body, err := json.Marshal(entry)
			Expect(err).ToNot(HaveOccurred())

			resp, err := http.Post(baseURL+"/api/v1/journal-entries", "application/json", bytes.NewReader(body))
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusCreated))

			respBody, err := io.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())
			resp.Body.Close()

			err = json.Unmarshal(respBody, &createdEntry)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should create entry with generated ID", func() {
			Expect(createdEntry.ID).ToNot(BeEmpty())
			Expect(createdEntry.Ticker).To(Equal("RELIANCE"))
			Expect(createdEntry.Sequence).To(Equal("mwd"))
			Expect(createdEntry.Images).To(HaveLen(2))
		})

		Context("Get", func() {
			var fetchedEntry barkat.Entry

			BeforeEach(func() {
				resp, err := http.Get(baseURL + "/api/v1/journal-entries/" + createdEntry.ID)
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				respBody, err := io.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				resp.Body.Close()

				err = json.Unmarshal(respBody, &fetchedEntry)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should get entry by ID with images", func() {
				Expect(fetchedEntry.ID).To(Equal(createdEntry.ID))
				Expect(fetchedEntry.Ticker).To(Equal("RELIANCE"))
				Expect(fetchedEntry.Images).To(HaveLen(2))
				Expect(fetchedEntry.Images[0].Path).To(ContainSubstring("RELIANCE"))
			})
		})

		Context("List", func() {
			var entryList barkat.EntryList

			BeforeEach(func() {
				resp, err := http.Get(baseURL + "/api/v1/journal-entries?limit=10")
				Expect(err).ToNot(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				respBody, err := io.ReadAll(resp.Body)
				Expect(err).ToNot(HaveOccurred())
				resp.Body.Close()

				err = json.Unmarshal(respBody, &entryList)
				Expect(err).ToNot(HaveOccurred())
			})

			It("should list entries with metadata", func() {
				Expect(entryList.Records).ToNot(BeEmpty())
				Expect(entryList.Metadata.Total).To(BeNumerically(">=", 1))
			})
		})
	})

	Context("Get Not Found", func() {
		It("should return 404 for missing entry", func() {
			resp, err := http.Get(baseURL + "/api/v1/journal-entries/nonexistent-id")
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			resp.Body.Close()
		})
	})
})
