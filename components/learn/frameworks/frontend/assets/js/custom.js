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
    
    // Register Alpine.js components for interactive pages
    Alpine.data('basicShowcase', () => ({
        notesLength: 0,
        notesMaxLength: 200,
        showDialog: false,
        
        updateNotesLength(event) {
            this.notesLength = event.target.value.length;
        },
        
        openDialog() {
            this.showDialog = true;
            this.$nextTick(() => this.$refs.dialog?.showModal());
        },
        
        closeDialog() {
            this.showDialog = false;
            this.$refs.dialog?.close();
        },
        
        submitForm(event) {
            if (!document.getElementById('terms')?.checked) {
                event.preventDefault();
                alert('Please accept terms and conditions.');
            }
        }
    }));
    
    console.log('Custom JS initialized');
});
