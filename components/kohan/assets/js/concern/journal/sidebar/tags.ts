import { createCollection } from '../../../lib/collection';
import { createSubmitter } from '../../../lib/submitter';
import { JournalTagType } from '../../../types/api/journal/enums';
import type { JournalTag } from '../../../types/api/journal/response';
import type { JournalDetailPageProvider } from '../../../types/journal/detail';

export function NewTagCollectionConcern(pg: JournalDetailPageProvider) {
	return {
		...createCollection<JournalTag>(),
		submitter: createSubmitter(),

		async delete(this: any, tagId: string) {
			if (!pg().journal.detail) return;
			const ok = await this.submitter.run(
				() => pg().tagClient.delete(pg().journal.detail!.id, tagId),
				{ success: 'Tag deleted' },
			);
			if (ok) this.remove(tagId);
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
