import { getErrorMessage } from '../../../shared/error';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewHeaderConcern(pg: JournalDetailPageProvider) {
	return {
		deleting: false,

		async deleteJournal(this: any) {
			if (!pg().current.journal || this.deleting) return;
			if (!window.confirm('Delete this journal? This cannot be undone.')) return;
			this.deleting = true;
			try {
				await pg().client.delete(pg().current.journalId);
				window.location.href = '/journal';
			} catch (err) {
				pg().current.errorMessage = getErrorMessage(err, 'Unable to delete journal.');
			} finally {
				this.deleting = false;
			}
		},
	};
}
