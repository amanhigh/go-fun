/**
 * =============================================================================
 * LESSON INDEX
 * =============================================================================
 * 1. Why this page exists
 *    - This file is the benchmark implementation for a real Alpine + templUI
 *      CRUD page in this demo project.
 * 2. Why the data types are explicit
 *    - The TypeScript interfaces define the contract for the page state and for
 *      the student API payloads, making the demo easy to extend safely.
 * 3. Why Alpine owns the live state
 *    - templUI renders the shell, while Alpine owns filtering, pagination, and
 *      the list data that changes without a full page reload.
 * 4. Why the dialogs use DOM helpers
 *    - templUI dialog.Content teleports into <body>, so direct DOM helpers are
 *      more reliable than trying to keep Alpine scope alive across the portal.
 * 5. Why toast feedback is manual
 *    - The benchmark favors explicitness over magic: API success becomes a toast
 *      through one tiny DOM injection that templUI toast.js can pick up.
 * 6. Why we export globals
 *    - The bundle is loaded as a side effect, so the helpers must also be
 *      available on window for Alpine expressions and inline handlers.
 * =============================================================================
 */

// SECTION 1 — WHY THIS PAGE EXISTS
// This is the benchmark page for CRUD + Alpine + templUI. The code is written as
// a teaching example, so the comments explain both the implementation and the why.

// SECTION 2 — WHY THE TYPES ARE EXPLICIT
// These types document the contract between the page model and the backend API.
// Keeping them here makes the demo safer to extend and easier to reason about.

type Grade = '' | 'Freshman' | 'Sophomore' | 'Junior' | 'Senior';
type ToastVariant = 'success' | 'destructive';

interface StudentRecord {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  age: number;
  grade: Grade;
}

interface StudentFormFields {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  age: number;
  grade: Grade;
}

interface StudentListPageState {
  students: StudentRecord[];
  loading: boolean;
  errorMessage: string;
  searchQuery: string;
  selectedGrade: Grade;
  currentPage: number;
  pageSize: number;
  readonly filteredStudents: StudentRecord[];
  readonly totalPages: number;
  readonly paginatedStudents: StudentRecord[];
  readonly startItem: number;
  readonly endItem: number;
  init(): void;
  fetchStudents(this: StudentListPageState): Promise<void>;
  onGradeFilterChange(this: StudentListPageState, event: Event): void;
  clearFilters(this: StudentListPageState): void;
  goToNextPage(this: StudentListPageState): void;
  goToPreviousPage(this: StudentListPageState): void;
  openCreateModal(this: StudentListPageState): void;
  openEditModal(this: StudentListPageState, student: StudentRecord): void;
  openDeleteModal(this: StudentListPageState, student: StudentRecord): void;
  showToast(this: StudentListPageState, message: string, variant?: ToastVariant): void;
  afterSave(this: StudentListPageState, message: string): Promise<void>;
  setError(this: StudentListPageState, message: string): void;
}

// SECTION 3 — WHY WE DECLARE GLOBALS
// The JS bundle is imported for side effects, so Alpine and inline handlers need
// a stable global surface. This avoids hidden coupling and makes the benchmark
// easier to inspect from the browser console.

declare global {
  interface Window {
    tui?: {
      dialog: {
        open(id: string): void;
        close(id: string): void;
        toggle(id: string): void;
        isOpen(id: string): boolean;
      };
    };
    studentListPage: typeof studentListPage;
    setFormFields: typeof setFormFields;
    studentFormSubmit: typeof studentFormSubmit;
    studentDeleteConfirm: typeof studentDeleteConfirm;
  }

  const Alpine: {
    data(name: string, callback: () => StudentListPageState): void;
  };
}

