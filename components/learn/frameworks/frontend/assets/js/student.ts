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
 *
 * The functions are attached to window at the bottom so they can be bundled into
 * app.js while still being callable from Alpine expressions and inline handlers.
 */

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

// ── studentListPage ───────────────────────────────────────────────────────────
// Page-level Alpine component. Owns the student list, filters, and pagination.
// Opens dialogs via window.tui.dialog.open() and passes data via CustomEvents.
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

    // ── Filters & Pagination ──────────────────────────────────────────────────
    // templUI selectbox fires a native 'change' event; Alpine reads e.target.value.
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

    // ── Open Dialogs ──────────────────────────────────────────────────────────
    // dialog.js teleports Content to <body>, destroying any Alpine scope on it.
    // We use DOM APIs to populate inputs directly before opening the dialog.
    // On submit/confirm, global functions read values back via FormData / getElementById.
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

    // ── Toast ─────────────────────────────────────────────────────────────────
    // Inject minimal [data-tui-toast] node; templUI toast.js handles the rest.
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

// ── Dialog helpers (global, called from onsubmit / onclick in dialog Content) ─
// These must be global because dialog.Content teleports to <body> and loses
// any Alpine scope. Plain JS functions are always accessible regardless of DOM position.
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

// Called by <form onsubmit="return studentFormSubmit(event)"> inside the form dialog.
// Reads FormData, calls the API, then notifies the parent Alpine component via CustomEvent.
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

// Called by the Delete button onclick inside the delete confirmation dialog.
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

document.addEventListener('alpine:init', () => {
  Alpine.data('studentListPage', studentListPage);
});

window.studentListPage = studentListPage;
window.setFormFields = setFormFields;
window.studentFormSubmit = studentFormSubmit;
window.studentDeleteConfirm = studentDeleteConfirm;

export {};
