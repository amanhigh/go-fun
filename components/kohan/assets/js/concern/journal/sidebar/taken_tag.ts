import { createSubmitter } from '../../../lib/submitter';
import { JournalType, JournalTagType } from '../../../types/api/journal/enums';
import type { JournalTag } from '../../../types/api/journal/response';
import type { JournalTagRequest } from '../../../types/api/journal/request';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function TakenTagConcern(pg: JournalDetailPageProvider) {
	return {
		submitter: createSubmitter(),
		tags: [
			{ id: 'ntr', tag: 'ntr', type: JournalTagType.MANAGEMENT },
			{ id: 'enl', tag: 'enl', type: JournalTagType.MANAGEMENT },
			{ id: 'slt', tag: 'slt', type: JournalTagType.MANAGEMENT },
			{ id: 'fz', tag: 'fz', type: JournalTagType.MANAGEMENT },
			{ id: 'nbe', tag: 'nbe', type: JournalTagType.MANAGEMENT },
			{ id: 'ws', tag: 'ws', type: JournalTagType.MANAGEMENT },
			{ id: 'important', tag: 'important', type: JournalTagType.MANAGEMENT },
			{ id: 'be', tag: 'be', type: JournalTagType.MANAGEMENT },
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
				{ success: `${tagValue} tag added.`, error: 'Unable to save management tag.' },
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
