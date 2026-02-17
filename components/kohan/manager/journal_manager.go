package manager

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// JournalManager provides business logic for journal entry operations.
type JournalManager interface {
	// CreateEntry creates a new journal entry with associations.
	CreateEntry(ctx context.Context, entry *barkat.Entry) common.HttpError
	// GetEntry retrieves a single entry by ID with all associations.
	GetEntry(ctx context.Context, id string) (barkat.Entry, common.HttpError)
	// ListEntries returns a filtered, paginated list of entry summaries.
	ListEntries(ctx context.Context, query barkat.EntryQuery) (barkat.EntryList, common.HttpError)
	// EntryExists checks if an entry with the given ID exists.
	EntryExists(ctx context.Context, entryID string) common.HttpError
	// DeleteEntry deletes a journal entry by ID.
	DeleteEntry(ctx context.Context, id string) common.HttpError
}

type JournalManagerImpl struct {
	repo repository.JournalRepository
}

var _ JournalManager = (*JournalManagerImpl)(nil)

// NewJournalManager creates a new JournalManager.
func NewJournalManager(repo repository.JournalRepository) *JournalManagerImpl {
	return &JournalManagerImpl{repo: repo}
}

// ---- Entry ----

func (m *JournalManagerImpl) CreateEntry(ctx context.Context, entry *barkat.Entry) common.HttpError {
	// BUG: Transactional Support is missing.
	if err := m.repo.CreateEntry(ctx, entry); err != nil {
		return common.NewServerError(fmt.Errorf("failed to create entry: %w", err))
	}
	return nil
}

func (m *JournalManagerImpl) GetEntry(ctx context.Context, id string) (barkat.Entry, common.HttpError) {
	entry, err := m.repo.GetEntry(ctx, id)
	if err != nil {
		return barkat.Entry{}, common.ErrNotFound
	}
	return entry, nil
}

func (m *JournalManagerImpl) ListEntries(ctx context.Context, query barkat.EntryQuery) (barkat.EntryList, common.HttpError) {
	entries, total, err := m.repo.ListEntries(ctx, query)
	if err != nil {
		return barkat.EntryList{}, common.NewServerError(fmt.Errorf("failed to list entries: %w", err))
	}
	return barkat.EntryList{
		Records:  entries,
		Metadata: common.PaginatedResponse{Total: total},
	}, nil
}

func (m *JournalManagerImpl) EntryExists(ctx context.Context, entryID string) common.HttpError {
	entry := &barkat.Entry{}
	return m.repo.FindById(ctx, entryID, entry)
}

func (m *JournalManagerImpl) DeleteEntry(ctx context.Context, id string) common.HttpError {
	entry := &barkat.Entry{}
	return m.repo.DeleteById(ctx, id, entry)
}
