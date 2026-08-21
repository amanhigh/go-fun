import type { PresentationConcern } from '../../types/core/present';
import { NewStatusPresenter } from './status';
import { NewTypePresenter } from './type';
import { NewTimeframePresenter } from './timeframe';
import { NewTagPresenter } from './tag';
import { NewDatePresenter } from './date';
import { NewReviewPresenter } from './review';

export function NewPresentationConcern(): PresentationConcern {
	return {
		status: NewStatusPresenter(),
		type: NewTypePresenter(),
		timeframe: NewTimeframePresenter(),
		tag: NewTagPresenter(),
		date: NewDatePresenter(),
		review: NewReviewPresenter(),
	};
}
