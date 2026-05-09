import { createFeedback, type Feedback } from '../../../lib/feedback';
import { formatDateInputValue } from '../../../lib/date';
import type { DisplaySpec } from '../../../types/present';
import type { QuickAction } from '../../../types/quick_action';
import type { Journal } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function localToday(): string {
	return formatDateInputValue(new Date());
}

// ===== Display Helpers =====

function reviewDisplay(journal: Journal): DisplaySpec {
	const hasReview = !!journal.reviewed_at;
	return {
		text: hasReview ? 'Mark Pending' : 'Mark Reviewed',
		class: hasReview ? 'journal-review-toggle-pending' : 'journal-review-toggle-reviewed',
	};
}

function statusDisplay(journal: Journal): DisplaySpec {
	switch (journal.type) {
		case 'TAKEN': return { text: 'Mark Just Loss', class: 'journal-quick-status-loss' };
		default: return { text: 'Mark Broken', class: 'journal-quick-status-broken' };
	}
}

// ===== Active Checks =====

function isStatusActive(journal: Journal): boolean {
	switch (journal.type) {
		case 'TAKEN': return journal.status !== 'JUST_LOSS';
		default: return journal.status !== 'BROKEN';
	}
}

// ===== Async Action Handlers =====

async function toggleReviewedAt(feedback: Feedback, pg: JournalDetailPageProvider): Promise<void> {
	const journal = pg().current.journal!;
	const reviewedAt = journal.reviewed_at ? null : localToday();
	const successMsg = reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.';
	await feedback.run(async () => {
		const envelope = await pg().client.updateReview(pg().current.journalId, { reviewed_at: reviewedAt });
		journal.reviewed_at = envelope.data.reviewed_at;
		journal.status = envelope.data.status;
	}, successMsg, 'Unable to update review date.');
}

async function applyReviewStatus(feedback: Feedback, pg: JournalDetailPageProvider): Promise<void> {
	const journal = pg().current.journal!;
	const isTaken = journal.type === 'TAKEN';
	const targetStatus = isTaken ? 'JUST_LOSS' : 'BROKEN';

	await feedback.run(async () => {
		const envelope = await pg().client.updateReview(pg().current.journalId, { status: targetStatus, reviewed_at: localToday() });
		journal.reviewed_at = envelope.data.reviewed_at;
		journal.status = envelope.data.status;
		await pg().sidebar.reviewQueue.load();
	}, `${isTaken ? 'Mark Just Loss' : 'Mark Broken'} applied and journal marked reviewed.`, 'Unable to update journal status.');
}

// ===== Exported Concern =====

export function NewReviewActionsConcern(pg: JournalDetailPageProvider) {
	return {
		...createFeedback(),

		actions(): QuickAction[] {
			const journal = pg().current.journal;
			if (!journal) return [];

			return [
				{ id: 'review-toggle', isActive: () => true, display: reviewDisplay(journal), apply: () => toggleReviewedAt(this, pg) },
				{ id: 'review-status', isActive: () => isStatusActive(journal), display: statusDisplay(journal), apply: () => applyReviewStatus(this, pg) },
			].filter((action) => action.isActive());
		},
	};
}
