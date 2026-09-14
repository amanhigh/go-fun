import { createSubmitter } from '../../../lib/submitter';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function NewHeaderConcern(pg: JournalDetailPageProvider) {
	return {
		submitter: createSubmitter(),

		async deleteJournal(this: any) {
			if (!pg().journal.detail) return;
			if (!window.confirm('Delete this journal? This cannot be undone.')) return;
			// No success notification: navigation happens immediately after
			// deletion, so a success notification could never be seen. Caught errors
			// are still surfaced by the Submitter as persistent notifications.
			await this.submitter.run(async () => {
				await pg().client.delete(pg().journal.detail!.id);
				window.location.href = '/journal';
			});
		},
	};
}
