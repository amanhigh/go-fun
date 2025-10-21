package publisher

import (
	"context"
	"fmt"

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

// PublishWithExtras marshals the payload, attaches context-derived metadata, and emits it on the given topic.
// It expects the correlation id to be present in the context (via common.WithCorrelation).
func (bp BasePublisher) PublishWithExtras(ctx context.Context, topic string, payload any, extras map[string]string) common.HttpError {
	metadata, err := bp.buildMetadata(ctx, extras)
	if err != nil {
		return common.NewServerError(err)
	}

	if err = util.PublishJSONMessage(ctx, bp.publisher, topic, payload, metadata); err != nil {
		return common.NewServerError(err)
	}
	return nil
}

func (bp BasePublisher) buildMetadata(ctx context.Context, extras map[string]string) (common.Metadata, error) {
	correlation := common.CorrelationFrom(ctx)
	if correlation == "" {
		return nil, fmt.Errorf("correlation id missing from context")
	}

	meta := common.MustBaseMetadata(correlation)

	if causation := common.CausationFrom(ctx); causation != "" {
		meta = meta.WithCausation(causation)
	}

	if len(extras) > 0 {
		meta = meta.WithPairs(extras)
	}

	return meta, nil
}
