package manager

import (
	"context"
	"fmt"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// JournalManager provides business logic for journal operations.
type JournalManager interface {
	// CreateJournal creates a new journal with associations.
	CreateJournal(ctx context.Context, journal *barkat.Journal) common.HttpError
	// GetJournal retrieves a single journal by ID with all associations.
	GetJournal(ctx context.Context, id string) (barkat.Journal, common.HttpError)
	// ListJournals returns a filtered, paginated list of journal summaries.
	ListJournals(ctx context.Context, query barkat.JournalQuery) (barkat.JournalList, common.HttpError)
	// DeleteJournal deletes a journal by ID.
	DeleteJournal(ctx context.Context, id string) common.HttpError
}

type JournalManagerImpl struct {
	repo repository.JournalRepository
}

var _ JournalManager = (*JournalManagerImpl)(nil)

// NewJournalManager creates a new JournalManager.
func NewJournalManager(repo repository.JournalRepository) *JournalManagerImpl {
	return &JournalManagerImpl{repo: repo}
}

// ---- Journal ----

func (m *JournalManagerImpl) CreateJournal(ctx context.Context, journal *barkat.Journal) common.HttpError {
	// Business rule: validate unique timeframes (PRD Section 3.1)
	if httpErr := m.validateUniqueTimeframes(journal.Images); httpErr != nil {
		return httpErr
	}

	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		return m.repo.Create(c, journal)
	})
}

// validateUniqueTimeframes ensures all image timeframes are unique within a journal.
// PRD Section 3.1: "Business rule: minimum 4 unique timeframes"
func (m *JournalManagerImpl) validateUniqueTimeframes(images []barkat.Image) common.HttpError {
	seen := make(map[string]bool)
	for _, img := range images {
		if seen[img.Timeframe] {
			return common.NewFieldHttpError("Images", "Duplicate timeframe not allowed: "+img.Timeframe)
		}
		seen[img.Timeframe] = true
	}
	return nil
}

func (m *JournalManagerImpl) GetJournal(ctx context.Context, id string) (barkat.Journal, common.HttpError) {
	journal, err := m.repo.GetJournal(ctx, id)
	if err != nil {
		return barkat.Journal{}, common.ErrNotFound
	}
	return journal, nil
}

func (m *JournalManagerImpl) ListJournals(ctx context.Context, query barkat.JournalQuery) (barkat.JournalList, common.HttpError) {
	journals, total, err := m.repo.ListJournals(ctx, query)
	if err != nil {
		return barkat.JournalList{}, common.NewServerError(fmt.Errorf("failed to list journals: %w", err))
	}
	return barkat.JournalList{
		Records: journals,
		Metadata: common.PaginatedResponse{
			Total:  total,
			Offset: query.Offset,
			Limit:  query.Limit,
		},
	}, nil
}

func (m *JournalManagerImpl) DeleteJournal(ctx context.Context, id string) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// First fetch the journal by external_id to get internal ID
		journal, httpErr := m.GetJournal(c, id)
		if httpErr != nil {
			return httpErr
		}

		// Now delete by internal ID
		return m.repo.DeleteById(c, journal.ID, &barkat.Journal{})
	})
}
