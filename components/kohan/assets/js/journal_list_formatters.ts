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

const normalizeTag = (value: string): string => (value ?? '').trim().toUpperCase();

export function createJournalListFormatters() {
	return {
		normalizeStatus(value: string) {
			return normalizeTag(value);
		},
		statusBadgeClass(value: string) {
			return statusBadgeClassMap[normalizeTag(value)] ?? defaultBadgeClass;
		},
		typeBadgeClass(value: string) {
			return typeBadgeClassMap[normalizeTag(value)] ?? defaultBadgeClass;
		},
		typeToggleButtonLabel(this: any) {
			return this.filterTracker.type === 'TAKEN' ? 'Rejected' : 'Taken';
		},
		typeToggleButtonClass(this: any) {
			return this.typeToggleButtonLabel() === 'Taken'
				? '!border-emerald-300 !bg-emerald-200 !text-emerald-800'
				: '!border-rose-300 !bg-rose-200 !text-rose-800';
		},
		formatTimestamp(value: string) {
			if (!value) return '—';
			const parsed = new Date(value);
			if (Number.isNaN(parsed.getTime())) return '—';
			return parsed.toLocaleString();
		},
	};
}
