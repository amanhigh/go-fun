import { formatTimestamp } from '../shared/date';
import { normalizeTag } from '../shared/tags';
import { resolveTypeToggle } from './filter_actions';
import type { JournalFilterState } from './filter';

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

export function createJournalListFormatters(filter: JournalFilterState) {
	return {
		normalizeStatus: normalizeTag,
		statusBadgeClass(value: string) {
			return statusBadgeClassMap[normalizeTag(value)] ?? defaultBadgeClass;
		},
		typeBadgeClass(value: string) {
			return typeBadgeClassMap[normalizeTag(value)] ?? defaultBadgeClass;
		},
		typeToggleLabel() {
			return resolveTypeToggle(filter.type).label;
		},
		typeToggleClass() {
			return resolveTypeToggle(filter.type).className;
		},
		formatTimestamp,
	};
}
