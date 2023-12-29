package play_test

import (
	"context"
	"time"

	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/learn/frameworks"
	"github.com/olivere/elastic"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/teris-io/shortid"
)

var _ = Describe("Elastic", Label(models.GINKGO_SETUP), func() {
	var (
		elasticClient *elastic.Client
		endpoint      = "http://docker:9200"
		err           error
		ctx           = context.Background()
	)

	BeforeEach(func() {
		elasticClient, err = elastic.NewClient(
			elastic.SetURL(endpoint),
			elastic.SetSniff(false),
		)
		Expect(err).To(BeNil())
		Expect(elasticClient).ToNot(BeNil())
	})

	It("should connect", func() {
		info, code, err := elasticClient.Ping(endpoint).Do(ctx)
		Expect(err).To(BeNil())
		Expect(code).To(Equal(200))
		Expect(info).ToNot(BeNil())
	})

	Context("Document", func() {
		var (
			indexName = "learn"
			typeName  = "learnType"
			docs      = []frameworks.DocumentRequest{
				{Title: "Aman", Content: "Preet"},
				{Title: "John", Content: "Doe"},
			}
			query = "aman"
			skip  = 1
			top   = 2
		)
		BeforeEach(func() {
			// Insert documents in bulk
			bulk := elasticClient.Bulk().Index(indexName).Type(typeName)
			for _, d := range docs {
				doc := frameworks.Document{
					ID:        shortid.MustGenerate(),
					Title:     d.Title,
					CreatedAt: time.Now().UTC(),
					Content:   d.Content,
				}
				bulk.Add(elastic.NewBulkIndexRequest().Id(doc.ID).Doc(doc))
			}
			_, err = bulk.Do(ctx)
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			//Delete Records
			_, err = elasticClient.DeleteByQuery(indexName).Do(ctx)
			Expect(err).To(BeNil())
		})

		It("should search", func() {
			esQuery := elastic.NewMultiMatchQuery(query, "title", "content").
				Fuzziness("2").MinimumShouldMatch("2")

			result, err := elasticClient.Search().
				Index(indexName).
				Query(esQuery).
				From(skip).Size(top).
				Do(ctx)

			Expect(err).To(BeNil())
			Expect(result.Hits.TotalHits).To(Equal(int64(1)))
		})

	})
})
