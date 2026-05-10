import type { DisplaySpec, DisplayValue, Presenter } from '../../types/present';
import { BasePresenter } from './base';

const typeCatalog: Record<string, DisplaySpec> = {
	TAKEN: { icon: '📈', text: 'TAKEN', class: 'journal-type-taken' },
	REJECTED: { icon: '📉', text: 'REJECTED', class: 'journal-type-rejected' },
};

const fallbackSpec: DisplaySpec = { icon: '🏷️', text: 'Unknown', class: 'journal-display-default' };

class TypePresenterImpl extends BasePresenter {
	protected catalog = typeCatalog;
	protected fallbackSpec = fallbackSpec;
}

export function NewTypePresenter(): Presenter {
	return new TypePresenterImpl();
}
