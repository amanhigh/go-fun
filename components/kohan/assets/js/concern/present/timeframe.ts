import type { DisplaySpec, DisplayValue, Presenter } from '../../types/present';
import { BasePresenter } from '../../types/present';
import type { JournalTimeframe } from '../../types/journal_api';

const timeframeCatalog: Record<JournalTimeframe, DisplaySpec> = {
	YR: { icon: '🗓️', text: 'YR', class: 'journal-timeframe-yr' },
	SMN: { icon: '📅', text: 'SMN', class: 'journal-timeframe-smn' },
	TMN: { icon: '📈', text: 'TMN', class: 'journal-timeframe-tmn' },
	MN: { icon: '📊', text: 'MN', class: 'journal-timeframe-mn' },
	WK: { icon: '📆', text: 'WK', class: 'journal-timeframe-wk' },
	DL: { icon: '🔍', text: 'DL', class: 'journal-timeframe-dl' },
};

const fallbackSpec: DisplaySpec = { text: '', class: 'journal-timeframe-default' };

class TimeframePresenterImpl extends BasePresenter {
	spec(value: DisplayValue): DisplaySpec {
		const key = (value ?? '').trim().toUpperCase() as JournalTimeframe;
		return timeframeCatalog[key] ?? { text: value ?? '', class: 'journal-timeframe-default' };
	}
}

export function NewTimeframePresenter(): Presenter {
	return new TimeframePresenterImpl();
}
