import { formatTimestamp } from '../../../shared/date';
import { normalizeTag } from '../../../shared/tags';
import type { PresentationConcern } from '../../../types/presentation_concern';

const defaultBadgeClass = 'border-slate-300 bg-slate-50 text-slate-700';

const statusBadgeClassMap: Record<string, string> = {
	SUCCESS: 'border-emerald-300 bg-emerald-50 text-emerald-800',
	FAIL: 'border-rose-300 bg-rose-50 text-rose-800',
	RUNNING: 'border-sky-300 bg-sky-50 text-sky-800',
	SET: 'border-amber-300 bg-amber-50 text-amber-800',
	JUST_LOSS: 'border-rose-300 bg-rose-50 text-rose-800',
	BROKEN: 'border-violet-300 bg-violet-50 text-violet-800',
	MISSED: 'border-slate-300 bg-slate-50 text-slate-700',
	REJECTED: 'border-violet-300 bg-violet-50 text-violet-800',
};

const typeBadgeClassMap: Record<string, string> = {
	REJECTED: 'border-rose-300 bg-rose-50 text-rose-800',
	RESULT: 'border-emerald-300 bg-emerald-50 text-emerald-800',
	SET: 'border-indigo-300 bg-indigo-50 text-indigo-800',
};

const shortMonthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

function resolveBadgeClass(map: Record<string, string>, value: string): string {
	return map[normalizeTag(value)] ?? defaultBadgeClass;
}

export function NewPresentationConcern(): PresentationConcern {
	return {
		normalizeStatus: normalizeTag,
		statusBadgeClass(value: string) {
			return resolveBadgeClass(statusBadgeClassMap, value);
		},
		typeBadgeClass(value: string) {
			return resolveBadgeClass(typeBadgeClassMap, value);
		},
		formatTimestamp,
		// Extended formatters
		feedbackClass(type: string) {
			return type === 'success' ? 'text-emerald-700' : 'text-rose-700';
		},
		reviewQueueItemClass(value: string) {
			const journalType = normalizeTag(value);
			if (journalType === 'TAKEN') {
				return 'border-emerald-300 bg-emerald-50/70 hover:bg-emerald-100/80 text-emerald-900';
			}
			if (journalType === 'REJECTED') {
				return 'border-rose-300 bg-rose-50/70 hover:bg-rose-100/80 text-rose-900';
			}
			return 'border-border bg-muted/30 hover:bg-muted/70 hover:text-foreground';
		},
		formatDate(value: string | null | undefined) {
			if (!value) return '—';
			const parsed = new Date(value);
			return Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleDateString();
		},
		formatReviewQueueDate(value: string | null | undefined) {
			if (!value) return '—';
			const parsed = new Date(value);
			if (Number.isNaN(parsed.getTime())) return '—';
			const day = parsed.getUTCDate();
			const month = shortMonthNames[parsed.getUTCMonth()] ?? '—';
			const year = `${parsed.getUTCFullYear()}`.slice(-2);
			return `${day} ${month}, ${year}`;
		},
		sequenceLabel(sequence: string | null | undefined) {
			if (!sequence) return '';
			return sequence === 'MWD' ? '🕐 ' + sequence : '📅 ' + sequence;
		},
	};
}
