// Custom JavaScript for component behaviors
// This file demonstrates how to extend Alpine.js with custom functionality

// Register alpine:init listener BEFORE Alpine.js loads
document.addEventListener('alpine:init', () => {
    
    // Custom keyboard shortcut: Press 'R' to reset all counters
    document.addEventListener('keydown', (event) => {
        if (event.target.tagName === 'INPUT' || event.target.tagName === 'TEXTAREA') return;
        
        if (event.key.toLowerCase() === 'r') {
            document.querySelectorAll('[x-data]').forEach(el => {
                if (el._x_dataStack?.[0]?.count !== undefined) {
                    el._x_dataStack[0].count = 0;
                    console.log('Counter reset via R key');
                }
            });
        }
    });
    
    console.log('Custom JS initialized');
});