// SECTION 4 — WHY ALPINE OWNS THE PAGE STATE
// templUI provides the structure, but Alpine owns the changing view-model data
// so the page can filter, paginate, and update without a full reload.
function studentListPage(): StudentListPageState {
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
      return this.students.filter((student) => {
        const name = `${student.first_name} ${student.last_name}`.toLowerCase();
        return (!q || name.includes(q) || student.email.toLowerCase().includes(q))
          && (!this.selectedGrade || student.grade === this.selectedGrade);
      });
    },
    get totalPages() { return Math.max(1, Math.ceil(this.filteredStudents.length / this.pageSize)); },
    get paginatedStudents() {
      const start = (this.currentPage - 1) * this.pageSize;
      return this.filteredStudents.slice(start, start + this.pageSize);
    },
    get startItem() { return this.filteredStudents.length ? (this.currentPage - 1) * this.pageSize + 1 : 0; },
    get endItem() { return this.filteredStudents.length ? Math.min(this.currentPage * this.pageSize, this.filteredStudents.length) : 0; },

    init() {
      void this.fetchStudents();
    },

    async fetchStudents() {
      this.loading = true;
      this.errorMessage = '';
      try {
        const response = await fetch('/api/students');
        if (!response.ok) throw new Error('Failed to fetch students');
        const payload = (await response.json()) as { data?: StudentRecord[] };
        this.students = payload.data || [];
        this.currentPage = 1;
      } catch {
        this.errorMessage = "Couldn't load students";
      } finally {
        this.loading = false;
      }
    },

    // SECTION 5 — WHY FILTERS AND PAGINATION LIVE HERE
    // The UI is small, but the behavior is stateful. Keeping it in this object is
    // the cleanest benchmark pattern for a page-level Alpine component.
    onGradeFilterChange(event: Event) {
      const target = event.target as HTMLSelectElement | null;
      this.selectedGrade = (target?.value ?? '') as Grade;
      this.currentPage = 1;
    },
    clearFilters() {
      this.searchQuery = '';
      this.selectedGrade = '';
      this.currentPage = 1;
    },
    goToNextPage() {
      if (this.currentPage < this.totalPages) this.currentPage += 1;
    },
    goToPreviousPage() {
      if (this.currentPage > 1) this.currentPage -= 1;
    },

    // SECTION 6 — WHY DIALOGS USE DOM HELPERS
    // templUI dialog.Content teleports to <body>. That is great for overlay UX,
    // but it means we should not rely on Alpine scope crossing the portal.
    // Direct DOM writes are simpler, explicit, and more reliable for this demo.
    openCreateModal() {
      setFormFields({ id: '', first_name: '', last_name: '', email: '', age: 18, grade: '' });
      window.tui?.dialog.open('student-form-dialog');
    },

    openEditModal(student) {
      setFormFields({
        id: student.id,
        first_name: student.first_name,
        last_name: student.last_name,
        email: student.email,
        age: student.age,
        grade: student.grade,
      });
      window.tui?.dialog.open('student-form-dialog');
    },

    openDeleteModal(student) {
      const deleteId = document.getElementById('s-delete-id') as HTMLInputElement | null;
      const deleteName = document.getElementById('s-delete-name');
      if (deleteId) deleteId.value = student.id;
      if (deleteName) deleteName.textContent = `${student.first_name} ${student.last_name}`;
      window.tui?.dialog.open('student-delete-dialog');
    },

    // SECTION 7 — WHY TOASTS ARE INJECTED MANUALLY
    // A tiny DOM insert is enough to hand off to templUI's toast observer. This
    // keeps the demo benchmark honest: the page owns the state transition, and
    // templUI owns the UI effect.
    showToast(message: string, variant: ToastVariant = 'success') {
      const el = document.getElementById('toast-container');
      if (!el) return;
      el.innerHTML = `<div data-tui-toast data-tui-toast-duration="3000" data-position="top-right" data-variant="${variant}">
        <div class="w-full bg-popover text-popover-foreground rounded-lg shadow-xs border pt-5 pb-4 px-4 flex items-center gap-3">
          <span class="flex-1 text-sm font-semibold">${message}</span>
        </div>
      </div>`;
    },

    async afterSave(message: string) {
      this.showToast(message, 'success');
      await this.fetchStudents();
    },
    setError(message: string) {
      this.errorMessage = message;
    },
  };
}

