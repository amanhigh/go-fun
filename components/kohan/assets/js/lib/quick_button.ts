import type { Presenter } from '../types/core/present';

export type QuickButtonResult = {
	label: string;
	className: string;
	nextValue: string;
};

export interface QuickButton {
	resolve(currentValue: string): QuickButtonResult;
}

export type AllSpec = {
	label: string;
	className: string;
};

export abstract class BaseQuickButton implements QuickButton {
	protected abstract states: string[];
	protected abstract getPresenter: () => Presenter;
	protected allSpec: AllSpec = { label: 'All', className: 'journal-display-default' };

	resolve(currentValue: string): QuickButtonResult {
		if (currentValue === '') {
			return this.resolveEmpty();
		}
		if (this.isLastOrUnknown(currentValue)) {
			return this.resolveAll();
		}
		return this.resolveNext(currentValue);
	}

	// Show first state in the cycle.
	private resolveEmpty(): QuickButtonResult {
		return this.presentResult(this.states[0]);
	}

	// Show the "All" spec at end of cycle.
	private resolveAll(): QuickButtonResult {
		return { label: this.allSpec.label, className: this.allSpec.className, nextValue: '' };
	}

	// Show the state immediately after currentValue.
	private resolveNext(currentValue: string): QuickButtonResult {
		const idx = this.states.indexOf(currentValue);
		return this.presentResult(this.states[idx + 1]);
	}

	// Build a result for a real enum value using the Presenter.
	private presentResult(value: string): QuickButtonResult {
		return {
			label: this.getPresenter().label(value),
			className: this.getPresenter().spec(value).class,
			nextValue: value,
		};
	}

	private isLastOrUnknown(currentValue: string): boolean {
		const idx = this.states.indexOf(currentValue);
		return idx === -1 || idx === this.states.length - 1;
	}
}
