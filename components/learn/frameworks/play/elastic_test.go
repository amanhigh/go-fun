package play_test

import (
	"github.com/amanhigh/go-fun/models"
	es "github.com/elastic/go-elasticsearch"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/getting-started-go.html

var _ = Describe("Elastic", Label(models.GINKGO_SETUP), func() {
	var (
		elasticClient *es.Client
		endpoint      = "http://docker:9200"
		err           error
		// ctx           = context.Background()
	)

	BeforeEach(func() {
		elasticClient, err = es.NewClient(es.Config{Addresses: []string{endpoint}})
		Expect(err).To(BeNil())
		Expect(elasticClient).ToNot(BeNil())
	})

	FIt("should connect", func() {
		info, err := elasticClient.Ping()
		Expect(err).To(BeNil())
		Expect(info).ToNot(BeNil())
	})

	Context("Document", func() {
		var (
		// indexName = "learn"
		// typeName  = "learnType"
		// docs      = []frameworks.DocumentRequest{
		// 	{Title: "Aman", Content: "Preet"},
		// 	{Title: "John", Content: "Doe"},
		// }
		// query = "aman"
		// skip  = 1
		// top   = 2
		)
		BeforeEach(func() {
			// Insert documents in bulk
			// bulk := elasticClient.Bulk().Index(indexName).Type(typeName)
			// for _, d := range docs {
			// 	doc := frameworks.Document{
			// 		ID:        shortid.MustGenerate(),
			// 		Title:     d.Title,
			// 		CreatedAt: time.Now().UTC(),
			// 		Content:   d.Content,
			// 	}
			// 	bulk.Add(elastic.NewBulkIndexRequest().Id(doc.ID).Doc(doc))
			// }
			// _, err = bulk.Do(ctx)
			// Expect(err).To(BeNil())
		})

		AfterEach(func() {
			//Delete Records
		})

		It("should search", func() {

		})

	})
})
