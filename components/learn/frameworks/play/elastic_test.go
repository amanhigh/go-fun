package play_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models"
	"github.com/amanhigh/go-fun/models/learn/frameworks"
	"github.com/bxcodec/faker/v3"
	es "github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
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
		// Create Test Container
		esContainer, err = util.ElasticSearchTestContainer(ctx)
		Expect(err).To(BeNil())

		// Get Mapped Port
		endpoint, err = esContainer.PortEndpoint(ctx, "9200/tcp", "")
		Expect(err).To(BeNil())
		log.Info().Str("Host", endpoint).Msg("Elastic Endpoint")

		// Elastic Client
		elasticClient, err = es.NewClient(es.Config{Addresses: []string{"http://" + endpoint}})
		Expect(err).To(BeNil())
		Expect(elasticClient).ToNot(BeNil())
	})

	AfterAll(func() {
		log.Warn().Msg("Elastic Shutting Down")
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
			// Create Index
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

	Context("Students", func() {
		var (
			students []Student
		)

		BeforeEach(func() {
			// Generate 100 students
			for i := 0; i < 100; i++ {
				var s Student
				Expect(faker.FakeData(&s)).NotTo(HaveOccurred())
				s.CreatedAt = time.Now().UTC()
				s.UpdatedAt = time.Now().UTC()
				students = append(students, s)
			}
		})

		It("generates and insert", func() {
			for i, s := range students {
				studentJSON, err := json.Marshal(s)
				Expect(err).NotTo(HaveOccurred())

				// Insert student into Elasticsearch
				// http://docker:9200/students
				// Kibana -> Analytics -> Discover -> Data View (Name:Students, Index Pattern: students, Timestamp: Created At.)
				res, err := esapi.IndexRequest{
					Index:      "students",
					DocumentID: fmt.Sprintf("%d", i),
					Body:       strings.NewReader(string(studentJSON)),
					Refresh:    "true",
				}.Do(ctx, elasticClient)

				Expect(err).NotTo(HaveOccurred())
				Expect(res.IsError()).To(BeFalse())
				res.Body.Close()
			}
		})
	})
})
