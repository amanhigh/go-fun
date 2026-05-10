import { JournalStatus } from '../../types/api/journal/enums';
import type { JournalPageProvider } from '../../types/journal/list';
import type { QuickButtonResult, QuickFilterButton } from '../../types/core/quick';
import { BaseQuickButton } from '../../lib/quick_button';

export class QuickStatusButton extends BaseQuickButton implements QuickFilterButton {
	private pg: JournalPageProvider;

	protected states = [JournalStatus.SET, JournalStatus.RUNNING];
	protected getPresenter = () => this.pg().present.status;

	constructor(pg: JournalPageProvider) {
		super();
		this.pg = pg;
	}

	button(): QuickButtonResult {
		return this.resolve(this.pg().filter.status);
	}

	toggle(): void {
		this.pg().filter.status = this.button().nextValue as JournalStatus;
		this.pg().filter.applyManualFilters();
	}
}
