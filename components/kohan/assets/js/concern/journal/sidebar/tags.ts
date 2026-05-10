import { createCollection } from '../../../lib/collection';
import { createLoader } from '../../../lib/loader';
import { JournalTagType } from '../../../types/api/journal/enums';
import type { JournalTag } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function NewTagCollectionConcern(pg: JournalDetailPageProvider) {
	return {
		...createCollection<JournalTag>(),
		loader: createLoader(),

		async delete(tagId: string) {
			if (!pg().journal.detail) return;
			await this.loader.tryRun(
				() => pg().tagClient.delete(pg().journal.detail!.id, tagId),
			);
			this.remove(tagId);
		},

		reason() {
			return this.all().filter((tag) => tag.type === JournalTagType.REASON);
		},
		directional() {
			return this.all().filter((tag) => tag.type === JournalTagType.DIRECTION);
		},
		management() {
			return this.all().filter((tag) => tag.type === JournalTagType.MANAGEMENT);
		},
	};
}
