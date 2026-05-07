import { formatTimestamp } from './date';
import { normalizeTag } from './tags';
import type { JournalTag } from '../types/journal_api';
import type { DisplaySpec, PresentationConcern } from '../types/presentation_concern';

// --- Type ---

const defaultTypeSpec: DisplaySpec = { icon: '🏷️', text: '🏷️', class: 'journal-display-default' };

const typeDisplayMap: Record<string, DisplaySpec> = {
	TAKEN: { icon: '📈', text: '📈 TAKEN', class: 'journal-type-taken' },
	REJECTED: { icon: '📉', text: '📉 REJECTED', class: 'journal-type-rejected' },
	RESULT: { icon: '🏷️', text: '🏷️ RESULT', class: 'journal-type-result' },
	SET: { icon: '🏷️', text: '🏷️ SET', class: 'journal-type-set' },
};

// --- Status ---

const defaultStatusSpec: DisplaySpec = { icon: '🏷️', text: '🏷️', class: 'journal-display-default' };

const statusDisplayMap: Record<string, DisplaySpec> = {
	SUCCESS: { icon: '✅', text: '✅ SUCCESS', class: 'journal-status-success' },
	FAIL: { icon: '❌', text: '❌ FAIL', class: 'journal-status-fail' },
	RUNNING: { icon: '🏃', text: '🏃 RUNNING', class: 'journal-status-running' },
	SET: { icon: '🎯', text: '🎯 SET', class: 'journal-status-set' },
	JUST_LOSS: { icon: '💔', text: '💔 JUST_LOSS', class: 'journal-status-just-loss' },
	BROKEN: { icon: '💥', text: '💥 BROKEN', class: 'journal-status-broken' },
	MISSED: { icon: '🚫', text: '🚫 MISSED', class: 'journal-status-missed' },
	REJECTED: { icon: '🏷️', text: '🏷️ REJECTED', class: 'journal-status-rejected' },
};

// --- Timeframe ---

const defaultTimeframeSpec: DisplaySpec = { icon: '', text: '', class: 'journal-timeframe-default' };

const timeframeDisplayMap: Record<string, DisplaySpec> = {
	YR: { icon: '🗓️', text: '🗓️ YR', class: 'journal-timeframe-yr' },
	SMN: { icon: '📅', text: '📅 SMN', class: 'journal-timeframe-smn' },
	TMN: { icon: '📈', text: '📈 TMN', class: 'journal-timeframe-tmn' },
	MN: { icon: '📊', text: '📊 MN', class: 'journal-timeframe-mn' },
	WK: { icon: '📆', text: '📆 WK', class: 'journal-timeframe-wk' },
	DL: { icon: '🔍', text: '🔍 DL', class: 'journal-timeframe-dl' },
};

// --- Sequence ---

const sequenceDisplayMap: Record<string, DisplaySpec> = {
	MWD: { icon: '🕐', text: '🕐 MWD', class: '' },
	YR: { icon: '📅', text: '📅 YR', class: '' },
};

// --- Date ---

const shortMonthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];

// --- Tag helpers ---

function reasonTagIcon(tagName: string): string {
	return tagName.toLowerCase().includes('trend') ? '📈' : '⚡';
}

export function NewPresentationConcern(): PresentationConcern {
	return {
		// --- Type ---
		type(value: string): DisplaySpec {
			const key = normalizeTag(value);
			return typeDisplayMap[key] ?? { ...defaultTypeSpec, text: `${defaultTypeSpec.icon} ${key}` };
		},

		// --- Status ---
		status(value: string): DisplaySpec {
			const key = normalizeTag(value);
			return statusDisplayMap[key] ?? { ...defaultStatusSpec, text: `${defaultStatusSpec.icon} ${key}` };
		},

		// --- Timeframe ---
		timeframe(value: string): DisplaySpec {
			const key = normalizeTag(value);
			return timeframeDisplayMap[key] ?? { ...defaultTimeframeSpec, text: key };
		},

		// --- Sequence ---
		sequence(value: string | null | undefined): DisplaySpec {
			const key = normalizeTag(value ?? '');
			if (!key) return { icon: '', text: '', class: '' };
			return sequenceDisplayMap[key] ?? { icon: '📅', text: `📅 ${key}`, class: '' };
		},

		// --- Reason Tag ---
		reasonTag(tag: JournalTag): DisplaySpec {
			const name = tag.tag ?? '';
			const icon = reasonTagIcon(name);
			const override = tag.override ? ` → ${tag.override}` : '';
			return { icon, text: `${icon} ${name}${override}`, class: '' };
		},

		// --- Directional Tag ---
		directionalTag(tag: JournalTag): DisplaySpec {
			const name = tag.tag ?? '';
			const icon = '🏷';
			return { icon, text: `${icon} ${name}`, class: '' };
		},

		// --- Review State ---
		reviewedAt(value: string | null | undefined): DisplaySpec {
			const label = value ? formatTimestamp(value) : '—';
			return { icon: '✅', text: `✅ ${label}`, class: '' };
		},
		pendingReview(): DisplaySpec {
			return { icon: '⏳', text: '⏳ Pending Review', class: '' };
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
			return `${day} ${month}, '${year}`;
		},
	};
}
