import type { JournalPageProvider } from '../../types/journal/list';
import type { QuickConcern } from '../../types/core/quick';
import { QuickStatusButton } from './status';
import { QuickTypeButton } from './type';

export function NewQuickConcern(pg: JournalPageProvider): QuickConcern {
	return {
		type: new QuickTypeButton(pg),
		status: new QuickStatusButton(pg),
	};
}
