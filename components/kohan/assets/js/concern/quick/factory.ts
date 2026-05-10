import type { JournalPageProvider } from '../../types/journal/list';
import type { QuickConcern } from '../../types/core/quick';
import { QuickStatusButton } from './status';

export function NewQuickConcern(pg: JournalPageProvider): QuickConcern {
	return {
		status: new QuickStatusButton(pg),
	};
}
