import { createSubmitter } from '../../../lib/submitter';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewHeaderConcern(pg: JournalDetailPageProvider) {
	return {
		submitter: createSubmitter(),

		async deleteJournal(this: any) {
			if (!pg().current.journal) return;
			if (!window.confirm('Delete this journal? This cannot be undone.')) return;
			await this.submitter.run(
				async () => {
					await pg().client.delete(pg().current.journalId);
					window.location.href = '/journal';
				},
				{ success: 'Journal deleted.', error: 'Unable to delete journal.' },
			);
		},
	};
}
