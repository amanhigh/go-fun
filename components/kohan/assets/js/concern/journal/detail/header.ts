import { createFeedback } from '../../../lib/feedback';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewHeaderConcern(pg: JournalDetailPageProvider) {
	return {
		...createFeedback(),

		async deleteJournal(this: any) {
			if (!pg().current.journal || this.submitting) return;
			if (!window.confirm('Delete this journal? This cannot be undone.')) return;
			await this.run(
				async () => {
					await pg().client.delete(pg().current.journalId);
					window.location.href = '/journal';
				},
				'Journal deleted.',
				'Unable to delete journal.',
			);
			if (this.feedbackKind === 'error') {
				pg().current.errorMessage = this.message;
			}
		},
	};
}
