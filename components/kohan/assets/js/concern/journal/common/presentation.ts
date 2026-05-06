import { formatTimestamp } from '../../../shared/date';
import { normalizeTag } from '../../../shared/tags';
import type { JournalPresentationState } from '../../../types/journal_common_state';

const defaultBadgeClass = 'border-slate-300 bg-slate-50 text-slate-700';

const statusBadgeClassMap: Record<string, string> = {
	SUCCESS: 'border-emerald-300 bg-emerald-50 text-emerald-800',
	FAIL: 'border-rose-300 bg-rose-50 text-rose-800',
	RUNNING: 'border-sky-300 bg-sky-50 text-sky-800',
	SET: 'border-amber-300 bg-amber-50 text-amber-800',
	REJECTED: 'border-violet-300 bg-violet-50 text-violet-800',
};

const typeBadgeClassMap: Record<string, string> = {
	REJECTED: 'border-rose-300 bg-rose-50 text-rose-800',
	RESULT: 'border-emerald-300 bg-emerald-50 text-emerald-800',
	SET: 'border-indigo-300 bg-indigo-50 text-indigo-800',
};

function resolveBadgeClass(map: Record<string, string>, value: string): string {
	return map[normalizeTag(value)] ?? defaultBadgeClass;
}

export function newPresentationConcern(): JournalPresentationState {
	return {
		normalizeStatus: normalizeTag,
		statusBadgeClass(value: string) {
			return resolveBadgeClass(statusBadgeClassMap, value);
		},
		typeBadgeClass(value: string) {
			return resolveBadgeClass(typeBadgeClassMap, value);
		},
		formatTimestamp,
	};
}
