import { createSubmitter } from '../../../lib/submitter';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

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
