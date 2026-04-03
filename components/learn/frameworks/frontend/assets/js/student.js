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
 *     dialog.Content teleports to <body>, destroying Alpine scope. We use plain
 *     DOM APIs (getElementById, FormData) inside dialogs, and CustomEvents to
 *     notify the parent Alpine component after API calls succeed/fail.
 *
 *  5. TOAST ON MUTATION  — After API success, inject a minimal [data-tui-toast]
 *     node. templUI toast.js MutationObserver auto-handles animation + dismiss.
 *
 * ONE Alpine component + TWO global helpers in this file:
 *   studentListPage()      — Alpine: data, filters, pagination, API, opens dialogs
 *   setFormFields()        — populates form inputs via getElementById before dialog open
 *   studentFormSubmit()    — reads FormData on submit, calls API, fires CustomEvent
 *   studentDeleteConfirm() — reads hidden ID, calls DELETE API, fires CustomEvent
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
    // dialog.js teleports Content to <body>, destroying any Alpine scope on it.
    // We use DOM APIs to populate inputs directly before opening the dialog.
    // On submit/confirm, global functions read values back via FormData / getElementById.
    openCreateModal() {
      setFormFields({ id: '', first_name: '', last_name: '', email: '', age: 18, grade: '' });
      window.tui.dialog.open('student-form-dialog');
    },

    openEditModal(s) {
      setFormFields({ id: s.id, first_name: s.first_name, last_name: s.last_name, email: s.email, age: s.age, grade: s.grade });
      window.tui.dialog.open('student-form-dialog');
    },

    openDeleteModal(s) {
      document.getElementById('s-delete-id').value   = s.id;
      document.getElementById('s-delete-name').textContent = `${s.first_name} ${s.last_name}`;
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

    async afterSave(message) { this.showToast(message, 'success'); await this.fetchStudents(); },
    setError(msg)             { this.errorMessage = msg; },
  };
}

// ── Dialog helpers (global, called from onsubmit / onclick in dialog Content) ─
// These must be global because dialog.Content teleports to <body> and loses
// any Alpine scope. Plain JS functions are always accessible regardless of DOM position.

function setFormFields(s) {
  const fields = {
    's-student-id': s.id, 's-first-name': s.first_name, 's-last-name': s.last_name,
    's-email': s.email, 's-age': s.age, 's-grade': s.grade,
  };
  for (const [id, val] of Object.entries(fields)) {
    const el = document.getElementById(id);
    if (el) el.value = val;
    else console.error(`[student] #${id} not found`);
  }
}

// Called by <form onsubmit="return studentFormSubmit(event)"> inside the form dialog.
// Reads FormData, calls the API, then notifies the parent Alpine component via CustomEvent.
async function studentFormSubmit(event) {
  event.preventDefault();
  const fd  = new FormData(event.target);
  const id  = document.getElementById('s-student-id').value;
  const isEdit = id !== '';
  const body = {
    first_name: fd.get('first_name'),
    last_name:  fd.get('last_name'),
    email:      fd.get('email'),
    age:        parseInt(fd.get('age'), 10) || 18,
    grade:      fd.get('grade'),
  };
  try {
    const r = await fetch(isEdit ? `/api/students/${id}` : '/api/students', {
      method:  isEdit ? 'PUT' : 'POST',
      headers: { 'Content-Type': 'application/json' },
      body:    JSON.stringify(body),
    });
    if (!r.ok) throw new Error();
    window.tui.dialog.close('student-form-dialog');
    window.dispatchEvent(new CustomEvent('student:saved', {
      detail: { message: isEdit ? 'Student updated ✓' : 'Student added ✓' },
    }));
  } catch {
    window.dispatchEvent(new CustomEvent('student:error', {
      detail: { message: isEdit ? 'Failed to update' : 'Failed to create' },
    }));
  }
  return false;
}

// Called by the Delete button onclick inside the delete confirmation dialog.
async function studentDeleteConfirm() {
  const id = document.getElementById('s-delete-id').value;
  try {
    const r = await fetch(`/api/students/${id}`, { method: 'DELETE' });
    if (!r.ok) throw new Error();
    window.tui.dialog.close('student-delete-dialog');
    window.dispatchEvent(new CustomEvent('student:saved', { detail: { message: 'Student deleted ✓' } }));
  } catch {
    window.dispatchEvent(new CustomEvent('student:error', { detail: { message: 'Failed to delete student' } }));
  }
}
