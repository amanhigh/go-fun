// Student page state + helpers for Alpine.js CRUD page.
// Loaded for side-effects: registers Alpine data factory and global form helpers.

type Grade = '' | 'Freshman' | 'Sophomore' | 'Junior' | 'Senior';
type StudentMutationAction = 'create' | 'update' | 'delete';

interface StudentErrorResponse {
  error?: string;
}

interface StudentMutationEventDetail {
  title: string;
  message: string;
  action?: StudentMutationAction;
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

interface StudentListResponse {
  data?: StudentApiRecord[];
  count?: number;
  offset?: number;
  limit?: number;
  total_pages?: number;
}

interface StudentDataState {
  students: StudentApiRecord[];
  loading: boolean;
  errorMessage: string;
  totalStudents: number;
}

interface StudentFilterState {
  name: string;
  selectedGrade: Grade;
}

interface StudentPaginationState {
  page: number;
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
  readonly currentPage: number;
  readonly totalPages: number;
  readonly paginatedStudents: StudentApiRecord[];
  readonly startItem: number;
  readonly endItem: number;
  init(): void;
  fetchStudents(this: StudentPageState): Promise<void>;
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
  showToast(this: StudentPageState, title: string, message: string, isError: boolean): void;
  afterSave(this: StudentPageState, action?: StudentMutationAction): Promise<void>;
  setError(this: StudentPageState, message: string): void;
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
    totalStudents: 0,
    name: '',
    selectedGrade: '',
    page: 1,
    pageSize: 4,
    pendingDeleteId: '',
    pendingDeleteSeconds: 0,
    pendingDeleteTimer: null,

    get currentPage() { return this.page; },
    get filteredStudents() {
      const query = this.name.toLowerCase().trim();
      return this.students.filter((student) => {
        const studentName = `${student.first_name} ${student.last_name}`.toLowerCase();
        return (!query || studentName.includes(query))
          && (!this.selectedGrade || student.grade === this.selectedGrade);
      });
    },
    get totalPages() { return Math.max(1, Math.ceil(this.totalStudents / this.pageSize)); },
    get paginatedStudents() { return this.filteredStudents; },
    get startItem() { return this.totalStudents ? ((this.page - 1) * this.pageSize) + 1 : 0; },
    get endItem() { return this.totalStudents ? Math.min(this.page * this.pageSize, this.totalStudents) : 0; },

    init() {
      window.addEventListener('student:saved', (e) => {
        const detail = (e as CustomEvent<StudentMutationEventDetail>).detail;
        this.showToast(detail.title, detail.message, false);
        void this.afterSave(detail.action);
      });
      window.addEventListener('student:error', (e) => {
        const detail = (e as CustomEvent<StudentMutationEventDetail>).detail;
        this.showToast(detail.title, detail.message, true);
      });
    },

    async fetchStudents() {
      this.loading = true;
      this.errorMessage = '';
      try {
        const offset = (this.page - 1) * this.pageSize;
        const params = new URLSearchParams({
          offset: String(offset),
          limit: String(this.pageSize),
        });
        if (this.name.trim() !== '') params.set('name', this.name.trim());
        if (this.selectedGrade !== '') params.set('grade', this.selectedGrade);

        const response = await fetch(`/api/students?${params.toString()}`);
        if (!response.ok) throw new Error('Failed to fetch students');
        const payload = (await response.json()) as StudentListResponse;
        this.students = payload.data ?? [];
        this.totalStudents = payload.count ?? this.students.length;
        const responsePageSize = payload.limit ?? this.pageSize;
        const responseOffset = payload.offset ?? offset;
        this.pageSize = responsePageSize;
        this.page = Math.floor(responseOffset / responsePageSize) + 1;
      } catch {
        this.errorMessage = "Couldn't load students";
      } finally {
        this.loading = false;
      }
    },

