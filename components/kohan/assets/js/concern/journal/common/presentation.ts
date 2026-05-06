import { formatTimestamp } from '../../../shared/date';
import { normalizeTag } from '../../../shared/tags';
import type { JournalTag } from '../../../types/journal_api';
import type { DisplaySpec, PresentationConcern } from '../../../types/presentation_concern';

// --- Type ---

const defaultTypeSpec: DisplaySpec = { icon: '🏷️', text: '🏷️', badgeClass: 'journal-display-default' };

const typeDisplayMap: Record<string, DisplaySpec> = {
	TAKEN: { icon: '📈', text: '📈 TAKEN', badgeClass: 'journal-type-taken' },
	REJECTED: { icon: '📉', text: '📉 REJECTED', badgeClass: 'journal-type-rejected' },
	RESULT: { icon: '🏷️', text: '🏷️ RESULT', badgeClass: 'journal-type-result' },
	SET: { icon: '🏷️', text: '🏷️ SET', badgeClass: 'journal-type-set' },
};

// --- Status ---

const defaultStatusSpec: DisplaySpec = { icon: '🏷️', text: '🏷️', badgeClass: 'journal-display-default' };

const statusDisplayMap: Record<string, DisplaySpec> = {
	SUCCESS: { icon: '✅', text: '✅ SUCCESS', badgeClass: 'journal-status-success' },
	FAIL: { icon: '❌', text: '❌ FAIL', badgeClass: 'journal-status-fail' },
	RUNNING: { icon: '🏃', text: '🏃 RUNNING', badgeClass: 'journal-status-running' },
	SET: { icon: '🎯', text: '🎯 SET', badgeClass: 'journal-status-set' },
	JUST_LOSS: { icon: '💔', text: '💔 JUST_LOSS', badgeClass: 'journal-status-just-loss' },
	BROKEN: { icon: '💥', text: '💥 BROKEN', badgeClass: 'journal-status-broken' },
	MISSED: { icon: '🚫', text: '🚫 MISSED', badgeClass: 'journal-status-missed' },
	REJECTED: { icon: '🏷️', text: '🏷️ REJECTED', badgeClass: 'journal-status-rejected' },
};

// --- Timeframe ---

const defaultTimeframeSpec: DisplaySpec = { icon: '', text: '', badgeClass: 'journal-timeframe-default' };

const timeframeDisplayMap: Record<string, DisplaySpec> = {
	YR: { icon: '🗓️', text: '🗓️ YR', badgeClass: 'journal-timeframe-yr' },
	SMN: { icon: '📅', text: '📅 SMN', badgeClass: 'journal-timeframe-smn' },
	TMN: { icon: '📈', text: '📈 TMN', badgeClass: 'journal-timeframe-tmn' },
	MN: { icon: '📊', text: '📊 MN', badgeClass: 'journal-timeframe-mn' },
	WK: { icon: '📆', text: '📆 WK', badgeClass: 'journal-timeframe-wk' },
	DL: { icon: '🔍', text: '🔍 DL', badgeClass: 'journal-timeframe-dl' },
};

// --- Sequence ---

const sequenceDisplayMap: Record<string, { icon: string; text: string }> = {
	MWD: { icon: '🕐', text: '🕐 MWD' },
};

const defaultSequenceText = '📅 ';

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
			if (!key) return { icon: '', text: '', badgeClass: '' };
			const spec = sequenceDisplayMap[key];
			if (spec) return { ...spec, badgeClass: '' };
			return { icon: '📅', text: `${defaultSequenceText}${key}`, badgeClass: '' };
		},

		// --- Reason Tag ---
		reasonTag(tag: JournalTag): DisplaySpec {
			const name = tag.tag ?? '';
			const icon = reasonTagIcon(name);
			const override = tag.override ? ` → ${tag.override}` : '';
			return { icon, text: `${icon} ${name}${override}`, badgeClass: '' };
		},

		// --- Directional Tag ---
		directionalTag(tag: JournalTag): DisplaySpec {
			const name = tag.tag ?? '';
			return { icon: '', text: name, badgeClass: '' };
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
	};
}
