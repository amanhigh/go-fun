// Input entry point for esbuild - imports and exports everything
import './app.js';

// Re-export anything that needs to be globally available
export { initCustomBehaviors, enhanceCounter } from './custom.js';
