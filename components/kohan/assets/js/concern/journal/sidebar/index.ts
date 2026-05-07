import { NewSidebarStateConcern } from './state';
import { NewReviewActionsConcern } from './review_actions';
import { NewReviewQueueConcern } from './review_queue';
import { NewNoteFormConcern } from './note_form';
import { NewNotesConcern } from './notes';
import { NewTagCollectionConcern } from './tags';
import { NewReasonTagFormConcern } from './reason_tag_form';
import { NewManagementTagsConcern } from './management_tags';
import type { JournalDetailPageProvider } from '../../../types/journal_detail_concern';

export function NewSidebarConcern(pg: JournalDetailPageProvider) {
	return {
		state: NewSidebarStateConcern(),
		reviewActions: NewReviewActionsConcern(pg),
		reviewQueue: NewReviewQueueConcern(pg),
		noteForm: NewNoteFormConcern(pg),
		notes: NewNotesConcern(pg),
		tags: NewTagCollectionConcern(pg),
		reasonTagForm: NewReasonTagFormConcern(pg),
		managementTags: NewManagementTagsConcern(pg),
	};
}
