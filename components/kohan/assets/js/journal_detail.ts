import { type Journal, type Envelope } from './journal_models';

declare const Alpine: {
	data(name: string, callback: () => ReturnType<typeof journalDetailPage>): void;
};

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

function journalDetailPage() {
	return {
		journalId: '',
		journal: null as Journal | null,
		loading: false,
		errorMessage: '',
		normalizeStatus(value: string) {
			return normalizeTag(value);
		},
		statusBadgeClass(value: string) {
			return statusBadgeClassMap[normalizeTag(value)] ?? defaultBadgeClass;
		},
		typeBadgeClass(value: string) {
			return typeBadgeClassMap[normalizeTag(value)] ?? defaultBadgeClass;
		},
		async loadJournal() {
			this.loading = true;
			this.errorMessage = '';
			try {
				const response = await fetch(`/v1/api/journals/${this.journalId}`);
				if (!response.ok) {
					if (response.status === 404) {
						throw new Error('Journal not found');
					}
					throw new Error('Failed to load journal');
				}
				const envelope = (await response.json()) as Envelope<Journal>;
				this.journal = envelope.data ?? null;
			} catch (err) {
				this.errorMessage = err instanceof Error ? err.message : 'Unable to load journal details. Please try again.';
			} finally {
				this.loading = false;
			}
		},
		hasError() {
			return this.errorMessage !== '';
		},
		resolveImageSrc(fileName: string) {
			if (!fileName) return '';
			if (fileName.startsWith('http://') || fileName.startsWith('https://') || fileName.startsWith('/')) {
				return fileName;
			}
			return '/assets/' + fileName;
		},
		formatTimestamp(value: string | null | undefined) {
			if (!value) return '—';
			const parsed = new Date(value);
			if (Number.isNaN(parsed.getTime())) return '—';
			return parsed.toLocaleString();
		},
	};
}

document.addEventListener('alpine:init', () => {
	Alpine.data('journalDetailPage', journalDetailPage);
});

export {};
