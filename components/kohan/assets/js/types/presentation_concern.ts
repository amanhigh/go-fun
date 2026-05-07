import type { JournalTag } from './journal_api';

export type DisplaySpec = {
	icon: string;
	text: string;
	class: string;
};

export type PresentationConcern = {
	type(value: string): DisplaySpec;
	status(value: string): DisplaySpec;
	timeframe(value: string): DisplaySpec;

	sequence(value: string | null | undefined): DisplaySpec;

	reasonTag(tag: JournalTag): DisplaySpec;
	directionalTag(tag: JournalTag): DisplaySpec;

	reviewedAt(value: string | null | undefined): DisplaySpec;
	pendingReview(): DisplaySpec;

	formatTimestamp(value: string | null | undefined): string;
	formatReviewQueueDate(value: string | null | undefined): string;
};
