// Type-only declaration module for the common notification runtime.
//
// This file documents the notification contract consumed by Kohan TypeScript
// callers. It is NOT a runtime file; the canonical runtime lives in
// lib/notification.ts and is bundled into app.js via esbuild. Importing this
// module with `import type` provides types without resolving the runtime file
// or bundling a second runtime implementation.
//
// Because this module contains only type declarations and ambient global
// declarations (no runtime code), it emits no JavaScript and can be safely
// imported as a type-only module by any caller.

/** Visual variant of a notification. */
export type NotificationVariant = "success" | "error";

/** Optional action button shown on a notification. */
export interface NotificationAction {
  /** Visible action button label. */
  label: string;
  /** Callback invoked when the action is clicked (after dismissal). */
  run: () => void;
}

/** Dismiss callback returned by notify(); clears the timer and removes the item. */
export type DismissNotification = () => void;

/** Notification payload accepted by notify(). */
export interface Notification {
  /** Main notification text (required). */
  message: string;
  /** Optional heading shown above the message. */
  title?: string;
  /** Visual variant and dismissal policy. */
  variant: NotificationVariant;
  /** Optional action button. */
  action?: NotificationAction;
  /** One-time callback fired only when the notification expires via its timer. */
  onExpire?: () => void;
}

/** Enqueue a notification into the shared viewport; returns a dismiss callback. */
export declare function notify(notification: Notification): DismissNotification;

declare global {
  /** Global handle exposed by the Kohan JS bundle runtime. */
  interface Window {
    notify: (notification: Notification) => DismissNotification;
  }
  /** Standalone global notify handle used by callers across components. */
  function notify(notification: Notification): DismissNotification;
}
