package play_test

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/learn/frameworks"
	es "github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/teris-io/shortid"
)

/*
 Guide - https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/getting-started-go.html
 Store Example - https://github.com/elastic/go-elasticsearch/blob/main/_examples/xkcdsearch/store.go
 Index - http://docker:9200/learn/_search?pretty
*/

var _ = Describe("Elastic", Label(models.GINKGO_SETUP), func() {
	var (
		elasticClient *es.Client
		endpoint      = "http://docker:9200"
		err           error
		ctx           = context.Background()
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

	Context("Index", func() {
		var (
			indexName = "learn"
			// typeName  = "learnType"
		// query = "aman"
		// skip  = 1
		// top   = 2
		)
		BeforeEach(func() {
			//Create Index
			_, err = elasticClient.Indices.Create(indexName)
			Expect(err).To(BeNil())
		})

		It("should exist", func() {
			_, err = elasticClient.Indices.Exists([]string{indexName})
			Expect(err).To(BeNil())
		})

		Context("Document", func() {
			var (
				docs = []frameworks.Document{
					{ID: shortid.MustGenerate(), Title: "Aman", Content: "Preet", CreatedAt: time.Now().UTC()},
					{ID: shortid.MustGenerate(), Title: "John", Content: "Doe", CreatedAt: time.Now().UTC()},
				}
			)

			BeforeEach(func() {
				// Insert documents in bulk
				for _, doc := range docs {
					payload, _ := json.Marshal(doc)

					_, err = esapi.CreateRequest{
						Index:      indexName,
						DocumentID: doc.ID,
						Body:       bytes.NewReader(payload),
					}.Do(ctx, elasticClient)
					Expect(err).To(BeNil())
				}
			})

			It("should exist", func() {
				// exists, err := esapi.Exists(indexName, docs[0].ID)
				// Expect(err).To(BeNil())
				// Expect(exists).To(BeTrue())
			})

		})
	})
})
