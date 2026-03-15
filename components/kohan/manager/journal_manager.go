package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/amanhigh/go-fun/components/kohan/repository"
	"github.com/amanhigh/go-fun/models/barkat"
	"github.com/amanhigh/go-fun/models/common"
)

// JournalManager provides business logic for journal operations.
type JournalManager interface {
	// CreateJournal creates a new journal with associations.
	CreateJournal(ctx context.Context, journal *barkat.Journal) common.HttpError
	// GetJournal retrieves a single journal by EXTERNAL_ID with all associations.
	GetJournal(ctx context.Context, journalExternalId string) (barkat.Journal, common.HttpError)
	// ListJournals returns a filtered, paginated list of journal summaries.
	ListJournals(ctx context.Context, query barkat.JournalQuery) (barkat.JournalList, common.HttpError)
	// DeleteJournal deletes a journal by EXTERNAL_ID.
	DeleteJournal(ctx context.Context, journalExternalId string) common.HttpError
	// UpdateReviewStatus updates the review status of a journal by EXTERNAL_ID.
	UpdateReviewStatus(ctx context.Context, journalExternalId string, update barkat.JournalReviewUpdate) (barkat.Journal, common.HttpError)
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
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		return m.repo.Create(c, journal)
	})
}

func (m *JournalManagerImpl) GetJournal(ctx context.Context, journalExternalId string) (barkat.Journal, common.HttpError) {
	journal, err := m.repo.GetJournal(ctx, journalExternalId)
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
		Journals: journals,
		Metadata: common.PaginatedResponse{
			Total:  total,
			Offset: query.Offset,
			Limit:  query.Limit,
		},
	}, nil
}

func (m *JournalManagerImpl) DeleteJournal(ctx context.Context, journalExternalId string) common.HttpError {
	return m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// First fetch the journal by external_id to get internal ID
		journal, httpErr := m.GetJournal(c, journalExternalId)
		if httpErr != nil {
			return httpErr
		}

		// Now delete by internal ID
		return m.repo.DeleteById(c, journal.ID, &barkat.Journal{})
	})
}

func (m *JournalManagerImpl) UpdateReviewStatus(ctx context.Context, journalExternalId string, update barkat.JournalReviewUpdate) (barkat.Journal, common.HttpError) {
	var updatedJournal barkat.Journal
	err := m.repo.UseOrCreateTx(ctx, func(c context.Context) common.HttpError {
		// Get journal to update
		journal, httpErr := m.GetJournal(c, journalExternalId)
		if httpErr != nil {
			return httpErr
		}

		// Update reviewed_at based on the update request
		if update.ReviewedAt != nil {
			// Convert civil.Date to time.Time for GORM using the In() function
			reviewedTime := update.ReviewedAt.In(time.UTC)
			journal.ReviewedAt = &reviewedTime
		} else {
			journal.ReviewedAt = nil
		}

		// Save the updated journal
		if httpErr := m.repo.Update(c, &journal); httpErr != nil {
			return httpErr
		}

		// Set the updated journal to return
		updatedJournal = journal
		return nil
	})

	if err != nil {
		return barkat.Journal{}, err
	}

	return updatedJournal, nil
}
