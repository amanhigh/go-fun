import type { DisplaySpec, Presenter } from '../../types/core/present';
import { BasePresenter } from './base';

const imageTypeCatalog: Record<string, DisplaySpec> = {
	SET: { icon: '🎯', text: 'SET', class: 'journal-image-type-set' },
	INFO: { icon: 'ℹ️', text: 'INFO', class: 'journal-image-type-info' },
	RESULT: { icon: '🏁', text: 'RESULT', class: 'journal-image-type-result' },
};

const fallbackSpec: DisplaySpec = { icon: '🏷️', text: 'Unknown', class: 'journal-display-default' };

class ImageTypePresenterImpl extends BasePresenter {
	protected catalog = imageTypeCatalog;
	protected fallbackSpec = fallbackSpec;
}

export function NewImageTypePresenter(): Presenter {
	return new ImageTypePresenterImpl();
}