    onGradeFilterChange(event: Event) {
      const target = event.target as HTMLInputElement | HTMLSelectElement | null;
      this.selectedGrade = (target?.value ?? '') as Grade;
      this.page = 1;
      void this.fetchStudents();
    },
    clearFilters() {
      this.name = '';
      this.selectedGrade = '';
      this.page = 1;

      resetTemplSelectboxValue('grade-filter');
      void this.fetchStudents();
    },
    goToNextPage() {
      if (this.currentPage < this.totalPages) {
        this.page += 1;
        void this.fetchStudents();
      }
    },
    goToPreviousPage() {
      if (this.currentPage > 1) {
        this.page -= 1;
        void this.fetchStudents();
      }
    },

    openCreateModal() {
      setFormFields(emptyStudentFormValues);
      window.tui?.dialog.open('student-form-dialog');
    },

    openEditModal(student) {
      setFormFields({
        id: student.id,
        firstName: student.first_name,
        lastName: student.last_name,
        email: student.email,
        age: student.age,
        grade: student.grade,
      });
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
        this.showToast('Student deleted', 'The student record was removed.', false);
        await this.afterSave('delete');
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to delete student';
        this.showToast('Delete failed', message, true);
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

    showToast(title: string, message: string, isError: boolean) {
      const toast = cloneToastTemplate(isError);
      if (!toast) return;

      replaceToastText(toast, title, message);
      document.body.appendChild(toast);
    },

    async afterSave(action: StudentMutationAction = 'update') {
      const nextPage = action === 'create'
        ? Math.max(1, Math.ceil((this.totalStudents + 1) / this.pageSize))
        : this.page;
      this.page = nextPage;
      await this.fetchStudents();
    },
    setError(message: string) {
      this.errorMessage = message;
      this.showToast('Error', message, true);
    },
  };
}

// Global helpers: dialog Content is teleported, so DOM helpers are used for pre-filling forms.
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

function cloneToastTemplate(isError: boolean): HTMLElement | null {
  const templateId = isError ? 'student-error-toast-template' : 'student-success-toast-template';
  const template = document.getElementById(templateId) as HTMLTemplateElement | null;
  if (!template) return null;

  const toast = template.content.firstElementChild?.cloneNode(true);
  return toast instanceof HTMLElement ? toast : null;
}

function replaceToastText(toast: HTMLElement, title: string, message: string): void {
  const walker = document.createTreeWalker(toast, NodeFilter.SHOW_TEXT);
  let currentNode = walker.nextNode();

  for (; currentNode; currentNode = walker.nextNode()) {
    const text = currentNode.textContent;
    if (!text) continue;

    currentNode.textContent = text
      .replaceAll('__STUDENT_TOAST_TITLE__', title)
      .replaceAll('__STUDENT_TOAST_MESSAGE__', message);
  }
}

function emitStudentEvent(
  type: 'student:saved' | 'student:error',
  title: string,
  message: string,
  action?: StudentMutationAction,
): void {
  const detail: StudentMutationEventDetail = action ? { title, message, action } : { title, message };
  window.dispatchEvent(new CustomEvent(type, { detail }));
}

// studentFormSubmit: reads FormData, calls API, dispatches student:saved/student:error event.
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
    emitStudentEvent(
      'student:saved',
      isEdit ? 'Student updated' : 'Student added',
      isEdit ? 'The student details were saved successfully.' : 'The student record was created successfully.',
      isEdit ? 'update' : 'create',
    );
  } catch (error) {
    const fallback = isEdit ? 'Failed to update student' : 'Failed to create student';
    const message = error instanceof Error ? error.message : fallback;
    emitStudentEvent('student:error', isEdit ? 'Update failed' : 'Create failed', message);
  }
  return false;
}

document.addEventListener('alpine:init', () => {
  Alpine.data('studentPage', studentPage);
});

window.studentPage = studentPage;
window.setFormFields = setFormFields;
window.studentFormSubmit = studentFormSubmit;

export {};
