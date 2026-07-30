package common_test

import (
	"context"

	"github.com/amanhigh/go-fun/models/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Metadata", func() {
	Context("root metadata", func() {
		var metadata common.Metadata

		BeforeEach(func() {
			metadata = common.NewRootMetadata("message-1", "correlation-1")
		})

		It("contains the message and correlation identifiers", func() {
			Expect(metadata).To(Equal(common.Metadata{
				MessageID:     "message-1",
				CorrelationID: "correlation-1",
			}))
			Expect(metadata.CausationID).To(BeEmpty())
		})
	})

	Context("child metadata", func() {
		var metadata common.Metadata

		BeforeEach(func() {
			parent := common.NewRootMetadata("parent-message", "correlation-1")
			metadata = common.NewChildMetadata("child-message", parent)
		})

		It("inherits correlation and records the parent message as causation", func() {
			Expect(metadata).To(Equal(common.Metadata{
				MessageID:     "child-message",
				CorrelationID: "correlation-1",
				CausationID:   "parent-message",
			}))
		})
	})

	Context("context storage", func() {
		var metadata common.Metadata
		var result common.Metadata

		BeforeEach(func() {
			metadata = common.Metadata{
				MessageID:     "message-1",
				CorrelationID: "correlation-1",
				CausationID:   "parent-message",
			}
			ctx := common.WithMetadata(context.Background(), metadata)
			result = common.MetadataFromContext(ctx)
		})

		It("round-trips the typed value object", func() {
			Expect(result).To(Equal(metadata))
		})
	})

	Context("absent context", func() {
		It("returns zero metadata for a nil context", func() {
			Expect(common.MetadataFromContext(context.TODO())).To(Equal(common.Metadata{}))
		})

		It("accepts nil context when storing metadata", func() {
			metadata := common.Metadata{MessageID: "message-1"}
			ctx := common.WithMetadata(context.TODO(), metadata)

			Expect(common.MetadataFromContext(ctx)).To(Equal(metadata))
		})
	})
})
