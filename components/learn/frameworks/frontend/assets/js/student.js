/**
 * Student CRUD Page — Alpine.js component (studentListPage)
 *
 * WHY Alpine.js is required here (and what templUI cannot do alone):
 *
 *  1. CLIENT-SIDE STATE  — templUI renders static HTML. Alpine holds the live
 *     student array, filter values, pagination index, and form data in memory
 *     without a full page reload.
 *
 *  2. ASYNC API CALLS    — fetch() to GET/POST/PUT/DELETE /api/students is JS-only.
 *     templUI has no built-in data-fetching; it pairs with HTMX or Alpine for that.
 *
 *  3. CLIENT FILTERING   — Filtering/searching the in-memory array and slicing
 *     pages is pure JS logic; templUI is a rendering library, not a data layer.
 *
 *  4. MODAL VISIBILITY   — templUI dialog.Dialog teleports its DOM to <body>,
 *     breaking Alpine x-model bindings. Plain x-show on a <div> keeps the modal
 *     inside Alpine's x-data scope so all bindings work correctly.
 *
 *  5. TOAST ON MUTATION  — After an API call succeeds we inject a toast node.
 *     templUI's toast.js watches MutationObserver on <body> and auto-dismisses
 *     any node with [data-tui-toast] — so we inject the minimal wrapper and let
 *     templUI handle the animation and dismiss timer (no custom CSS/JS needed).
 */

// Plain function export — Alpine calls studentListPage() lazily when it processes
// x-data="studentListPage()" on the DOM element, so load order doesn't matter.
// Alpine.data() registration was tried but requires Alpine global to exist at
// alpine:init time, which is unreliable when bundled before Alpine loads.
function studentListPage() {
  return {
    // ── State ─────────────────────────────────────────────────────────────────
    students: [],
    loading: false,
    errorMessage: '',
    searchQuery: '',
    selectedGrade: '',
    currentPage: 1,
    pageSize: 5,

    // Alpine owns modal visibility; templUI dialog cannot be used here (see WHY #4)
    showFormModal: false,
    showDeleteModal: false,

    editingStudentID: '',
    formSubmitting: false,
    form: { first_name: '', last_name: '', email: '', age: 18, grade: '' },
    deleteStudentID: '',
    deleteStudentName: '',

    // ── Computed (Alpine getters drive x-text / x-for / x-bind in the template) ──
    get filteredStudents() {
      const q = this.searchQuery.toLowerCase().trim();
      return this.students.filter((s) => {
        const name = `${s.first_name} ${s.last_name}`.toLowerCase();
        return (!q || name.includes(q) || s.email.toLowerCase().includes(q))
          && (!this.selectedGrade || s.grade === this.selectedGrade);
      });
    },
    get totalPages()      { return Math.max(1, Math.ceil(this.filteredStudents.length / this.pageSize)); },
    get paginatedStudents() { const s = (this.currentPage - 1) * this.pageSize; return this.filteredStudents.slice(s, s + this.pageSize); },
    get startItem()       { return this.filteredStudents.length ? (this.currentPage - 1) * this.pageSize + 1 : 0; },
    get endItem()         { return this.filteredStudents.length ? Math.min(this.currentPage * this.pageSize, this.filteredStudents.length) : 0; },

    // ── Lifecycle ─────────────────────────────────────────────────────────────
    init() { this.fetchStudents(); },

    // ── API ───────────────────────────────────────────────────────────────────
    async fetchStudents() {
      this.loading = true; this.errorMessage = '';
      try {
        const r = await fetch('/api/students');
        if (!r.ok) throw new Error();
        this.students = (await r.json()).data || [];
        this.currentPage = 1;
      } catch { this.errorMessage = "Couldn't load students"; }
      finally  { this.loading = false; }
    },

    async submitStudent() {
      this.formSubmitting = true;
      const update = this.editingStudentID !== '';
      try {
        const r = await fetch(update ? `/api/students/${this.editingStudentID}` : '/api/students', {
          method:  update ? 'PUT' : 'POST',
          headers: { 'Content-Type': 'application/json' },
          body:    JSON.stringify(this.form),
        });
        if (!r.ok) throw new Error();
        this.showFormModal = false;
        this.showToast(update ? 'Student updated ✓' : 'Student added ✓', 'success');
        await this.fetchStudents();
      } catch { this.errorMessage = update ? 'Failed to update' : 'Failed to create'; }
      finally  { this.formSubmitting = false; }
    },

    async confirmDelete() {
      try {
        const r = await fetch(`/api/students/${this.deleteStudentID}`, { method: 'DELETE' });
        if (!r.ok) throw new Error();
        this.closeDeleteModal();
        this.showToast('Student deleted ✓', 'success');
        await this.fetchStudents();
      } catch { this.errorMessage = 'Failed to delete student'; }
    },

    // ── Filters & Pagination ──────────────────────────────────────────────────
    // templUI selectbox fires a native 'change' event on its hidden <input>.
    // Alpine captures it via x-on:change and reads the value from the event target.
    onGradeFilterChange(e) {
      this.selectedGrade = e.target?.value ?? '';
      this.currentPage = 1;
    },

    clearFilters()   { this.searchQuery = ''; this.selectedGrade = ''; this.currentPage = 1; },
    goToNextPage()   { if (this.currentPage < this.totalPages) this.currentPage++; },
    goToPreviousPage() { if (this.currentPage > 1) this.currentPage--; },

    // ── Modals ────────────────────────────────────────────────────────────────
    openCreateModal() { this.editingStudentID = ''; this.form = { first_name: '', last_name: '', email: '', age: 18, grade: '' }; this.showFormModal = true; },
    openEditModal(s)  { this.editingStudentID = s.id; this.form = { first_name: s.first_name, last_name: s.last_name, email: s.email, age: s.age, grade: s.grade }; this.showFormModal = true; },
    closeFormModal()  { if (!this.formSubmitting) this.showFormModal = false; },
    openDeleteModal(s) { this.deleteStudentID = s.id; this.deleteStudentName = `${s.first_name} ${s.last_name}`; this.showDeleteModal = true; },
    closeDeleteModal() { this.showDeleteModal = false; this.deleteStudentID = ''; this.deleteStudentName = ''; },

    // ── Toast ─────────────────────────────────────────────────────────────────
    // We inject the minimal [data-tui-toast] wrapper into #toast-container.
    // templUI's toast.js MutationObserver detects the new node automatically
    // and handles positioning, animation, and auto-dismiss — no custom CSS needed.
    showToast(message, variant = 'success') {
      const el = document.getElementById('toast-container');
      if (!el) return;
      el.innerHTML = `<div data-tui-toast data-tui-toast-duration="3000" data-position="top-right" data-variant="${variant}">
        <div class="w-full bg-popover text-popover-foreground rounded-lg shadow-xs border pt-5 pb-4 px-4 flex items-center gap-3">
          <span class="flex-1 text-sm font-semibold">${message}</span>
        </div>
      </div>`;
    },
  };
}
