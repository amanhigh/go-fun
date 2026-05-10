import { createSubmitter } from '../../../lib/submitter';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

function managementTagKey(value: string): string {
	return value.trim().toUpperCase();
}

export function TakenTagConcern(pg: JournalDetailPageProvider) {
	return {
		submitter: createSubmitter(),
		tags: [
			{ id: 'ntr', tag: 'ntr', type: 'MANAGEMENT' },
			{ id: 'enl', tag: 'enl', type: 'MANAGEMENT' },
			{ id: 'slt', tag: 'slt', type: 'MANAGEMENT' },
			{ id: 'fz', tag: 'fz', type: 'MANAGEMENT' },
			{ id: 'nbe', tag: 'nbe', type: 'MANAGEMENT' },
			{ id: 'ws', tag: 'ws', type: 'MANAGEMENT' },
			{ id: 'important', tag: 'important', type: 'MANAGEMENT' },
			{ id: 'be', tag: 'be', type: 'MANAGEMENT' },
		],

		show() {
			return pg().current.journal?.type === 'TAKEN';
		},
		hasTag(value: string) {
			const normalizedValue = managementTagKey(value);
			return pg().sidebar.tags.management().some((tag: JournalTag) => managementTagKey(tag.tag ?? '') === normalizedValue);
		},
		async submit(tagValue: string) {
			if (!pg().current.journal) return;
			await this.submitter.run(
				() => this.addTag(tagValue),
				{ success: `${managementTagKey(tagValue)} tag added.`, error: 'Unable to save management tag.' },
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
