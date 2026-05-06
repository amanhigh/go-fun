export type PresentationConcern = {
	normalizeStatus(value: string): string;
	statusBadgeClass(value: string): string;
	typeBadgeClass(value: string): string;
	formatTimestamp(value: string | null | undefined): string;

	feedbackClass(type: string): string;
	reviewQueueItemClass(value: string): string;
	formatDate(value: string | null | undefined): string;
	formatReviewQueueDate(value: string | null | undefined): string;
	sequenceLabel(sequence: string | null | undefined): string;
};
