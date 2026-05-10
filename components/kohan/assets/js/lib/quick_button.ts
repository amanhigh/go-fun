import type { Presenter } from '../types/core/present';

export type QuickButtonResult<TValue extends string> = {
	label: string;
	className: string;
	nextValue: TValue | '';
};

export interface QuickButton<TValue extends string> {
	resolve(currentValue: TValue | ''): QuickButtonResult<TValue>;
}

export type AllSpec = {
	label: string;
	className: string;
};

export class BaseQuickButton<TValue extends string> implements QuickButton<TValue> {
	private states: TValue[];
	private getPresenter: () => Presenter;
	private allSpec: AllSpec;

	constructor(states: TValue[], getPresenter: () => Presenter, allSpec: AllSpec) {
		this.states = states;
		this.getPresenter = getPresenter;
		this.allSpec = allSpec;
	}

	resolve(currentValue: TValue | ''): QuickButtonResult<TValue> {
		if (currentValue === '') {
			return this.resolveEmpty();
		}
		if (this.isLastOrUnknown(currentValue as TValue)) {
			return this.resolveAll();
		}
		return this.resolveNext(currentValue as TValue);
	}

	// Show first state in the cycle.
	private resolveEmpty(): QuickButtonResult<TValue> {
		return this.presentResult(this.states[0]);
	}

	// Show the "All" spec at end of cycle.
	private resolveAll(): QuickButtonResult<TValue> {
		return { label: this.allSpec.label, className: this.allSpec.className, nextValue: '' };
	}

	// Show the state immediately after currentValue.
	private resolveNext(currentValue: TValue): QuickButtonResult<TValue> {
		const idx = this.states.indexOf(currentValue);
		return this.presentResult(this.states[idx + 1]);
	}

	// Build a result for a real enum value using the Presenter.
	private presentResult(value: TValue): QuickButtonResult<TValue> {
		return {
			label: this.getPresenter().label(value),
			className: this.getPresenter().spec(value).class,
			nextValue: value,
		};
	}

	private isLastOrUnknown(currentValue: TValue): boolean {
		const idx = this.states.indexOf(currentValue);
		return idx === -1 || idx === this.states.length - 1;
	}
}
