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
	util.BaseDbRepository
	// ListTags returns all tags for a journal, optionally filtered by type.
	// Type can be "reason", "management", or empty for all tags.
	ListTags(ctx context.Context, journalID uint64, tagType string) ([]barkat.Tag, common.HttpError)
}

type TagRepositoryImpl struct {
	util.BaseDbRepository
}

var _ TagRepository = (*TagRepositoryImpl)(nil)

// NewTagRepository creates a new TagRepository backed by GORM.
func NewTagRepository(baseRepository util.BaseDbRepository) *TagRepositoryImpl {
	return &TagRepositoryImpl{BaseDbRepository: baseRepository}
}

func (r *TagRepositoryImpl) ListTags(ctx context.Context, journalID uint64, tagType string) ([]barkat.Tag, common.HttpError) {
	var tags []barkat.Tag
	var txErr error
	where := barkat.Tag{JournalID: journalID}
	if tagType != "" {
		where.Type = tagType
	}
	query := r.SafeTx(ctx).Where(&where)
	query = util.ApplySort(query, util.SortOptions{
		DefaultSortBy:    "created_at",
		DefaultSortOrder: common.SortOrderAsc,
	})
	if txErr = query.Find(&tags).Error; txErr != nil && !errors.Is(txErr, gorm.ErrRecordNotFound) {
		return nil, util.GormErrorMapper(txErr)
	}
	return tags, nil
}
