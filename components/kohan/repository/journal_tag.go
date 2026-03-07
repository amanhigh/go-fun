//nolint:dupl
package repository

// TagRepository provides persistence operations for journal tags.
// Tags are categorical labels that organize and classify journal entries.

import (
	"context"
	"errors"

	"github.com/amanhigh/go-fun/common/util"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
	"gorm.io/gorm"
)

// TagRepository provides persistence operations for journal tags.
// Tags provide categorical organization with type-based filtering capabilities.
type TagRepository interface {
	util.BaseDbRepositoryInterface
	// ListTags returns all tags for an entry, optionally filtered by type.
	// Type can be "reason", "management", or empty for all tags.
	ListTags(ctx context.Context, journalID string, tagType string) ([]barkat.Tag, common.HttpError)
}

type TagRepositoryImpl struct {
	util.BaseDbRepository
}

var _ TagRepository = (*TagRepositoryImpl)(nil)

// NewTagRepository creates a new TagRepository backed by GORM.
func NewTagRepository(db *gorm.DB) *TagRepositoryImpl {
	return &TagRepositoryImpl{BaseDbRepository: util.NewBaseDbRepository(db)}
}

func (r *TagRepositoryImpl) ListTags(ctx context.Context, journalID, tagType string) ([]barkat.Tag, common.HttpError) {
	var tags []barkat.Tag
	var txErr error
	query := r.SafeTx(ctx).Where("journal_id = ?", journalID)
	if tagType != "" {
		query = query.Where("type = ?", tagType)
	}
	if txErr = query.Order("created_at").Find(&tags).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		return nil, util.GormErrorMapper(txErr)
	}
	return tags, nil
}
