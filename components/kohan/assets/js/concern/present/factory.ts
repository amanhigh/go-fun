import type { PresentationConcern } from '../../types/present';
import { NewStatusPresenter } from './status';

export function NewPresentationConcern(): PresentationConcern {
	return {
		status: NewStatusPresenter(),
	};
}
