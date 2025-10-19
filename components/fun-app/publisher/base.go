package publisher

import (
	"context"

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

// Publish marshals the payload, attaches metadata, and emits it on the given topic.
func (bp BasePublisher) Publish(ctx context.Context, topic string, payload any, metadata map[string]string) common.HttpError {
	if err := util.PublishJSONMessage(ctx, bp.publisher, topic, payload, metadata); err != nil {
		return common.NewServerError(err)
	}
	return nil
}
