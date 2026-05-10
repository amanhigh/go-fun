import type { DisplaySpec, DisplayValue, Presenter } from '../../types/present';

export abstract class BasePresenter implements Presenter<DisplayValue> {
	protected abstract catalog: Record<string, DisplaySpec>;
	protected abstract fallbackSpec: DisplaySpec;

	protected fallbackText(_value: DisplayValue, key: string): string {
		return key || this.fallbackSpec.text;
	}

	spec(value: DisplayValue): DisplaySpec {
		const key = value ?? '';
		return this.catalog[key] ?? { ...this.fallbackSpec, text: this.fallbackText(value, key) };
	}

	label(value: DisplayValue): string {
		const s = this.spec(value);
		return s.icon ? `${s.icon} ${s.text}` : s.text;
	}
}
