/**
 * Student List Page - Alpine.js Component
 *
 * This file contains all client-side logic for the Student CRUD page.
 * It works alongside student_list.templ which provides the HTML structure
 * using templUI components (table, badge, skeleton, pagination, button, input, form).
 *
 * Architecture:
 *   - templUI components handle all UI rendering and styling
 *   - Alpine.js handles reactive state, API calls, and dynamic behavior
 *   - Modals use Alpine x-show (not templUI dialog) because templUI dialog
 *     moves DOM nodes to document.body, breaking Alpine's x-data scope
 *   - Grade filter uses templUI selectbox (read-only, no 2-way bind needed)
 *   - Form grade uses native <select> styled with Tailwind (needs x-model binding)
 *
 * Data Flow:
 *   Alpine x-data ──> templ components render with x-bind/x-text/x-show
 *   User actions ──> Alpine methods ──> fetch() API ──> update Alpine state
 *   State change ──> Alpine reactivity ──> DOM auto-updates
 */
function studentListPage() {
  return {
    // --- Reactive State ---
    students: [],
    loading: false,
    errorMessage: '',
    searchQuery: '',
    selectedGrade: '',
    currentPage: 1,
    pageSize: 5,

    // Modal state (Alpine-managed since templUI dialog breaks x-data scope)
    showFormModal: false,
    showDeleteModal: false,

    // Form state
    editingStudentID: '',
    formSubmitting: false,
    form: { first_name: '', last_name: '', email: '', age: 18, grade: '' },

    // Delete confirmation state
    deleteStudentID: '',
    deleteStudentName: '',

    // --- Computed Properties ---

    /** Filter students by search query and grade */
    get filteredStudents() {
      const query = this.searchQuery.toLowerCase().trim();
      return this.students.filter((s) => {
        const name = `${s.first_name} ${s.last_name}`.toLowerCase();
        const matchSearch = !query || name.includes(query) || s.email.toLowerCase().includes(query);
        const matchGrade = !this.selectedGrade || s.grade === this.selectedGrade;
        return matchSearch && matchGrade;
      });
    },

    get totalPages() {
      return Math.max(1, Math.ceil(this.filteredStudents.length / this.pageSize));
    },

    get paginatedStudents() {
      const start = (this.currentPage - 1) * this.pageSize;
      return this.filteredStudents.slice(start, start + this.pageSize);
    },

    get startItem() {
      return this.filteredStudents.length === 0 ? 0 : (this.currentPage - 1) * this.pageSize + 1;
    },

    get endItem() {
      return this.filteredStudents.length === 0
        ? 0
        : Math.min(this.currentPage * this.pageSize, this.filteredStudents.length);
    },

    // --- Lifecycle ---

    init() {
      this.fetchStudents();
    },

    // --- API Methods ---

    /** Fetch all students from the API */
    async fetchStudents() {
      this.loading = true;
      this.errorMessage = '';
      try {
        const res = await fetch('/api/students');
        if (!res.ok) throw new Error('Failed');
        const payload = await res.json();
        this.students = payload.data || [];
        this.currentPage = 1;
      } catch {
        this.errorMessage = "Couldn't load students";
      } finally {
        this.loading = false;
      }
    },

    /** Create or update a student via API */
    async submitStudent() {
      this.formSubmitting = true;
      const isUpdate = this.editingStudentID !== '';
      const url = isUpdate ? `/api/students/${this.editingStudentID}` : '/api/students';
      try {
        const res = await fetch(url, {
          method: isUpdate ? 'PUT' : 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(this.form),
        });
        if (!res.ok) throw new Error('Save failed');
        this.showFormModal = false;
        this.showToast(isUpdate ? 'Student updated ✓' : 'Student added ✓');
        await this.fetchStudents();
      } catch {
        this.errorMessage = isUpdate ? 'Failed to update student' : 'Failed to create student';
      } finally {
        this.formSubmitting = false;
      }
    },

    /** Delete a student via API */
    async confirmDelete() {
      if (!this.deleteStudentID) return;
      try {
        const res = await fetch(`/api/students/${this.deleteStudentID}`, { method: 'DELETE' });
        if (!res.ok) throw new Error('Delete failed');
        this.closeDeleteModal();
        this.showToast('Student deleted ✓');
        await this.fetchStudents();
      } catch {
        this.errorMessage = 'Failed to delete student';
      }
    },

    // --- Filter & Pagination ---

    /** Called by templUI selectbox change event for grade filter */
    onGradeFilterChange(event) {
      const input = event.target?.querySelector?.('input[type="hidden"]') || event.target;
      this.selectedGrade = input?.value ?? '';
      this.currentPage = 1;
    },

    clearFilters() {
      this.searchQuery = '';
      this.selectedGrade = '';
      this.currentPage = 1;
    },

    goToNextPage() {
      if (this.currentPage < this.totalPages) this.currentPage++;
    },

    goToPreviousPage() {
      if (this.currentPage > 1) this.currentPage--;
    },

    // --- Modal Management ---
    // Uses Alpine x-show instead of templUI dialog because dialog.Dialog
    // moves DOM nodes to document.body via JS, disconnecting them from
    // Alpine's x-data scope and breaking all reactive bindings.

    openCreateModal() {
      this.editingStudentID = '';
      this.form = { first_name: '', last_name: '', email: '', age: 18, grade: '' };
      this.showFormModal = true;
    },

    openEditModal(student) {
      this.editingStudentID = student.id;
      this.form = {
        first_name: student.first_name,
        last_name: student.last_name,
        email: student.email,
        age: student.age,
        grade: student.grade,
      };
      this.showFormModal = true;
    },

    closeFormModal() {
      if (!this.formSubmitting) this.showFormModal = false;
    },

    openDeleteModal(student) {
      this.deleteStudentID = student.id;
      this.deleteStudentName = `${student.first_name} ${student.last_name}`;
      this.showDeleteModal = true;
    },

    closeDeleteModal() {
      this.showDeleteModal = false;
      this.deleteStudentID = '';
      this.deleteStudentName = '';
    },

    // --- Toast Notification ---
    // Injects a templUI-styled toast into #toast-container.
    // Uses templUI's data-tui-toast attributes so its JS auto-dismisses it.

    showToast(message) {
      const el = document.getElementById('toast-container');
      if (!el) return;
      el.innerHTML = `
        <div data-tui-toast data-tui-toast-duration="3000" data-position="top-right"
          class="z-50 fixed pointer-events-auto p-4 w-full md:max-w-[420px] top-0 right-0
                 animate-in fade-in slide-in-from-top-4 duration-300">
          <div class="w-full bg-popover text-popover-foreground rounded-lg shadow-xs border
                      pt-5 pb-4 px-4 flex items-center gap-3">
            <svg xmlns="http://www.w3.org/2000/svg" class="size-5 text-green-500 flex-shrink-0"
              viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
              <path d="m9 11 3 3L22 4"/>
            </svg>
            <span class="text-sm font-medium">${message}</span>
          </div>
        </div>`;
    },
  };
}
