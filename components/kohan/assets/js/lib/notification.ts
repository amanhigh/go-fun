// Shared notification runtime for the common UI notification viewport.
//
// This file is bundled into Kohan's app.js via esbuild. The notification
// viewport .templ template provides the DOM scaffold; this runtime owns
// queue, timer, and dismissal state privately.
//
// Public API (intentionally small):
//   - notify(notification)  enqueue a notification; returns a dismiss callback
// Callers import notify directly so the runtime dependency remains explicit.
//
// Variant icons and the optional action button are rendered by the Templ
// template as TemplUI scaffolds with static classes. This runtime only selects
// the active variant icon (via data-variant-icon) and clones the action scaffold
// (via data-notification-action-scaffold); it never builds SVG or CSS strings.

import type { Notification } from '../types/notification';

// Success notifications are transient; errors remain until dismissed.
const SUCCESS_NOTIFICATION_DURATION = 3000;

// Private state: one entry per active notification id.
const notificationState = new Map<string, { timer: ReturnType<typeof setTimeout> | undefined; onExpire?: () => void }>();

// Monotonic id source for stacking and lookup.
let notificationSeq = 0;

/**
 * Resolve the shared viewport element rendered by the Notification() template.
 */
function getViewport(): Element | null {
  return document.querySelector("[data-notification-viewport]");
}

/**
 * Build a DOM node for a notification by cloning the template scaffold and
 * populating it from the supplied data. Wires the dismiss control and optional
 * action button, then returns the prepared node (not yet attached).
 */
function buildNotificationNode(notification: Notification, id: string): Element {
  const viewport = getViewport()!;
  const template = viewport.querySelector("[data-notification-template]") as HTMLTemplateElement;
  const node = template.content.firstElementChild!.cloneNode(true) as Element;

  node.setAttribute("data-notification-id", id);

  const variant = notification.variant;

  // Mark the alert with the active variant so the static Tailwind variant
  // classes (declared in the .templ source) apply color. The cloned node is
  // the TemplUI Alert root itself.
  node.setAttribute("data-variant", variant);

  // Reveal the TemplUI icon scaffold for the requested variant; the rest stay
  // hidden. All variant icons are static in the template, so Tailwind can scan
  // their classes.
  const iconSpan = node.querySelector('[data-variant-icon="' + variant + '"]');
  if (iconSpan) {
    iconSpan.removeAttribute("hidden");
  }

  // Title (hidden when absent).
  const titleEl = node.querySelector("[data-notification-title]");
  if (titleEl) {
    if (notification.title) {
      titleEl.textContent = notification.title;
      titleEl.removeAttribute("hidden");
    } else {
      titleEl.setAttribute("hidden", "");
    }
  }

  // Message (always present).
  const messageEl = node.querySelector("[data-notification-message]");
  if (messageEl) {
    messageEl.textContent = notification.message;
  }

  // Optional action button: clone the hidden TemplUI Button scaffold, set the
  // label, and wire a handler that dismisses before invoking the callback.
  const actionSlot = node.querySelector("[data-notification-action]");
  const actionScaffold = actionSlot
    ? actionSlot.querySelector("[data-notification-action-scaffold]")
    : null;
  const action = notification.action;
  if (actionSlot && actionScaffold && action && action.label) {
    const actionButton = actionScaffold.cloneNode(true) as Element;
    actionButton.removeAttribute("hidden");
    actionButton.removeAttribute("data-notification-action-scaffold");
    const labelEl = actionButton.querySelector("[data-notification-action-label]");
    if (labelEl) {
      labelEl.textContent = action.label;
    }
    actionButton.addEventListener("click", function () {
      // Dismiss first, then run the action callback.
      removeNotification(id);
      action.run();
    });
    actionSlot.appendChild(actionButton);
  }

  // Explicit dismiss control.
  const dismissButton = node.querySelector("[data-notification-dismiss]");
  if (dismissButton) {
    dismissButton.addEventListener("click", function () {
      removeNotification(id);
    });
  }

  return node;
}

/**
 * Remove a notification from the viewport and clear its private state.
 * Safe to call multiple times; the timer (if any) is always cleared.
 */
function removeNotification(id: string): void {
  const state = notificationState.get(id);
  if (state && state.timer !== undefined) {
    clearTimeout(state.timer);
  }
  notificationState.delete(id);

  const el = document.querySelector('[data-notification-id="' + id + '"]');
  if (el) {
    el.remove();
  }
}

/**
 * Enqueue a notification into the shared viewport.
 *
 * Behavior:
 *   - Stacks new notifications below existing ones in the viewport.
 *   - Auto-dismisses success notifications after 3000 ms by default; a caller-supplied duration overrides the variant default, while errors persist unless a duration is provided.
 *   - Calls `onExpire` at most once, only when the timer expires (not on explicit dismiss).
 *   - Supports an optional action button and an explicit dismiss control.
 *
 * @returns a dismiss callback that clears the timer and removes the item
 */
export function notify(notification: Notification): () => void {
  const viewport = getViewport();
  if (!viewport || !notification || typeof notification.message !== "string") {
    // Always return the documented DismissNotification contract, even on no-op.
    return function () {};
  }

  notificationSeq += 1;
  const id = "notification-" + notificationSeq;

  const node = buildNotificationNode(notification, id);
  viewport.appendChild(node);

  const duration = notification.duration ?? (notification.variant === "success" ? SUCCESS_NOTIFICATION_DURATION : 0);

  const state: { timer: ReturnType<typeof setTimeout> | undefined; onExpire?: () => void } = { timer: undefined, onExpire: notification.onExpire };
  notificationState.set(id, state);

  if (duration > 0) {
    state.timer = window.setTimeout(function () {
      if (typeof state.onExpire === "function") {
        try {
          state.onExpire();
        } catch (err) {
          console.error("notification onExpire failed", err);
        }
      }
      removeNotification(id);
    }, duration);
  }

  // Return the planned dismiss callback rather than an internal id.
  return function dismiss() {
    removeNotification(id);
  };
}
