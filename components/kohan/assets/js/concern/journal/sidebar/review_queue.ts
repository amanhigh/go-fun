import { getErrorMessage } from '../../../shared/error';
import type { Journal } from '../../../types/journal_api';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewReviewQueueConcern(pg: JournalDetailPageProvider) {
	return {
		items: [] as Journal[],
		loading: false,
		error: '',

		async load() {
			this.loading = true;
			this.error = '';
			try {
				const envelope = await pg().client.list(0, 10, { reviewed: 'false', sortBy: 'created_at', sortOrder: 'asc' });
				this.items = envelope.data?.journals ?? [];
			} catch (err) {
				this.error = getErrorMessage(err, 'Unable to load review queue.');
			} finally {
				this.loading = false;
			}
		},
	};
}
