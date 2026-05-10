import { JournalStatus } from '../../types/api/journal/enums';
import type { JournalPageProvider } from '../../types/journal/list';
import type { QuickFilterButton } from '../../types/core/quick';
import { BaseQuickButton } from '../../lib/quick_button';

export class QuickStatusButton extends BaseQuickButton<JournalStatus> implements QuickFilterButton<JournalStatus> {
	private pg: JournalPageProvider;

	protected states = [JournalStatus.SET, JournalStatus.RUNNING];
	protected getPresenter = () => this.pg().present.status;

	constructor(pg: JournalPageProvider) {
		super();
		this.pg = pg;
	}

	protected currentValue(): JournalStatus | '' {
		return this.pg().filter.status;
	}

	toggle(): void {
		this.pg().filter.status = this.nextValue;
		this.pg().filter.applyManualFilters();
	}
}
