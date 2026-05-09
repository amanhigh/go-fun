import type { PresentationConcern } from '../../types/present';
import { NewStatusPresenter } from './status';
import { NewTypePresenter } from './type';
import { NewTimeframePresenter } from './timeframe';

export function NewPresentationConcern(): PresentationConcern {
	return {
		status: NewStatusPresenter(),
		type: NewTypePresenter(),
		timeframe: NewTimeframePresenter(),
	};
}
