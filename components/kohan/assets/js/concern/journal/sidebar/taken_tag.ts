import { createSubmitter } from '../../../lib/submitter';
import { JournalType, JournalTagType } from '../../../types/api/journal/enums';
import type { JournalTag } from '../../../types/api/journal/response';
import type { JournalTagRequest } from '../../../types/api/journal/request';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function TakenTagConcern(pg: JournalDetailPageProvider) {
	return {
		submitter: createSubmitter(),
		tags: [
			{ id: 'ntr', tag: 'ntr', type: JournalTagType.MANAGEMENT, created_at: '' },
			{ id: 'enl', tag: 'enl', type: JournalTagType.MANAGEMENT, created_at: '' },
			{ id: 'slt', tag: 'slt', type: JournalTagType.MANAGEMENT, created_at: '' },
			{ id: 'fz', tag: 'fz', type: JournalTagType.MANAGEMENT, created_at: '' },
			{ id: 'nbe', tag: 'nbe', type: JournalTagType.MANAGEMENT, created_at: '' },
			{ id: 'ws', tag: 'ws', type: JournalTagType.MANAGEMENT, created_at: '' },
			{ id: 'important', tag: 'important', type: JournalTagType.MANAGEMENT, created_at: '' },
			{ id: 'be', tag: 'be', type: JournalTagType.MANAGEMENT, created_at: '' },
		],

		show() {
			return pg().current.journal?.type === JournalType.TAKEN;
		},
		hasTag(value: string) {
			return pg().sidebar.tags.management().some((tag: JournalTag) => tag.tag === value);
		},
		async submit(tagValue: string) {
			if (!pg().current.journal) return;
			await this.submitter.run(
				() => this.addTag(tagValue),
				{ success: `${tagValue} tag added.` },
			);
		},

		async addTag(tagValue: string) {
			const payload: JournalTagRequest = { tag: tagValue, type: JournalTagType.MANAGEMENT };
			const page = pg();
			const envelope = await page.tagClient.create(page.current.journalId, payload);
			page.sidebar.tags.prepend(envelope.data as JournalTag);
		},
	};
}
