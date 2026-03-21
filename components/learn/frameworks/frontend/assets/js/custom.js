// Custom JavaScript for specific component behaviors

// Counter keyboard shortcuts - only reset, Alpine handles the rest
export function initCounterKeyboardShortcuts() {
    document.addEventListener('keydown', (event) => {
        // Only handle 'r' key when not typing in input fields
        if (event.target.tagName === 'INPUT' || event.target.tagName === 'TEXTAREA') {
            return;
        }
        
        if (event.key.toLowerCase() === 'r') {
            // Reset all counters
            document.querySelectorAll('[x-data*="count"]').forEach(counterEl => {
                if (counterEl._x_dataStack && counterEl._x_dataStack[0].count !== undefined) {
                    counterEl._x_dataStack[0].count = 0;
                }
            });
            console.log('Counter reset via keyboard shortcut');
        }
    });
}

// Custom counter enhancements
export function enhanceCounter() {
    // Add visual feedback for keyboard shortcuts
    const style = document.createElement('style');
    style.textContent = `
        .counter-keyboard-hint {
            position: absolute;
            top: -25px;
            right: 0;
            font-size: 0.75rem;
            color: #6b7280;
            background: rgba(255, 255, 255, 0.9);
            padding: 2px 6px;
            border-radius: 4px;
            opacity: 0;
            transition: opacity 0.2s;
        }
        .counter-container:hover .counter-keyboard-hint {
            opacity: 1;
        }
    `;
    document.head.appendChild(style);
    
    // Add keyboard hints to counter containers
    document.querySelectorAll('[x-data*="count"]').forEach(counterEl => {
        if (!counterEl.querySelector('.counter-keyboard-hint')) {
            const hint = document.createElement('div');
            hint.className = 'counter-keyboard-hint';
            hint.textContent = 'R: Reset';
            counterEl.style.position = 'relative';
            counterEl.appendChild(hint);
        }
    });
    
    console.log('Counter enhancements loaded with keyboard shortcuts');
}

// Main application JavaScript entry point
document.addEventListener('alpine:init', () => {
    // Initialize counter keyboard shortcuts
    initCounterKeyboardShortcuts();
    
    // Initialize counter enhancements
    enhanceCounter();
    
    // Basic showcase Alpine.js component logic
    Alpine.data('basicShowcase', () => ({
        notesLength: 0,
        notesMaxLength: 200,
        showDialog: false,
        
        updateNotesLength(event) {
            this.notesLength = event.target.value.length;
        },
        
        openDialog() {
            this.showDialog = true;
            this.$nextTick(() => {
                this.$refs.dialog.showModal();
            });
        },
        
        closeDialog() {
            this.showDialog = false;
            this.$refs.dialog.close();
        },
        
        submitForm(event) {
            const terms = document.getElementById('terms');
            if (!terms.checked) {
                event.preventDefault();
                alert('Please accept terms and conditions.');
            }
        }
    }))
    
    console.log('Alpine.js and custom counter features initialized');
})
