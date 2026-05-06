import { formatTimestamp } from '../../../shared/date';
import { normalizeTag } from '../../../shared/tags';
import type { JournalTag } from '../../../types/journal_api';
import type { DisplaySpec, PresentationConcern } from '../../../types/presentation_concern';

// --- Type ---

const defaultTypeSpec: DisplaySpec = { icon: '🏷️', badgeClass: 'border-slate-300 bg-slate-50 text-slate-700' };

const typeDisplayMap: Record<string, DisplaySpec> = {
	TAKEN: { icon: '📈', badgeClass: 'border-slate-300 bg-slate-50 text-slate-700' },
	REJECTED: { icon: '📉', badgeClass: 'border-rose-300 bg-rose-50 text-rose-800' },
	RESULT: { icon: '🏷️', badgeClass: 'border-emerald-300 bg-emerald-50 text-emerald-800' },
	SET: { icon: '🏷️', badgeClass: 'border-indigo-300 bg-indigo-50 text-indigo-800' },
};

// --- Status ---

const defaultStatusSpec: DisplaySpec = { icon: '🏷️', badgeClass: 'border-slate-300 bg-slate-50 text-slate-700' };

const statusDisplayMap: Record<string, DisplaySpec> = {
	SUCCESS: { icon: '✅', badgeClass: 'border-emerald-300 bg-emerald-50 text-emerald-800' },
	FAIL: { icon: '❌', badgeClass: 'border-rose-300 bg-rose-50 text-rose-800' },
	RUNNING: { icon: '🏃', badgeClass: 'border-sky-300 bg-sky-50 text-sky-800' },
	SET: { icon: '🎯', badgeClass: 'border-amber-300 bg-amber-50 text-amber-800' },
	JUST_LOSS: { icon: '💔', badgeClass: 'border-rose-300 bg-rose-50 text-rose-800' },
	BROKEN: { icon: '💥', badgeClass: 'border-violet-300 bg-violet-50 text-violet-800' },
	MISSED: { icon: '🚫', badgeClass: 'border-slate-300 bg-slate-50 text-slate-700' },
	REJECTED: { icon: '🏷️', badgeClass: 'border-violet-300 bg-violet-50 text-violet-800' },
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

// --- Review Queue ---

const reviewQueueItemClassMap: Record<string, string> = {
	TAKEN: 'border-emerald-300 bg-emerald-50/70 hover:bg-emerald-100/80 text-emerald-900',
	REJECTED: 'border-rose-300 bg-rose-50/70 hover:bg-rose-100/80 text-rose-900',
};

const defaultReviewQueueItemClass = 'border-border bg-muted/30 hover:bg-muted/70 hover:text-foreground';

// --- Date ---

const shortMonthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

// --- Tag helpers ---

function reasonTagIcon(tagName: string): string {
	return tagName.toLowerCase().includes('trend') ? '📈 ' : '⚡ ';
}

export function NewPresentationConcern(): PresentationConcern {
	return {
		// --- Type ---
		typeBadgeClass(value: string) {
			return (typeDisplayMap[normalizeTag(value)] ?? defaultTypeSpec).badgeClass;
		},
		typeDisplay(value: string) {
			const key = normalizeTag(value);
			const spec = typeDisplayMap[key] ?? defaultTypeSpec;
			return `${spec.icon} ${key}`;
		},

		// --- Status ---
		statusBadgeClass(value: string) {
			return (statusDisplayMap[normalizeTag(value)] ?? defaultStatusSpec).badgeClass;
		},
		statusDisplay(value: string) {
			const key = normalizeTag(value);
			const spec = statusDisplayMap[key] ?? defaultStatusSpec;
			return `${spec.icon} ${key}`;
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
			const override = tag.override ? ` → ${tag.override}` : '';
			return `${reasonTagIcon(name)}${name}${override}`;
		},
		directionalTagLabel(tag: JournalTag) {
			return tag.tag ?? '';
		},

		// --- Timestamp / Date ---
		formatTimestamp,
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
			return reviewQueueItemClassMap[normalizeTag(value)] ?? defaultReviewQueueItemClass;
		},
	};
}
