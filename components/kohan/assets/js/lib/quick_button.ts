import type { Presenter, DisplaySpec } from '../types/core/present';
import type { QuickFilterButton } from '../types/core/quick';

type QuickButtonView = {
	label: string;
	className: string;
};

export abstract class BaseQuickButton<T extends string> implements QuickFilterButton<T> {
	protected abstract states: readonly T[];
	protected abstract getPresenter: () => Presenter;
	protected abstract getFilter: () => T | '';
	protected abstract applyFilter: (value: T | '') => void;
	private readonly allSpec: DisplaySpec = { text: 'All', class: 'journal-display-default' };

	/** Current label for the button display. */
	get label(): string {
		return this.current().label;
	}

	/** Current class for the button display. */
	get className(): string {
		return this.current().className;
	}

	/** Cycle to the next state: empty → first state, middle → next, last/unknown → clear. */
	toggle(): void {
		const currentValue = this.getFilter();

		if (currentValue === '') {
			this.applyFilter(this.states[0]);
			return;
		}

		const idx = this.states.indexOf(currentValue);

		if (idx === -1 || idx === this.states.length - 1) {
			this.applyFilter('');
			return;
		}

		this.applyFilter(this.states[idx + 1]);
	}

	private current(): QuickButtonView {
		const currentValue = this.getFilter();

		if (currentValue === '') {
			return this.presentView(this.states[0]);
		}

		const idx = this.states.indexOf(currentValue);

		if (idx === -1 || idx === this.states.length - 1) {
			return { label: this.allSpec.text, className: this.allSpec.class };
		}

		return this.presentView(this.states[idx + 1]);
	}

	private presentView(value: T): QuickButtonView {
		return {
			label: this.getPresenter().label(value),
			className: this.getPresenter().spec(value).class,
		};
	}
}
