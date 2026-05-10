import { createSubmitter, type Submitter } from '../../../lib/submitter';
import type { DisplaySpec } from '../../../types/present';
import type { QuickAction } from '../../../types/quick_action';
import { JournalType, JournalStatus } from '../../../types/api/journal/enums';
import type { Journal } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

function localToday(pg: JournalDetailPageProvider): string {
	return pg().present.date.humanDate(new Date());
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
		case JournalType.TAKEN: return { text: 'Mark Just Loss', class: 'journal-quick-status-loss' };
		default: return { text: 'Mark Broken', class: 'journal-quick-status-broken' };
	}
}

// ===== Active Checks =====

function isStatusActive(journal: Journal): boolean {
	switch (journal.type) {
		case JournalType.TAKEN: return journal.status !== JournalStatus.JUST_LOSS;
		default: return journal.status !== JournalStatus.BROKEN;
	}
}

// ===== Async Action Handlers =====

async function toggleReviewedAt(submitter: Submitter, pg: JournalDetailPageProvider): Promise<void> {
	const journal = pg().current.journal!;
	const reviewedAt = journal.reviewed_at ? null : localToday(pg);
	const successMsg = reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.';
	await submitter.run(async () => {
		const envelope = await pg().client.updateReview(pg().current.journalId, { reviewed_at: reviewedAt });
		journal.reviewed_at = envelope.data.reviewed_at;
		journal.status = envelope.data.status;
	}, { success: successMsg, error: 'Unable to update review date.' });
}

async function applyReviewStatus(submitter: Submitter, pg: JournalDetailPageProvider): Promise<void> {
	const journal = pg().current.journal!;
	const isTaken = journal.type === JournalType.TAKEN;
	const targetStatus = isTaken ? JournalStatus.JUST_LOSS : JournalStatus.BROKEN;

	await submitter.run(async () => {
		const envelope = await pg().client.updateReview(pg().current.journalId, { status: targetStatus, reviewed_at: localToday(pg) });
		journal.reviewed_at = envelope.data.reviewed_at;
		journal.status = envelope.data.status;
		await pg().sidebar.reviewQueue.load();
	}, { success: `${isTaken ? 'Mark Just Loss' : 'Mark Broken'} applied and journal marked reviewed.`, error: 'Unable to update journal status.' });
}

// ===== Exported Concern =====

export function NewReviewActionsConcern(pg: JournalDetailPageProvider) {
	return {
		submitter: createSubmitter(),

		actions(): QuickAction[] {
			const journal = pg().current.journal;
			if (!journal) return [];

			return [
				{ id: 'review-toggle', isActive: () => true, display: reviewDisplay(journal), apply: () => toggleReviewedAt(this.submitter, pg) },
				{ id: 'review-status', isActive: () => isStatusActive(journal), display: statusDisplay(journal), apply: () => applyReviewStatus(this.submitter, pg) },
			].filter((action) => action.isActive());
		},
	};
}
