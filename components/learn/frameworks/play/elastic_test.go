package play_test

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/learn/frameworks"
	es "github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/fatih/color"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/teris-io/shortid"
	"github.com/testcontainers/testcontainers-go"
)

/*
 Guide - https://www.elastic.co/guide/en/elasticsearch/client/go-api/current/getting-started-go.html
 Store Example - https://github.com/elastic/go-elasticsearch/blob/main/_examples/xkcdsearch/store.go
 Index - http://docker:9200/learn/_search?pretty
*/

var _ = Describe("Elastic", Ordered, Label(models.GINKGO_SLOW), func() {
	var (
		elasticClient *es.Client
		endpoint      = "docker:9200"
		err           error
		ctx           = context.Background()
		esContainer   testcontainers.Container
	)

	BeforeAll(func() {
		//Create Test Container
		esContainer, err = util.ElasticSearchTestContainer(ctx)
		Expect(err).To(BeNil())

		//Get Mapped Port
		endpoint, err = esContainer.Endpoint(ctx, "")
		Expect(err).To(BeNil())
		color.Green("Elastic Endpoint: %s", endpoint)

		//Elastic Client
		elasticClient, err = es.NewClient(es.Config{Addresses: []string{"http://" + endpoint}})
		Expect(err).To(BeNil())
		Expect(elasticClient).ToNot(BeNil())
	})

	AfterAll(func() {
		color.Red("Elastic Shutting Down")
		err = esContainer.Terminate(ctx)
		Expect(err).To(BeNil())
	})

	It("should connect", func() {
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

			AfterEach(func() {
				// Delete documents
				for _, doc := range docs {
					_, err = esapi.DeleteRequest{
						Index:      indexName,
						DocumentID: doc.ID,
					}.Do(ctx, elasticClient)
					Expect(err).To(BeNil())
				}
			})

			It("should get", func() {
				// Get document
				resp, err := esapi.GetRequest{
					Index:      indexName,
					DocumentID: docs[0].ID,
				}.Do(ctx, elasticClient)
				Expect(err).To(BeNil())
				Expect(resp).ToNot(BeNil())
			})

			It("should search", func() {
				// Search documents
				resp, err := esapi.SearchRequest{
					Index: []string{indexName},
					Body: strings.NewReader(`{
						"query": {
							"match": {
								"title": "aman"
							}
						}
					}`),
				}.Do(ctx, elasticClient)
				Expect(err).To(BeNil())
				Expect(resp).ToNot(BeNil())
			})

		})
	})
})
