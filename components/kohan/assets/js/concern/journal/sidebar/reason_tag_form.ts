import { createAsyncFeedback } from '../../../lib/async_feedback';
import type { JournalTag, JournalTagRequest } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewReasonTagFormConcern(pg: JournalDetailPageProvider) {
	return {
		...createAsyncFeedback(),
		input: '',
		override: '',

		focusOverride() {
			(pg() as any).$nextTick?.(() => {
				(pg() as any).$refs?.reasonTagOverride?.focus?.();
			});
		},

		async submit() {
			if (!pg().current.journal || this.submitting) return;
			const tag = this.input.trim();
			if (!tag) {
				this.setError('Tag is required.');
				return;
			}
			const override = this.override.trim();
			await this.run(async () => {
				const payload: JournalTagRequest = {
					tag,
					type: 'REASON',
					...(override ? { override } : {}),
				};
				const envelope = await pg().tagClient.create(pg().current.journalId, payload);
				pg().sidebar.tags.prepend(envelope.data as JournalTag);
				this.input = '';
				this.override = '';
			}, 'Reason tag added.', 'Unable to save reason tag.');
		},
	};
}
