import { createCollection } from '../../../lib/collection';
import { createSubmitter } from '../../../lib/submitter';
import type { JournalNote } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function NewNotesConcern(pg: JournalDetailPageProvider) {
	return {
		...createCollection<JournalNote>(),
		submitter: createSubmitter(),

		async delete(this: any, noteId: string) {
			if (!pg().journal.detail) return;
			const ok = await this.submitter.run(
				() => pg().noteClient.delete(pg().journal.detail!.id, noteId),
				{ success: 'Note deleted' },
			);
			if (ok) this.remove(noteId);
		},

		sorted() {
			return [...this.items].sort((left, right) => {
				const leftTime = left.created_at ? new Date(left.created_at).getTime() : 0;
				const rightTime = right.created_at ? new Date(right.created_at).getTime() : 0;
				return rightTime - leftTime;
			});
		},
	};
}
