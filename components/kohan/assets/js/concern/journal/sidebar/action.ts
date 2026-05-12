import { createSubmitter, type Submitter } from '../../../lib/submitter';
import type { DisplaySpec } from '../../../types/core/present';
import type { QuickAction } from '../../../types/journal/sidebar';
import { JournalType, JournalStatus } from '../../../types/api/journal/enums';
import type { Journal } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

function localToday(pg: JournalDetailPageProvider): string {
	return pg().present.date.humanDate(new Date());
}

// ===== Transition Table =====

const statusTransitions: Record<string, Record<string, JournalStatus[]>> = {
	[JournalType.TAKEN]: {
		[JournalStatus.SET]: [JournalStatus.RUNNING],
		[JournalStatus.RUNNING]: [JournalStatus.SUCCESS, JournalStatus.FAIL],
		[JournalStatus.SUCCESS]: [JournalStatus.SET],
		[JournalStatus.FAIL]: [JournalStatus.JUST_LOSS],
		[JournalStatus.MISSED]: [JournalStatus.JUST_LOSS],
		[JournalStatus.JUST_LOSS]: [JournalStatus.SET],
	},
	[JournalType.REJECTED]: {
		[JournalStatus.SET]: [JournalStatus.RUNNING],
		[JournalStatus.RUNNING]: [JournalStatus.FAIL],
		[JournalStatus.FAIL]: [JournalStatus.SUCCESS, JournalStatus.BROKEN],
		[JournalStatus.SUCCESS]: [JournalStatus.FAIL],
		[JournalStatus.BROKEN]: [JournalStatus.FAIL],
	},
};

function nextStatuses(journal: Journal): JournalStatus[] {
	const typeTransitions = statusTransitions[journal.type];
	if (!typeTransitions) return [];
	return typeTransitions[journal.status] ?? [];
}

// ===== Display Helpers =====

function reviewDisplay(journal: Journal): DisplaySpec {
	const hasReview = !!journal.reviewed_at;
	return {
		text: hasReview ? 'Unreviewed' : 'Reviewed',
		class: hasReview ? 'journal-review-toggle-pending' : 'journal-review-toggle-reviewed',
	};
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

// ===== Action Builders =====

function statusAction(pg: JournalDetailPageProvider, submitter: Submitter, targetStatus: JournalStatus): QuickAction {
	const spec = pg().present.status.spec(targetStatus);
	return {
		id: `status-${targetStatus.toLowerCase()}`,
		isActive: () => true,
		display: { text: spec.text, class: spec.class },
		apply: () => applyStatusOnly(submitter, pg, targetStatus, `Marked as ${spec.text}`),
	};
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
				...nextStatuses(journal).map((target) => statusAction(pg, this.submitter, target)),
			].filter((action) => action.isActive());
		},
	};
}
