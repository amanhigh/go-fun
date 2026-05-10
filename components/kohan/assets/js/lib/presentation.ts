import type { JournalTag, JournalTimeframe } from '../types/journal_api';
import type { DisplaySpec, PresentationConcern } from '../types/presentation_concern';

function display(spec: DisplaySpec): string {
	return spec.icon ? `${spec.icon} ${spec.text}` : spec.text;
}

// --- Type ---

const defaultTypeSpec: DisplaySpec = { icon: '🏷️', text: 'Unknown', class: 'journal-display-default' };

const typeDisplayMap: Record<string, DisplaySpec> = {
	TAKEN: { icon: '📈', text: 'TAKEN', class: 'journal-type-taken' },
	REJECTED: { icon: '📉', text: 'REJECTED', class: 'journal-type-rejected' },
	RESULT: { icon: '🏷️', text: 'RESULT', class: 'journal-type-result' },
	SET: { icon: '🏷️', text: 'SET', class: 'journal-type-set' },
};

// --- Status ---

const defaultStatusSpec: DisplaySpec = { icon: '🏷️', text: 'Unknown', class: 'journal-display-default' };

const statusDisplayMap: Record<string, DisplaySpec> = {
	SUCCESS: { icon: '✅', text: 'SUCCESS', class: 'journal-status-success' },
	FAIL: { icon: '❌', text: 'FAIL', class: 'journal-status-fail' },
	RUNNING: { icon: '🏃', text: 'RUNNING', class: 'journal-status-running' },
	SET: { icon: '🎯', text: 'SET', class: 'journal-status-set' },
	JUST_LOSS: { icon: '💔', text: 'JUST_LOSS', class: 'journal-status-just-loss' },
	BROKEN: { icon: '💥', text: 'BROKEN', class: 'journal-status-broken' },
	MISSED: { icon: '🚫', text: 'MISSED', class: 'journal-status-missed' },
	REJECTED: { icon: '🏷️', text: 'REJECTED', class: 'journal-status-rejected' },
};

// --- Timeframe ---

const defaultTimeframeSpec: DisplaySpec = { text: '', class: 'journal-timeframe-default' };

const timeframeDisplayMap: Record<JournalTimeframe, DisplaySpec> = {
	YR: { icon: '🗓️', text: 'YR', class: 'journal-timeframe-yr' },
	SMN: { icon: '📅', text: 'SMN', class: 'journal-timeframe-smn' },
	TMN: { icon: '📈', text: 'TMN', class: 'journal-timeframe-tmn' },
	MN: { icon: '📊', text: 'MN', class: 'journal-timeframe-mn' },
	WK: { icon: '📆', text: 'WK', class: 'journal-timeframe-wk' },
	DL: { icon: '🔍', text: 'DL', class: 'journal-timeframe-dl' },
};

// --- Tag helpers ---

function reasonTagIcon(tagName: string): string {
	return tagName.toLowerCase().includes('trend') ? '📈' : '⚡';
}

export function NewPresentationConcern(): PresentationConcern {
	return {
		display,

		// --- Type ---
		type(value: string): DisplaySpec {
			return typeDisplayMap[value] ?? { ...defaultTypeSpec, text: value };
		},

		// --- Status ---
		status(value: string): DisplaySpec {
			return statusDisplayMap[value] ?? { ...defaultStatusSpec, text: value };
		},

		// --- Timeframe ---
		timeframe(value: string): DisplaySpec {
			return timeframeDisplayMap[value as JournalTimeframe] ?? { text: value, class: 'journal-timeframe-default' };
		},

		// --- Sequence ---
		sequence(value: string | null | undefined): DisplaySpec {
			const key = value ?? '';
			if (!key) return { text: '', class: '' };
			const seqCatalog: Record<string, DisplaySpec> = {
				MWD: { icon: '🕐', text: 'MWD', class: '' },
				YR: { icon: '📅', text: 'YR', class: '' },
			};
			return seqCatalog[key] ?? { icon: '📅', text: key, class: '' };
		},

		// --- Reason Tag ---
		reasonTag(tag: JournalTag): DisplaySpec {
			const name = tag.tag ?? '';
			const icon = reasonTagIcon(name);
			const override = tag.override ? ` → ${tag.override}` : '';
			return { icon, text: `${name}${override}`, class: '' };
		},

		// --- Directional Tag ---
		directionalTag(tag: JournalTag): DisplaySpec {
			const name = tag.tag ?? '';
			return { icon: '🏷', text: name, class: '' };
		},

		// --- Review State ---
		reviewedAt(value: string | null | undefined): DisplaySpec {
			let label = '—';
			if (value) {
				const parsed = new Date(value);
				label = Number.isNaN(parsed.getTime()) ? '—' : parsed.toLocaleString();
			}
			return { icon: '✅', text: label, class: '' };
		},
		pendingReview(): DisplaySpec {
			return { icon: '⏳', text: 'Pending Review', class: '' };
		},
	};
}
