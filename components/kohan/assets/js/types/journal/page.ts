import type { JournalClient } from '../../client/journal';
import type { PresentationConcern } from '../core/present';

/** Shared base contract for all journal Alpine page objects. */
export interface JournalPageBase {
	client: JournalClient;
	present: PresentationConcern;
	init(): void;
}

export type PageProvider<T extends JournalPageBase> = () => T;
