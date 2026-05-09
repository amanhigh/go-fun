import type { PresentationConcern } from '../../types/present';
import { NewStatusPresenter } from './status';
import { NewTypePresenter } from './type';
import { NewTimeframePresenter } from './timeframe';
import { NewTagPresenter } from './tag';
import { NewSequencePresenter } from './sequence';

export function NewPresentationConcern(): PresentationConcern {
	return {
		status: NewStatusPresenter(),
		type: NewTypePresenter(),
		timeframe: NewTimeframePresenter(),
		tag: NewTagPresenter(),
		sequence: NewSequencePresenter(),
	};
}
