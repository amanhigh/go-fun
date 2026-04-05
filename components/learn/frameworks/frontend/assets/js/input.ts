// =============================================================================
// LESSON INDEX
// =============================================================================
// 1. Purpose of this bundle entry
//    - This file is the single source of truth for everything that must end up
//      in app.js.
// 2. Why the imports are ordered here
//    - Shared Alpine helpers are loaded first, then page-specific behavior.
// 3. Why this file stays tiny
//    - The benchmark for this demo is clarity: the entry file should explain the
//      pipeline, not duplicate business logic.
// =============================================================================

// SECTION 1 — WHY THIS FILE EXISTS
// This is the esbuild entrypoint used by Makefile generate-js.
// Only modules imported here are bundled into app.js.

import './custom';
import './student';

// SECTION 2 — WHY THERE IS NO EXTRA CODE HERE
// app.js should remain a composition layer. Real behavior belongs in the
// imported modules so each file can be reasoned about independently.

export {};
