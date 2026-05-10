import { createCollection } from '../../../lib/collection';
import { createLoader } from '../../../lib/loader';
import type { JournalNote } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function NewNotesConcern(pg: JournalDetailPageProvider) {
	return {
		...createCollection<JournalNote>(),
		loader: createLoader(),

		async delete(noteId: string) {
			if (!pg().current.journal) return;
			await this.loader.tryRun(
				() => pg().noteClient.delete(pg().current.journalId, noteId),
			);
			this.remove(noteId);
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
