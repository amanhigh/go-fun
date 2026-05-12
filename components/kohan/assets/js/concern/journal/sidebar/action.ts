import { createSubmitter, type Submitter } from '../../../lib/submitter';
import type { DisplaySpec } from '../../../types/core/present';
import type { QuickAction } from '../../../types/journal/sidebar';
import { JournalType, JournalStatus } from '../../../types/api/journal/enums';
import type { Journal } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

function localToday(pg: JournalDetailPageProvider): string {
	return pg().present.date.humanDate(new Date());
}

// ===== Terminal Status Check =====

const TERMINAL_STATUSES: ReadonlySet<JournalStatus> = new Set([
	JournalStatus.SUCCESS,
	JournalStatus.FAIL,
	JournalStatus.MISSED,
	JournalStatus.JUST_LOSS,
	JournalStatus.BROKEN,
]);

function isTerminalStatus(journal: Journal): boolean {
	return TERMINAL_STATUSES.has(journal.status);
}

// ===== Display Helpers =====

function reviewDisplay(journal: Journal): DisplaySpec {
	const hasReview = !!journal.reviewed_at;
	return {
		text: hasReview ? 'Unreviewed' : 'Reviewed',
		class: hasReview ? 'journal-review-toggle-pending' : 'journal-review-toggle-reviewed',
	};
}

function statusDisplay(journal: Journal): DisplaySpec {
	switch (journal.type) {
		case JournalType.TAKEN: return { text: 'Just Loss', class: 'journal-quick-status-loss' };
		default: return { text: 'Broken', class: 'journal-quick-status-broken' };
	}
}

// ===== Active Checks =====

function isStatusActive(journal: Journal): boolean {
	switch (journal.type) {
		case JournalType.TAKEN: return journal.status !== JournalStatus.JUST_LOSS;
		default: return journal.status !== JournalStatus.BROKEN;
	}
}

// FIXME: Add new QuickActions to transition journals to SUCCESS or FAIL status.
// Currently only JUST_LOSS (TAKEN) / BROKEN (default) are supported via review-status.
// Runner-up trades from RUNNING have no buttons for explicit SUCCESS/FAIL transitions.
// ===== Async Action Handlers =====

async function toggleReviewedAt(submitter: Submitter, pg: JournalDetailPageProvider): Promise<void> {
	const journal = pg().journal.detail!;
	const reviewedAt = journal.reviewed_at ? null : localToday(pg);
	const successMsg = reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.';
	await submitter.run(async () => {
		const envelope = await pg().client.updateReview(pg().journal.detail!.id, { reviewed_at: reviewedAt });
		journal.reviewed_at = envelope.data.reviewed_at;
		journal.status = envelope.data.status;
	}, { success: successMsg });
}

async function applyReviewStatus(submitter: Submitter, pg: JournalDetailPageProvider): Promise<void> {
	const journal = pg().journal.detail!;
	const isTaken = journal.type === JournalType.TAKEN;
	const targetStatus = isTaken ? JournalStatus.JUST_LOSS : JournalStatus.BROKEN;

	await submitter.run(async () => {
		const envelope = await pg().client.updateReview(pg().journal.detail!.id, { status: targetStatus, reviewed_at: localToday(pg) });
		journal.reviewed_at = envelope.data.reviewed_at;
		journal.status = envelope.data.status;
		await pg().sidebar.reviewQueue.load();
	}, { success: `${isTaken ? 'Mark Just Loss' : 'Mark Broken'} applied and journal marked reviewed.` });
}

// ===== Exported Concern =====

export function NewReviewBarConcern(pg: JournalDetailPageProvider) {
	return {
		submitter: createSubmitter(),

		actions(): QuickAction[] {
			const journal = pg().journal.detail;
			if (!journal) return [];

			return [
				{ id: 'review-toggle', isActive: () => true, display: reviewDisplay(journal), apply: () => toggleReviewedAt(this.submitter, pg) },
				{ id: 'review-status', isActive: () => isStatusActive(journal), display: statusDisplay(journal), apply: () => applyReviewStatus(this.submitter, pg) },
			].filter((action) => action.isActive());
		},
	};
}
