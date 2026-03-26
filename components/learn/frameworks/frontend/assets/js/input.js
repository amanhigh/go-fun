// Input entry point for esbuild - imports and exports everything
// Page-specific Alpine components (e.g. student.js) are loaded as separate
// static script tags per page — global functions are tree-shaken by esbuild.
import './custom.js';
