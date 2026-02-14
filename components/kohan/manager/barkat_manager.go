package manager

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// BarkatManager provides business logic for barkat journal operations.
//
//go:generate mockery --name BarkatManager
// HACK: Match name with entity it should be JournalManager (RenameEntity to Journal)
type BarkatManager interface {
	// CreateEntry creates a new barkat entry with optional images.
	CreateEntry(ctx context.Context, entry *barkat.Entry) common.HttpError
	// GetEntry retrieves a single entry by ID.
	GetEntry(ctx context.Context, id string) (barkat.Entry, common.HttpError)
	// ListEntries returns a filtered, paginated list of entries.
	ListEntries(ctx context.Context, query barkat.EntryQuery) (barkat.EntryList, common.HttpError)
}

type BarkatManagerImpl struct {
	repo repository.BarkatRepository
}

var _ BarkatManager = (*BarkatManagerImpl)(nil)

// NewBarkatManager creates a new BarkatManager.
func NewBarkatManager(repo repository.BarkatRepository) *BarkatManagerImpl {
	return &BarkatManagerImpl{repo: repo}
}

func (m *BarkatManagerImpl) CreateEntry(ctx context.Context, entry *barkat.Entry) common.HttpError {
	if err := m.repo.CreateEntry(ctx, entry); err != nil {
		return common.NewServerError(fmt.Errorf("failed to create entry: %w", err))
	}
	return nil
}

func (m *BarkatManagerImpl) GetEntry(ctx context.Context, id string) (barkat.Entry, common.HttpError) {
	entry, err := m.repo.GetEntry(ctx, id)
	if err != nil {
		return barkat.Entry{}, common.ErrNotFound
	}
	return entry, nil
}

func (m *BarkatManagerImpl) ListEntries(ctx context.Context, query barkat.EntryQuery) (barkat.EntryList, common.HttpError) {
	entries, total, err := m.repo.ListEntries(ctx, query)
	if err != nil {
		return barkat.EntryList{}, common.NewServerError(fmt.Errorf("failed to list entries: %w", err))
	}
	return barkat.EntryList{
		Records:  entries,
		Metadata: common.PaginatedResponse{Total: total},
	}, nil
}
