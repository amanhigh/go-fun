import type { DisplaySpec, DisplayValue, Presenter } from '../../types/present';
import { BasePresenter } from './base';

const timeframeCatalog: Record<string, DisplaySpec> = {
	YR: { icon: '🗓️', text: 'YR', class: 'journal-timeframe-yr' },
	SMN: { icon: '📅', text: 'SMN', class: 'journal-timeframe-smn' },
	TMN: { icon: '📈', text: 'TMN', class: 'journal-timeframe-tmn' },
	MN: { icon: '📊', text: 'MN', class: 'journal-timeframe-mn' },
	WK: { icon: '📆', text: 'WK', class: 'journal-timeframe-wk' },
	DL: { icon: '🔍', text: 'DL', class: 'journal-timeframe-dl' },
};

const fallbackSpec: DisplaySpec = { text: '', class: 'journal-timeframe-default' };

class TimeframePresenterImpl extends BasePresenter {
	protected catalog = timeframeCatalog;
	protected fallbackSpec = fallbackSpec;

	protected fallbackText(value: DisplayValue): string {
		return value ?? '';
	}
}

export function NewTimeframePresenter(): Presenter {
	return new TimeframePresenterImpl();
}
