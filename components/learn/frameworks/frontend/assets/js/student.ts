/**
 * =============================================================================
 * LESSON INDEX
 * =============================================================================
 * 1. Type declarations
 *    - These define the shape of page state, API payloads, and shared variants.
 * 2. Page state factory
 *    - This builds the Alpine data object that drives filtering, pagination, and
 *      CRUD requests.
 * 3. Dialog and toast helpers
 *    - templUI dialog Content teleports to <body>, so direct DOM helpers keep the
 *      interaction predictable and easy to audit.
 * 4. Alpine registration and globals
 *    - The bundle is loaded for side effects, so we register the factory and
 *      export the helpers on window.
 * =============================================================================
 */

// SECTION 1 — TYPE DECLARATIONS
// These types keep the demo honest: each one documents a contract used by the
// page state, the form dialog, or the student API payloads.

type Grade = '' | 'Freshman' | 'Sophomore' | 'Junior' | 'Senior';
// One-line purpose: the page only allows these grades in filters and forms.
type ToastVariant = 'success' | 'error';
// One-line purpose: page-level toast variants are mapped to templUI variants.

type StudentMutationAction = 'create' | 'update' | 'delete';
// One-line purpose: submit behavior changes based on whether the action adds, edits, or removes a student.

interface StudentErrorResponse {
  error?: string;
}

interface StudentMutationEventDetail {
  message: string;
  action?: StudentMutationAction;
}


interface FetchStudentsOptions {
  page?: number | 'last';
}

interface StudentApiRecord {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  age: number;
  grade: Grade;
}

interface StudentFormValues {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  age: number;
  grade: Grade;
}

interface StudentApiPayload {
  first_name: string;
  last_name: string;
  email: string;
  age: number;
  grade: Grade;
}

interface StudentDataState {
  students: StudentApiRecord[];
  loading: boolean;
  errorMessage: string;
}

interface StudentFilterState {
  searchQuery: string;
  selectedGrade: Grade;
}

interface StudentPaginationState {
  currentPage: number;
  pageSize: number;
}

interface StudentDeleteState {
  pendingDeleteId: string;
  pendingDeleteSeconds: number;
  pendingDeleteTimer: number | null;
}

const studentFormFieldIds = {
  id: 's-student-id',
  firstName: 's-first-name',
  lastName: 's-last-name',
  email: 's-email',
  age: 's-age',
  grade: 's-grade',
} as const;

const studentDialogFieldIds = {
  gradeFilter: 'grade-filter',
} as const;

const emptyStudentFormValues: StudentFormValues = {
  id: '',
  firstName: '',
  lastName: '',
  email: '',
  age: 18,
  grade: '',
};

interface StudentPageState extends StudentDataState, StudentFilterState, StudentPaginationState, StudentDeleteState {
  readonly filteredStudents: StudentApiRecord[];
  readonly totalPages: number;
  readonly paginatedStudents: StudentApiRecord[];
  readonly startItem: number;
  readonly endItem: number;
  init(): void;
  fetchStudents(this: StudentPageState, options?: FetchStudentsOptions): Promise<void>;
  onGradeFilterChange(this: StudentPageState, event: Event): void;
  clearFilters(this: StudentPageState): void;
  goToNextPage(this: StudentPageState): void;
  goToPreviousPage(this: StudentPageState): void;
  openCreateModal(this: StudentPageState): void;
  openEditModal(this: StudentPageState, student: StudentApiRecord): void;
  requestDelete(this: StudentPageState, student: StudentApiRecord): void;
  undoDelete(this: StudentPageState): void;
  confirmPendingDelete(this: StudentPageState): Promise<void>;
  clearPendingDelete(this: StudentPageState): void;
  afterSave(this: StudentPageState, message: string, action?: StudentMutationAction): Promise<void>;
  setError(this: StudentPageState, message: string): void;
}

// SECTION 2 — PAGE STATE FACTORY
// This is the benchmark-level Alpine state object. It owns filtering, paging,
// API calls, and the dialog entry points.

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
    studentPage: typeof studentPage;
    setFormFields: typeof setFormFields;
    studentFormSubmit: typeof studentFormSubmit;
  }

  const Alpine: {
    data(name: string, callback: () => StudentPageState): void;
  };
}

