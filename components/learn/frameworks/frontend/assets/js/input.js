// Input entry point for esbuild - imports and exports everything
// Shared JS and page-specific bundles that expose globals are imported here so
// they end up in app.js and load automatically with the base layout.
import './custom.js';
import './student.js';
