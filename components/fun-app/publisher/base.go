package publisher

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/common"
)

// BasePublisher offers shared publishing helpers for saga event publishers.
type BasePublisher struct {
	publisher message.Publisher
}

// NewBasePublisher constructs a BasePublisher around a watermill publisher.
func NewBasePublisher(p message.Publisher) BasePublisher {
	return BasePublisher{publisher: p}
}

// PublishRoot publishes a message that starts a new correlation.
func (bp BasePublisher) PublishRoot(ctx context.Context, topic string, payload any) common.HttpError {
	metadata := common.NewRootMetadata(watermill.NewUUID(), watermill.NewUUID())
	if err := util.PublishJSONMessage(ctx, bp.publisher, topic, payload, metadata); err != nil {
		return common.NewServerError(err)
	}
	return nil
}

// PublishChild publishes a message linked to the typed parent metadata in ctx.
func (bp BasePublisher) PublishChild(ctx context.Context, topic string, payload any) common.HttpError {
	parent := common.MetadataFromContext(ctx)
	if parent.CorrelationID == "" {
		parent.CorrelationID = watermill.NewUUID()
	}
	metadata := common.NewChildMetadata(watermill.NewUUID(), parent)
	if err := util.PublishJSONMessage(ctx, bp.publisher, topic, payload, metadata); err != nil {
		return common.NewServerError(err)
	}
	return nil
}
