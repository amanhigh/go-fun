import type { Presenter } from '../types/core/present';
import type { QuickFilterButton } from '../types/core/quick';

type QuickButtonState<T extends string> = {
	label: string;
	className: string;
	nextValue: T;
};

export type AllSpec = {
	label: string;
	className: string;
};

export abstract class BaseQuickButton<T extends string> implements QuickFilterButton<T> {
	protected abstract states: T[];
	protected abstract getPresenter: () => Presenter;
	protected abstract currentValue(): T | '';
	protected allSpec: AllSpec = { label: 'All', className: 'journal-display-default' };

	/** Current label for the button display. */
	get label(): string {
		return this.current().label;
	}

	/** Current class for the button display. */
	get className(): string {
		return this.current().className;
	}

	/** The state value the toggle will cycle to next. */
	get nextValue(): T {
		return this.current().nextValue;
	}

	resolve(currentValue: T | ''): QuickButtonState<T> {
		if (currentValue === '') {
			return this.resolveEmpty();
		}
		if (this.isLastOrUnknown(currentValue)) {
			return this.resolveAll();
		}
		return this.resolveNext(currentValue);
	}

	abstract toggle(): void;

	private current(): QuickButtonState<T> {
		return this.resolve(this.currentValue());
	}

	// Show first state in the cycle.
	private resolveEmpty(): QuickButtonState<T> {
		return this.presentResult(this.states[0]);
	}

	// Show the "All" spec at end of cycle; nextValue wraps to first state.
	private resolveAll(): QuickButtonState<T> {
		return { label: this.allSpec.label, className: this.allSpec.className, nextValue: this.states[0] };
	}

	// Show the state immediately after currentValue.
	private resolveNext(currentValue: T): QuickButtonState<T> {
		const idx = this.states.indexOf(currentValue);
		return this.presentResult(this.states[idx + 1]);
	}

	// Build a result for a real enum value using the Presenter.
	private presentResult(value: T): QuickButtonState<T> {
		return {
			label: this.getPresenter().label(value),
			className: this.getPresenter().spec(value).class,
			nextValue: value,
		};
	}

	private isLastOrUnknown(currentValue: T): boolean {
		const idx = this.states.indexOf(currentValue);
		return idx === -1 || idx === this.states.length - 1;
	}
}
