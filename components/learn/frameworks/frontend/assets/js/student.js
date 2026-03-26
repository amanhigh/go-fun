/**
 * Student CRUD Page — Alpine.js components
 *
 * WHY Alpine.js is required here (and what templUI cannot do alone):
 *
 *  1. CLIENT-SIDE STATE  — templUI renders static HTML. Alpine holds the live
 *     student array, filter values, and pagination in memory without page reloads.
 *
 *  2. ASYNC API CALLS    — fetch() GET/POST/PUT/DELETE is JS-only. templUI is a
 *     rendering library; it pairs with HTMX or Alpine for data fetching.
 *
 *  3. CLIENT FILTERING   — Filtering, searching, and paginating an in-memory
 *     array is pure JS logic that templUI has no equivalent for.
 *
 *  4. DIALOG COORDINATION — templUI dialog.Dialog IS used (window.tui.dialog API).
 *     dialog.Content teleports to <body>, so each dialog carries its own x-data
 *     scope. The parent communicates via window CustomEvents before opening.
 *
 *  5. TOAST ON MUTATION  — After API success, inject a minimal [data-tui-toast]
 *     node. templUI toast.js MutationObserver auto-handles animation + dismiss.
 *
 * THREE Alpine components in this file:
 *   studentListPage()   — page-level: data, filters, pagination, API, opens dialogs
 *   studentFormDialog() — scoped to form dialog Content (teleports to <body>)
 *   studentDeleteDialog() — scoped to delete dialog Content (teleports to <body>)
 */

// ── studentListPage ───────────────────────────────────────────────────────────
// Page-level Alpine component. Owns the student list, filters, and pagination.
// Opens dialogs via window.tui.dialog.open() and passes data via CustomEvents.
function studentListPage() {
  return {
    students: [],
    loading: false,
    errorMessage: '',
    searchQuery: '',
    selectedGrade: '',
    currentPage: 1,
    pageSize: 5,

    get filteredStudents() {
      const q = this.searchQuery.toLowerCase().trim();
      return this.students.filter((s) => {
        const name = `${s.first_name} ${s.last_name}`.toLowerCase();
        return (!q || name.includes(q) || s.email.toLowerCase().includes(q))
          && (!this.selectedGrade || s.grade === this.selectedGrade);
      });
    },
    get totalPages()        { return Math.max(1, Math.ceil(this.filteredStudents.length / this.pageSize)); },
    get paginatedStudents() { const s = (this.currentPage - 1) * this.pageSize; return this.filteredStudents.slice(s, s + this.pageSize); },
    get startItem()         { return this.filteredStudents.length ? (this.currentPage - 1) * this.pageSize + 1 : 0; },
    get endItem()           { return this.filteredStudents.length ? Math.min(this.currentPage * this.pageSize, this.filteredStudents.length) : 0; },

    init() { this.fetchStudents(); },

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

    // ── Filters & Pagination ──────────────────────────────────────────────────
    // templUI selectbox fires a native 'change' event; Alpine reads e.target.value.
    onGradeFilterChange(e) { this.selectedGrade = e.target?.value ?? ''; this.currentPage = 1; },
    clearFilters()          { this.searchQuery = ''; this.selectedGrade = ''; this.currentPage = 1; },
    goToNextPage()          { if (this.currentPage < this.totalPages) this.currentPage++; },
    goToPreviousPage()      { if (this.currentPage > 1) this.currentPage--; },

    // ── Open Dialogs ──────────────────────────────────────────────────────────
    // Dispatch CustomEvent BEFORE opening so dialog x-data has the payload ready.
    // window.tui.dialog.open/close is the public API exposed by templUI dialog.js.
    openCreateModal() {
      window.dispatchEvent(new CustomEvent('student:open-form', {
        detail: { id: '', form: { first_name: '', last_name: '', email: '', age: 18, grade: '' } },
      }));
      window.tui.dialog.open('student-form-dialog');
    },

    openEditModal(s) {
      window.dispatchEvent(new CustomEvent('student:open-form', {
        detail: { id: s.id, form: { first_name: s.first_name, last_name: s.last_name, email: s.email, age: s.age, grade: s.grade } },
      }));
      window.tui.dialog.open('student-form-dialog');
    },

    openDeleteModal(s) {
      window.dispatchEvent(new CustomEvent('student:open-delete', {
        detail: { id: s.id, name: `${s.first_name} ${s.last_name}` },
      }));
      window.tui.dialog.open('student-delete-dialog');
    },

    // ── Toast ─────────────────────────────────────────────────────────────────
    // Inject minimal [data-tui-toast] node; templUI toast.js handles the rest.
    showToast(message, variant = 'success') {
      const el = document.getElementById('toast-container');
      if (!el) return;
      el.innerHTML = `<div data-tui-toast data-tui-toast-duration="3000" data-position="top-right" data-variant="${variant}">
        <div class="w-full bg-popover text-popover-foreground rounded-lg shadow-xs border pt-5 pb-4 px-4 flex items-center gap-3">
          <span class="flex-1 text-sm font-semibold">${message}</span>
        </div>
      </div>`;
    },

    // Called by dialog components after API success to refresh the list + show toast.
    async afterSave(message) {
      this.showToast(message, 'success');
      await this.fetchStudents();
    },

    setError(msg) { this.errorMessage = msg; },
  };
}

