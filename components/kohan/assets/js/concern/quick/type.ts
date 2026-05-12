import { JournalType } from '../../types/api/journal/enums';
import type { JournalPageProvider } from '../../types/journal/list';
import { BaseQuickButton } from '../../lib/quick_button';

export class QuickTypeButton extends BaseQuickButton<JournalType> {
	private pg: JournalPageProvider;

	protected states = [JournalType.TAKEN, JournalType.REJECTED];
	protected getPresenter = () => this.pg().present.type;

	protected getFilter = () => this.pg().filter.type;

	protected applyFilter = (value: JournalType | '') => {
		this.pg().filter.type = value;
		this.pg().filter.applyManualFilters();
	};

	constructor(pg: JournalPageProvider) {
		super();
		this.pg = pg;
	}
}