function studentPage(): StudentPageState {
  return {
    students: [],
    loading: false,
    errorMessage: '',
    searchQuery: '',
    selectedGrade: '',
    currentPage: 1,
    pageSize: 4,
    pendingDeleteId: '',
    pendingDeleteSeconds: 0,
    pendingDeleteTimer: null,

    get filteredStudents() {
      const query = this.searchQuery.toLowerCase().trim();
      return this.students.filter((student) => {
        const name = `${student.first_name} ${student.last_name}`.toLowerCase();
        return (!query || name.includes(query) || student.email.toLowerCase().includes(query))
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

    async fetchStudents(options?: FetchStudentsOptions) {
      this.loading = true;
      this.errorMessage = '';
      try {
        const response = await fetch('/api/students');
        if (!response.ok) throw new Error('Failed to fetch students');
        const payload = (await response.json()) as { data?: StudentApiRecord[] };
        this.students = payload.data ?? [];
        const requestedPage = options?.page ?? this.currentPage;
        const targetPage = requestedPage === 'last' ? this.totalPages : requestedPage;
        this.currentPage = Math.min(Math.max(targetPage, 1), this.totalPages);
      } catch {
        this.errorMessage = "Couldn't load students";
      } finally {
        this.loading = false;
      }
    },

    // SECTION 3 — WHY FILTERS AND PAGINATION LIVE HERE
    // The UI is small, but the behavior is stateful. Keeping it in this object is
    // the benchmark pattern for Alpine page state.
    onGradeFilterChange(event: Event) {
      const target = event.target as HTMLInputElement | HTMLSelectElement | null;
      this.selectedGrade = (target?.value ?? '') as Grade;
      this.currentPage = 1;
    },
    clearFilters() {
      this.searchQuery = '';
      this.selectedGrade = '';
      this.currentPage = 1;

      resetTemplSelectboxValue(studentDialogFieldIds.gradeFilter);
    },
    goToNextPage() {
      if (this.currentPage < this.totalPages) this.currentPage += 1;
    },
    goToPreviousPage() {
      if (this.currentPage > 1) this.currentPage -= 1;
    },

    // SECTION 4 — WHY DIALOGS USE DOM HELPERS
    // templUI dialog.Content teleports to <body>, so DOM writes are more reliable
    // than trying to preserve Alpine scope through the portal.
    openCreateModal() {
      setFormFields(emptyStudentFormValues);
      window.tui?.dialog.open('student-form-dialog');
    },

    openEditModal(student) {
      setFormFields(toStudentFormValues(student));
      window.tui?.dialog.open('student-form-dialog');
    },

    requestDelete(student) {
      this.clearPendingDelete();
      this.pendingDeleteId = student.id;
      this.pendingDeleteSeconds = 3;

      this.pendingDeleteTimer = window.setInterval(() => {
        if (this.pendingDeleteSeconds <= 1) {
          void this.confirmPendingDelete();
          return;
        }
        this.pendingDeleteSeconds -= 1;
      }, 1000);
    },

    undoDelete() {
      this.clearPendingDelete();
    },

    async confirmPendingDelete() {
      const id = this.pendingDeleteId;
      if (!id) return;

      this.clearPendingDelete();

      try {
        const response = await fetch(`/api/students/${id}`, { method: 'DELETE' });
        if (!response.ok) throw new Error('Failed to delete student');
        emitStudentEvent('student:saved', 'Student deleted ✓', 'delete');
      } catch {
        emitStudentEvent('student:error', 'Failed to delete student');
      }
    },

    clearPendingDelete() {
      if (this.pendingDeleteTimer !== null) {
        window.clearInterval(this.pendingDeleteTimer);
      }
      this.pendingDeleteId = '';
      this.pendingDeleteSeconds = 0;
      this.pendingDeleteTimer = null;
    },

    // SECTION 5 — WHY THERE IS NO TOAST LOGIC HERE
    // Toasts are rendered by templUI templates; this file only emits business events.

    async afterSave(message: string, action: StudentMutationAction = 'update') {
      await this.fetchStudents({ page: action === 'create' ? 'last' : this.currentPage });
    },
    setError(message: string) {
      this.errorMessage = message;
    },
  };
}

// SECTION 6 — GLOBAL HELPERS FOR DIALOGS
// These helpers stay global because dialog Content is teleported and the bundle
// is intentionally loaded for side effects.
function setFormFields(fields: StudentFormValues): void {
  const mappings: Array<[string, string | number]> = [
    [studentFormFieldIds.id, fields.id],
    [studentFormFieldIds.firstName, fields.firstName],
    [studentFormFieldIds.lastName, fields.lastName],
    [studentFormFieldIds.email, fields.email],
    [studentFormFieldIds.age, fields.age],
    [studentFormFieldIds.grade, fields.grade],
  ];

  for (const [id, value] of mappings) {
    setInputValue(id, value);
  }
}

function toStudentFormValues(student: StudentApiRecord): StudentFormValues {
  return {
    id: student.id,
    firstName: student.first_name,
    lastName: student.last_name,
    email: student.email,
    age: student.age,
    grade: student.grade,
  };
}

function setInputValue(id: string, value: string | number): void {
  const el = document.getElementById(id) as HTMLInputElement | HTMLSelectElement | null;
  if (el) el.value = String(value);
}

function resetTemplSelectboxValue(triggerId: string): void {
  const trigger = document.getElementById(triggerId) as HTMLButtonElement | null;
  const hiddenValue = trigger?.querySelector('input[type="hidden"]') as HTMLInputElement | null;
  if (!hiddenValue) return;

  hiddenValue.value = '';
  hiddenValue.dispatchEvent(new Event('input', { bubbles: true }));
  hiddenValue.dispatchEvent(new Event('change', { bubbles: true }));
}

function readInputValue(id: string): string {
  return (document.getElementById(id) as HTMLInputElement | null)?.value ?? '';
}

function readStudentFormPayload(form: HTMLFormElement): StudentApiPayload {
  const formData = new FormData(form);
  const ageValue = String(formData.get('age') ?? '').trim();
  const ageNumber = Number(ageValue);

  return {
    first_name: String(formData.get('first_name') ?? ''),
    last_name: String(formData.get('last_name') ?? ''),
    email: String(formData.get('email') ?? ''),
    age: Number.isFinite(ageNumber) ? ageNumber : 0,
    grade: String(formData.get('grade') ?? '') as Grade,
  };
}

async function readErrorMessage(response: Response, fallback: string): Promise<string> {
  try {
    const payload = (await response.json()) as StudentErrorResponse;
    return payload.error || fallback;
  } catch {
    return fallback;
  }
}

function emitStudentEvent(
  type: 'student:saved' | 'student:error',
  message: string,
  action?: StudentMutationAction,
): void {
  const detail: StudentMutationEventDetail = action ? { message, action } : { message };
  window.dispatchEvent(new CustomEvent(type, { detail }));
}

// SECTION 7 — WHY SUBMIT IS HANDLED THIS WAY
// One helper keeps the markup declarative while still doing a real API request.
async function studentFormSubmit(event: SubmitEvent): Promise<boolean> {
  event.preventDefault();
  const form = event.target as HTMLFormElement | null;
  if (!form) return false;

  const id = readInputValue(studentFormFieldIds.id);
  const isEdit = id !== '';
  const body = readStudentFormPayload(form);

  try {
    const response = await fetch(isEdit ? `/api/students/${id}` : '/api/students', {
      method: isEdit ? 'PUT' : 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });
    if (!response.ok) {
      const fallback = isEdit ? 'Failed to update student' : 'Failed to create student';
      const errorMessage = await readErrorMessage(response, fallback);
      throw new Error(errorMessage);
    }
    window.tui?.dialog.close('student-form-dialog');
    emitStudentEvent('student:saved', isEdit ? 'Student updated ✓' : 'Student added ✓', isEdit ? 'update' : 'create');
  } catch (error) {
    const fallback = isEdit ? 'Failed to update student' : 'Failed to create student';
    const message = error instanceof Error ? error.message : fallback;
    emitStudentEvent('student:error', message);
  }
  return false;
}

// SECTION 9 — ALPINE REGISTRATION
// Alpine must know about the page factory before it processes x-data.
document.addEventListener('alpine:init', () => {
  Alpine.data('studentPage', studentPage);
});

// SECTION 10 — GLOBAL EXPORTS
// These globals keep inline handlers and browser-console inspection working.
window.studentPage = studentPage;
window.setFormFields = setFormFields;
window.studentFormSubmit = studentFormSubmit;

// SECTION 11 — WHY THERE IS NO MODULE EXPORT API
// This file exists for side effects and page wiring, so the empty export keeps
// the module boundary explicit without encouraging direct imports elsewhere.
export {};
