// DeferredAction provides a minimal, self-contained deferred-action primitive:
// it tracks an active countdown and invokes a single expiry callback unless
// cancelled. It is intentionally decoupled from any specific UI (journal,
// sidebar, etc.) so it can be reused anywhere Alpine state is used.

export type DeferredAction = {
	// active reports whether a deferred action is currently counting down.
	// Alpine can observe this property directly for reactive UI updates.
	active: boolean;
	// remainingSeconds reports the whole seconds left before expiry. Alpine
	// can observe this property directly for reactive UI updates.
	remainingSeconds: number;
	// message holds the human-readable label supplied to start(). Alpine can
	// observe this property directly for reactive UI updates.
	message: string;
	// start begins (or restarts) the countdown with the given message and
	// expiry callback. Any previously scheduled action is cancelled first.
	start(this: DeferredAction, message: string, action: () => void): void;
	// cancel stops the countdown and prevents the expiry callback.
	cancel(this: DeferredAction): void;
};

export function createDeferredAction(durationSeconds = 3): DeferredAction {
	// Timer ID and the pending expiry callback remain private implementation
	// state; only the observable fields above are exposed on the instance.
	let expiryAction: (() => void) | undefined;
	let intervalId: ReturnType<typeof setInterval> | undefined;

	function clearTimer(): void {
		if (intervalId !== undefined) {
			clearInterval(intervalId);
			intervalId = undefined;
		}
	}

	const instance: DeferredAction = {
		active: false,
		remainingSeconds: 0,
		message: '',
		start(this: DeferredAction, startMessage: string, action: () => void): void {
			// Clear any existing timer before (re)starting the lifecycle.
			clearTimer();
			this.message = startMessage;
			this.remainingSeconds = durationSeconds;
			this.active = true;
			expiryAction = action;

			// A single interval drives both the countdown and the expiry
			// callback, keeping the lifecycle in one place. The arrow function
			// captures `this` so it mutates the observable instance fields.
			intervalId = setInterval(() => {
				this.remainingSeconds -= 1;
				if (this.remainingSeconds <= 0) {
					clearTimer();
					// Invoke the expiry callback exactly once.
					this.active = false;
					const cb = expiryAction;
					expiryAction = undefined;
					if (cb) {
						cb();
					}
				}
			}, 1000);
		},
		cancel(this: DeferredAction): void {
			// Cancellation prevents the expiry callback from firing.
			clearTimer();
			this.active = false;
			this.remainingSeconds = 0;
			this.message = '';
			expiryAction = undefined;
		},
	};

	return instance;
}
