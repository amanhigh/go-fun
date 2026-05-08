import { createAsyncFeedback } from '../../../lib/async_feedback';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';
import type { AlpineContext } from '../../../types/platform';

export function TagFormConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedback(),
		input: '',
		override: '',

		focusOverride() {
			const alpine = pg() as unknown as AlpineContext;
			alpine.$nextTick?.(() => {
				alpine.$refs?.reasonTagOverride?.focus?.();
			});
		},

		async submit() {
			const page = pg();
			const journal = page.current.journal;
			if (!journal) return;

			const tag = this.input;
			if (!tag) {
				this.setError('Tag is required.');
				return;
			}

			const override = this.override;

			await this.run(
				() => this.createTag(tag, override),
				'Reason tag added.',
				'Unable to save reason tag.',
			);
		},

		async createTag(tag: string, override: string) {
			const page = pg();
			const payload: JournalTagRequest = { tag, type: 'REASON' };
			if (override) {
				payload.override = override;
			}
			const envelope = await page.tagClient.create(page.current.journalId, payload);
			page.sidebar.tags.prepend(envelope.data as JournalTag);
			this.input = '';
			this.override = '';
		},
	};
}
