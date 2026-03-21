// Basic showcase Alpine.js component logic
document.addEventListener('alpine:init', () => {
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
})
