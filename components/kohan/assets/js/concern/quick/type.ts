import { JournalType } from '../../types/api/journal/enums';
import type { JournalPageProvider } from '../../types/journal/list';
import type { QuickButtonResult, QuickFilterButton } from '../../types/core/quick';
import { BaseQuickButton } from '../../lib/quick_button';

export class QuickTypeButton extends BaseQuickButton implements QuickFilterButton {
	private pg: JournalPageProvider;

	protected states = [JournalType.TAKEN, JournalType.REJECTED];
	protected getPresenter = () => this.pg().present.type;

	constructor(pg: JournalPageProvider) {
		super();
		this.pg = pg;
	}

	button(): QuickButtonResult {
		return this.resolve(this.pg().filter.type);
	}

	toggle(): void {
		this.pg().filter.type = this.button().nextValue as JournalType;
		this.pg().filter.applyManualFilters();
	}
}
