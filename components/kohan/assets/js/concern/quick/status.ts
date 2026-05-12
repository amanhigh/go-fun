import { JournalStatus } from '../../types/api/journal/enums';
import type { JournalPageProvider } from '../../types/journal/list';
import { BaseQuickButton } from '../../lib/quick_button';

export class QuickStatusButton extends BaseQuickButton<JournalStatus> {
	private pg: JournalPageProvider;

	protected states = [JournalStatus.SET, JournalStatus.RUNNING];
	protected getPresenter = () => this.pg().present.status;

	protected getFilter = () => this.pg().filter.status;

	protected applyFilter = (value: JournalStatus | '') => {
		this.pg().filter.status = value;
		this.pg().filter.applyManualFilters();
	};

	constructor(pg: JournalPageProvider) {
		super();
		this.pg = pg;
	}
}
