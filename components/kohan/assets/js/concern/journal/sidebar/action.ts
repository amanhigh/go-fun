import { createSubmitter, type Submitter } from '../../../lib/submitter';
import type { DisplaySpec } from '../../../types/core/present';
import type { QuickAction } from '../../../types/journal/sidebar';
import { JournalType, JournalStatus } from '../../../types/api/journal/enums';
import type { Journal } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

function localToday(pg: JournalDetailPageProvider): string {
	return pg().present.date.humanDate(new Date());
}

// ===== Type Helpers =====

function isTaken(journal: Journal): boolean {
	return journal.type === JournalType.TAKEN;
}

function isRejected(journal: Journal): boolean {
	return journal.type === JournalType.REJECTED;
}

// ===== Display Helpers =====

function reviewDisplay(journal: Journal): DisplaySpec {
	const hasReview = !!journal.reviewed_at;
	return {
		text: hasReview ? 'Unreviewed' : 'Reviewed',
		class: hasReview ? 'journal-review-toggle-pending' : 'journal-review-toggle-reviewed',
	};
}

function runningDisplay(): DisplaySpec {
	return { text: 'Running', class: 'journal-quick-status-running' };
}

function successDisplay(): DisplaySpec {
	return { text: 'Success', class: 'journal-quick-status-success' };
}

function failDisplay(): DisplaySpec {
	return { text: 'Fail', class: 'journal-quick-status-fail' };
}

function justLossDisplay(): DisplaySpec {
	return { text: 'Just Loss', class: 'journal-quick-status-loss' };
}

function brokenDisplay(): DisplaySpec {
	return { text: 'Broken', class: 'journal-quick-status-broken' };
}

function setDisplay(): DisplaySpec {
	return { text: 'Set', class: 'journal-quick-status-set' };
}

// ===== Async Action Handlers =====

/** Reviewed/Unreviewed toggle — changes ONLY reviewed_at, never status. */
async function toggleReviewedAt(submitter: Submitter, pg: JournalDetailPageProvider): Promise<void> {
	const journal = pg().journal.detail!;
	const reviewedAt = journal.reviewed_at ? null : localToday(pg);
	const successMsg = reviewedAt ? 'Journal marked reviewed.' : 'Journal marked not reviewed.';
	await submitter.run(async () => {
		const envelope = await pg().client.updateReview(pg().journal.detail!.id, { reviewed_at: reviewedAt });
		journal.reviewed_at = envelope.data.reviewed_at;
		// Intentionally NOT updating journal.status — review toggle only touches reviewed_at.
	}, { success: successMsg });
}

/** Status-only update — changes ONLY status, never reviewed_at. */
async function applyStatusOnly(submitter: Submitter, pg: JournalDetailPageProvider, targetStatus: JournalStatus, successMsg: string): Promise<void> {
	await submitter.run(async () => {
		const envelope = await pg().client.updateReview(pg().journal.detail!.id, { status: targetStatus });
		pg().journal.detail!.status = envelope.data.status;
	}, { success: successMsg });
}

// ===== Exported Concern =====

export function NewReviewBarConcern(pg: JournalDetailPageProvider) {
	return {
		submitter: createSubmitter(),

		actions(): QuickAction[] {
			const journal = pg().journal.detail;
			if (!journal) return [];

			return [
				// Every journal gets the review toggle (reviewed_at only).
				{ id: 'review-toggle', isActive: () => true, display: reviewDisplay(journal), apply: () => toggleReviewedAt(this.submitter, pg) },

				// SET → RUNNING (both TAKEN and REJECTED).
				{ id: 'status-running', isActive: () => journal.status === JournalStatus.SET, display: runningDisplay(), apply: () => applyStatusOnly(this.submitter, pg, JournalStatus.RUNNING, 'Marked as Running') },

				// TAKEN: RUNNING → SUCCESS
				// REJECTED: FAIL → SUCCESS
				{ id: 'status-success', isActive: () =>
					(isTaken(journal) && journal.status === JournalStatus.RUNNING) ||
					(isRejected(journal) && journal.status === JournalStatus.FAIL),
					display: successDisplay(), apply: () => applyStatusOnly(this.submitter, pg, JournalStatus.SUCCESS, 'Marked as Success') },

				// TAKEN: RUNNING → FAIL
				// REJECTED: RUNNING → FAIL, SUCCESS → FAIL, BROKEN → FAIL
				{ id: 'status-fail', isActive: () =>
					(isTaken(journal) && journal.status === JournalStatus.RUNNING) ||
					(isRejected(journal) && (journal.status === JournalStatus.RUNNING || journal.status === JournalStatus.SUCCESS || journal.status === JournalStatus.BROKEN)),
					display: failDisplay(), apply: () => applyStatusOnly(this.submitter, pg, JournalStatus.FAIL, 'Marked as Fail') },

				// TAKEN: FAIL / MISSED → JUST_LOSS (SUCCESS goes to SET instead)
				{ id: 'status-just-loss', isActive: () =>
					isTaken(journal) && ([JournalStatus.FAIL, JournalStatus.MISSED] as JournalStatus[]).includes(journal.status),
					display: justLossDisplay(), apply: () => applyStatusOnly(this.submitter, pg, JournalStatus.JUST_LOSS, 'Marked as Just Loss') },

				// TAKEN: SUCCESS → SET, JUST_LOSS → SET (full cycle reset)
				{ id: 'status-set', isActive: () => isTaken(journal) && (journal.status === JournalStatus.SUCCESS || journal.status === JournalStatus.JUST_LOSS), display: setDisplay(), apply: () => applyStatusOnly(this.submitter, pg, JournalStatus.SET, 'Reset to Set') },

				// REJECTED: FAIL → BROKEN
				{ id: 'status-broken', isActive: () => isRejected(journal) && journal.status === JournalStatus.FAIL, display: brokenDisplay(), apply: () => applyStatusOnly(this.submitter, pg, JournalStatus.BROKEN, 'Marked as Broken') },
			].filter((action) => action.isActive());
		},
	};
}