// SECTION 8 — WHY THE HELPERS ARE GLOBAL
// These helpers are attached to window because dialog Content is teleported and
// because the page intentionally keeps the bundle simple: one file to load, one
// place to inspect, and no hidden dependency injection layer.
function setFormFields(fields: StudentFormFields): void {
  const entries: Array<[keyof StudentFormFields, string | number]> = [
    ['id', fields.id],
    ['first_name', fields.first_name],
    ['last_name', fields.last_name],
    ['email', fields.email],
    ['age', fields.age],
    ['grade', fields.grade],
  ];

  for (const [key, value] of entries) {
    const domKey = `s-${String(key).replace(/_/g, '-')}`;
    const el = document.getElementById(domKey) as HTMLInputElement | HTMLSelectElement | null;
    if (el) {
      el.value = String(value);
    }
  }
}

// SECTION 9 — WHY SUBMIT IS HANDLED THIS WAY
// The form submits through a single helper so the page can stay declarative in
// templUI markup while still performing a real API request and emitting events.
async function studentFormSubmit(event: SubmitEvent): Promise<boolean> {
  event.preventDefault();
  const form = event.target as HTMLFormElement | null;
  if (!form) return false;

  const fd = new FormData(form);
  const id = (document.getElementById('s-student-id') as HTMLInputElement | null)?.value ?? '';
  const isEdit = id !== '';
  const ageValue = fd.get('age');
  const body = {
    first_name: String(fd.get('first_name') ?? ''),
    last_name: String(fd.get('last_name') ?? ''),
    email: String(fd.get('email') ?? ''),
    age: Number(ageValue) || 18,
    grade: String(fd.get('grade') ?? ''),
  };

  try {
    const response = await fetch(isEdit ? `/api/students/${id}` : '/api/students', {
      method: isEdit ? 'PUT' : 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });
    if (!response.ok) throw new Error('Failed to save student');
    window.tui?.dialog.close('student-form-dialog');
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

// SECTION 10 — WHY DELETE USES THE SAME PATTERN
// Delete is intentionally parallel to save: same style, same state handoff, same
// benchmark principle of keeping behavior predictable and easy to audit.
async function studentDeleteConfirm(): Promise<void> {
  const id = (document.getElementById('s-delete-id') as HTMLInputElement | null)?.value ?? '';
  try {
    const response = await fetch(`/api/students/${id}`, { method: 'DELETE' });
    if (!response.ok) throw new Error('Failed to delete student');
    window.tui?.dialog.close('student-delete-dialog');
    window.dispatchEvent(new CustomEvent('student:saved', { detail: { message: 'Student deleted ✓' } }));
  } catch {
    window.dispatchEvent(new CustomEvent('student:error', { detail: { message: 'Failed to delete student' } }));
  }
}

// SECTION 11 — WHY REGISTRATION HAPPENS ON alpine:init
// Alpine must know about the page data factory before it processes x-data.
document.addEventListener('alpine:init', () => {
  Alpine.data('studentListPage', studentListPage);
});

// SECTION 12 — WHY THE GLOBAL EXPORTS REMAIN
// These globals are the compatibility layer for inline handlers and for pages
// that want to inspect the benchmark implementation directly from the console.
window.studentListPage = studentListPage;
window.setFormFields = setFormFields;
window.studentFormSubmit = studentFormSubmit;
window.studentDeleteConfirm = studentDeleteConfirm;

// SECTION 13 — WHY THERE IS NO MODULE EXPORT API
// The file exists for side effects and page wiring, so an empty export keeps the
// module boundary explicit without encouraging direct imports elsewhere.
export {};