// ── studentFormDialog ─────────────────────────────────────────────────────────
// Scoped Alpine component for the form dialog Content node.
// This x-data travels with the Content node when dialog.js teleports it to <body>.
// Receives form payload via 'student:open-form' CustomEvent from parent.
function studentFormDialog() {
  return {
    isEdit: false,
    studentID: '',
    submitting: false,
    form: { first_name: '', last_name: '', email: '', age: 18, grade: '' },

    // Called by x-on:student:open-form.window when parent dispatches before open
    receive(detail) {
      this.studentID = detail.id;
      this.isEdit    = detail.id !== '';
      this.form      = { ...detail.form };
    },

    async submit() {
      this.submitting = true;
      try {
        const r = await fetch(this.isEdit ? `/api/students/${this.studentID}` : '/api/students', {
          method:  this.isEdit ? 'PUT' : 'POST',
          headers: { 'Content-Type': 'application/json' },
          body:    JSON.stringify(this.form),
        });
        if (!r.ok) throw new Error();
        window.tui.dialog.close('student-form-dialog');
        // Notify parent to refresh list and show toast
        window.dispatchEvent(new CustomEvent('student:saved', {
          detail: { message: this.isEdit ? 'Student updated ✓' : 'Student added ✓' },
        }));
      } catch {
        window.dispatchEvent(new CustomEvent('student:error', {
          detail: { message: this.isEdit ? 'Failed to update' : 'Failed to create' },
        }));
      } finally { this.submitting = false; }
    },
  };
}

// ── studentDeleteDialog ───────────────────────────────────────────────────────
// Scoped Alpine component for the delete confirmation dialog Content node.
// Receives {id, name} via 'student:open-delete' CustomEvent from parent.
function studentDeleteDialog() {
  return {
    studentID: '',
    name: '',

    receive(detail) { this.studentID = detail.id; this.name = detail.name; },

    async confirmDelete() {
      try {
        const r = await fetch(`/api/students/${this.studentID}`, { method: 'DELETE' });
        if (!r.ok) throw new Error();
        window.tui.dialog.close('student-delete-dialog');
        window.dispatchEvent(new CustomEvent('student:saved', {
          detail: { message: 'Student deleted ✓' },
        }));
      } catch {
        window.dispatchEvent(new CustomEvent('student:error', {
          detail: { message: 'Failed to delete student' },
        }));
      }
    },
  };
}
