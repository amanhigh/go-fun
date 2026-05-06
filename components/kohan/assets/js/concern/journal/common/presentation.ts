import { formatTimestamp } from '../../../shared/date';
import { normalizeTag } from '../../../shared/tags';
import type { JournalTag } from '../../../types/journal_api';
import type { PresentationConcern } from '../../../types/presentation_concern';

const defaultBadgeClass = 'border-slate-300 bg-slate-50 text-slate-700';

function resolveBadgeClass(map: Record<string, string>, value: string): string {
	return map[normalizeTag(value)] ?? defaultBadgeClass;
}

// --- Type ---

const typeBadgeClassMap: Record<string, string> = {
	REJECTED: 'border-rose-300 bg-rose-50 text-rose-800',
	RESULT: 'border-emerald-300 bg-emerald-50 text-emerald-800',
	SET: 'border-indigo-300 bg-indigo-50 text-indigo-800',
};

const typeIconMap: Record<string, string> = {
	TAKEN: '📈',
	REJECTED: '📉',
};

// --- Status ---

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

const statusIconMap: Record<string, string> = {
	RUNNING: '🏃',
	SET: '🎯',
	SUCCESS: '✅',
	FAIL: '❌',
	BROKEN: '💥',
	MISSED: '🚫',
	JUST_LOSS: '💔',
};

// --- Timeframe ---

const timeframeChipClassMap: Record<string, string> = {
	YR: 'border-fuchsia-400 bg-fuchsia-200 text-fuchsia-950',
	SMN: 'border-sky-400 bg-sky-200 text-sky-950',
	TMN: 'border-emerald-400 bg-emerald-200 text-emerald-950',
	MN: 'border-orange-400 bg-orange-200 text-orange-950',
	WK: 'border-yellow-400 bg-yellow-200 text-yellow-950',
	DL: 'border-slate-400 bg-slate-200 text-slate-950',
};

// --- Date ---

const shortMonthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

export function NewPresentationConcern(): PresentationConcern {
	return {
		// --- Type ---
		typeBadgeClass(value: string) {
			return resolveBadgeClass(typeBadgeClassMap, value);
		},
		typeDisplay(value: string) {
			const icon = typeIconMap[normalizeTag(value)] ?? '🏷️';
			return `${icon} ${value}`;
		},

		// --- Status ---
		normalizeStatus: normalizeTag,
		statusBadgeClass(value: string) {
			return resolveBadgeClass(statusBadgeClassMap, value);
		},
		statusDisplay(value: string) {
			const icon = statusIconMap[normalizeTag(value)] ?? '🏷️';
			return `${icon} ${value}`;
		},

		// --- Timeframe ---
		timeframeChipClass(value: string) {
			return timeframeChipClassMap[normalizeTag(value)] ?? 'border-zinc-300 bg-zinc-100 text-zinc-900';
		},

		// --- Sequence ---
		sequenceLabel(sequence: string | null | undefined) {
			if (!sequence) return '';
			return sequence === 'MWD' ? '🕐 ' + sequence : '📅 ' + sequence;
		},

		// --- Tag Labels ---
		reasonTagLabel(tag: JournalTag) {
			const name = tag.tag ?? '';
			const prefix = name.toLowerCase().includes('trend') ? '📈 ' : '⚡ ';
			const override = tag.override ? ` → ${tag.override}` : '';
			return `${prefix}${name}${override}`;
		},
		directionalTagLabel(tag: JournalTag) {
			return tag.tag ?? '';
		},

		// --- Timestamp / Date ---
		formatTimestamp,
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

		// --- Review Queue ---
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

		// --- Feedback ---
		feedbackClass(type: string) {
			return type === 'success' ? 'text-emerald-700' : 'text-rose-700';
		},
	};
}
