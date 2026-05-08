import { NewSidebarStateConcern } from './state';
import { NewReviewActionsConcern } from './action';
import { NewReviewQueueConcern } from './queue';
import { NewNoteFormConcern } from './note_form';
import { NewNotesConcern } from './notes';
import { NewTagCollectionConcern } from './tags';
import { TagFormConcern } from './tag_form';
import { MgmntTagConcern } from './mgmnt_tag';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewSidebarConcern(pg: JournalDetailPageProvider) {
	return {
		state: NewSidebarStateConcern(),
		reviewActions: NewReviewActionsConcern(pg),
		reviewQueue: NewReviewQueueConcern(pg),
		noteForm: NewNoteFormConcern(pg),
		notes: NewNotesConcern(pg),
		tags: NewTagCollectionConcern(pg),
		reasonTagForm: TagFormConcern(pg),
		managementTags: MgmntTagConcern(pg),
	};
}
