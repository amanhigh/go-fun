import { createFeedback } from '../../../lib/feedback';
import { normalizeTag } from '../../../lib/tags';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function TakenTagConcern(pg: JournalDetailPageProvider) {
	return {
		...createFeedback(),
		tags: [
			{ text: 'ntr', class: 'journal-management-base-emerald' },
			{ text: 'enl', class: 'journal-management-base-sky' },
			{ text: 'slt', class: 'journal-management-base-rose' },
			{ text: 'fz', class: 'journal-management-base-violet' },
			{ text: 'nbe', class: 'journal-management-base-amber' },
			{ text: 'ws', class: 'journal-management-base-slate' },
			{ text: 'important', class: 'journal-management-base-fuchsia' },
			{ text: 'be', class: 'journal-management-base-orange' },
		],

		show() {
			return normalizeTag(pg().current.journal?.type ?? '') === 'TAKEN';
		},
		hasTag(value: string) {
			const normalizedValue = normalizeTag(value);
			return pg().sidebar.tags.management().some((tag: JournalTag) => normalizeTag(tag.tag ?? '') === normalizedValue);
		},
		async submit(tagValue: string) {
			if (!pg().current.journal || this.submitting) return;
			await this.run(
				() => this.addTag(tagValue),
				`${normalizeTag(tagValue)} tag added.`,
				'Unable to save management tag.',
			);
		},

		async addTag(tagValue: string) {
			const payload: JournalTagRequest = { tag: tagValue, type: 'MANAGEMENT' };
			const page = pg();
			const envelope = await page.tagClient.create(page.current.journalId, payload);
			page.sidebar.tags.prepend(envelope.data as JournalTag);
		},
	};
}
