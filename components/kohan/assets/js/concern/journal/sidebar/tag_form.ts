import { createSubmitter } from '../../../lib/submitter';
import { JournalTagType } from '../../../types/api/journal/enums';
import type { JournalTag } from '../../../types/api/journal/response';
import type { JournalTagRequest } from '../../../types/api/journal/request';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function TagFormConcern(pg: JournalDetailPageProvider) {
	return {
		submitter: createSubmitter(),
		input: '',
		override: '',

		canSubmit() {
			return this.input.trim() !== '';
		},

		async submit() {
			const page = pg();
			const journal = page.journal.detail;
			if (!journal) return;

			const tag = this.input.trim();
			if (!tag) {
				this.submitter.setError('Tag is required.');
				return;
			}

			const override = this.override.trim();

			await this.submitter.run(
				() => this.createTag(tag, override),
				{ success: 'Reason tag added.' },
			);
		},

		async createTag(tag: string, override: string) {
			const page = pg();
			const payload: JournalTagRequest = { tag, type: JournalTagType.REASON };
			if (override) {
				payload.override = override;
			}
			const envelope = await page.tagClient.create(page.journal.detail!.id, payload);
			page.sidebar.tags.prepend(envelope.data as JournalTag);
			this.input = '';
			this.override = '';
		},
	};
}
