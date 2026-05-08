import { createAsyncFeedback } from '../../../lib/async_feedback';
import { normalizeTag } from '../../../lib/tags';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function MgmntTagConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedback(),
		presets: [
			{ icon: '', text: 'ntr', class: 'journal-management-base-emerald' },
			{ icon: '', text: 'enl', class: 'journal-management-base-sky' },
			{ icon: '', text: 'slt', class: 'journal-management-base-rose' },
			{ icon: '', text: 'fz', class: 'journal-management-base-violet' },
			{ icon: '', text: 'nbe', class: 'journal-management-base-amber' },
			{ icon: '', text: 'ws', class: 'journal-management-base-slate' },
			{ icon: '', text: 'important', class: 'journal-management-base-fuchsia' },
			{ icon: '', text: 'be', class: 'journal-management-base-orange' },
		],
		pendingValue: '',

		showTakenTags() {
			return normalizeTag(pg().current.journal?.type ?? '') === 'TAKEN';
		},
		hasTag(value: string) {
			const normalizedValue = normalizeTag(value);
			return pg().sidebar.tags.management().some((tag: JournalTag) => normalizeTag(tag.tag ?? '') === normalizedValue);
		},
		async submit(tagValue: string) {
			if (!pg().current.journal || this.submitting) return;
			this.pendingValue = tagValue;
			await this.run(
				() => this.addTag(tagValue),
				`${normalizeTag(tagValue)} tag added.`,
				'Unable to save management tag.',
			);
			this.pendingValue = '';
		},

		async addTag(tagValue: string) {
			const payload: JournalTagRequest = { tag: tagValue, type: 'MANAGEMENT' };
			const page = pg();
			const envelope = await page.tagClient.create(page.current.journalId, payload);
			page.sidebar.tags.prepend(envelope.data as JournalTag);
		},
	};
}
