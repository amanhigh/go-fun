// =============================================================================
// LESSON INDEX
// =============================================================================
// 1. Why this file exists
//    - This is the smallest shared Alpine extension point in the demo.
// 2. Why alpine:init matters
//    - The listener must register before Alpine starts so the hook is reliable.
// 3. Why the shortcut is global
//    - It demonstrates a reusable cross-page behavior without page-specific code.
// 4. Why this file stays tiny
//    - The benchmark is clarity: one file, one focused behavior, one obvious reason.
// =============================================================================

// SECTION 1 — WHY THIS FILE EXISTS
// This file shows how to add a small, intentionally scoped Alpine behavior.
// In a benchmark/demo project, this is the pattern we want future components to copy.

type AlpineCounterElement = HTMLElement & {
  _x_dataStack?: Array<{ count?: number }>;
};

// SECTION 2 — WHY alpine:init IS THE ENTRY POINT
// The listener must be registered before Alpine boots so every page can see it.
document.addEventListener('alpine:init', () => {
  registerCounterResetShortcut();

  console.log('Custom JS initialized');
});

// SECTION 3 — THE SHARED KEYBOARD SHORTCUT
// Pressing R is intentionally generic: it proves how shared DOM behavior can
// be layered on top of Alpine without turning the page component into a utility.
function registerCounterResetShortcut(): void {
  document.addEventListener('keydown', (event: KeyboardEvent) => {
    const target = event.target as HTMLElement | null;
    if (target?.tagName === 'INPUT' || target?.tagName === 'TEXTAREA') return;

    if (event.key.toLowerCase() === 'r') {
      document.querySelectorAll<AlpineCounterElement>('[x-data]').forEach((el) => {
        if (el._x_dataStack?.[0]?.count !== undefined) {
          el._x_dataStack[0].count = 0;
          console.log('Counter reset via R key');
        }
      });
    }
  });
}

// SECTION 4 — WHY THERE IS NO EXTRA EXPORT LOGIC
// This file is imported only for its side effects, so the empty export makes the
// module explicit without introducing accidental API surface.
export {};
