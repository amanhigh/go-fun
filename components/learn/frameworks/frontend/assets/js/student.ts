// =============================================================================
// SECTION 1: Type Definitions
// =============================================================================
// Core types for Student CRUD operations - API requests, responses, and form data.

type Grade = '' | 'Freshman' | 'Sophomore' | 'Junior' | 'Senior';
type StudentMutationAction = 'create' | 'update' | 'delete';

// Student Response - student data from API
interface StudentResponse {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  age: number;
  grade: Grade;
}

// Student Form Data - Alpine x-model binding
interface StudentFormData {
  firstName: string;
  lastName: string;
  email: string;
  age: number;
  grade: Grade;
}

// Student Request - request body for create/update
interface StudentRequest {
  first_name: string;
  last_name: string;
  email: string;
  age: number;
  grade: Grade;
}

// Student List Response - paginated API response
interface StudentListResponse {
  data?: StudentResponse[];
  count?: number;
  offset?: number;
  limit?: number;
  total_pages?: number;
}

// Empty form values for reset
const emptyFormValues: StudentFormData = {
  firstName: '',
  lastName: '',
  email: '',
  age: 18,
  grade: '' as Grade,
};

// =============================================================================
// SECTION 2: Pagination Tracker
// =============================================================================

interface PaginationTracker {
  page: number;
  pageSize: number;
  totalStudents: number;
  getPage(): number;
  getPageSize(): number;
  getTotalStudents(): number;
  getTotalPages(): number;
  hasNext(): boolean;
  hasPrev(): boolean;
  setTotalStudents(count: number): void;
  setPageFromResponse(offset: number): void;
  nextPage(): void;
  prevPage(): void;
  resetPage(): void;
}

function createPaginationTracker(pageSize: number): PaginationTracker {
  return {
    page: 1,
    pageSize,
    totalStudents: 0,
    getPage() { return this.page; },
    getPageSize() { return this.pageSize; },
    getTotalStudents() { return this.totalStudents; },
    getTotalPages() { return Math.max(1, Math.ceil(this.totalStudents / this.pageSize)); },
    hasNext() { return this.page < this.getTotalPages(); },
    hasPrev() { return this.page > 1; },
    setTotalStudents(count: number) { this.totalStudents = count; },
    setPageFromResponse(offset: number) { this.page = Math.floor(offset / this.pageSize) + 1; },
    nextPage() { if (this.hasNext()) this.page += 1; },
    prevPage() { if (this.hasPrev()) this.page -= 1; },
    resetPage() { this.page = 1; },
  };
}

// =============================================================================
// SECTION 3: Delete Tracker
// =============================================================================

interface DeleteTracker {
  pendingDeleteId: string;
  pendingDeleteSeconds: number;
  pendingDeleteTimer: number | null;
  deletingId: string;
  isPendingDelete(studentId: string): boolean;
  isDeletingId(studentId: string): boolean;
  getPendingDeleteSeconds(): number;
  startCountdown(studentId: string, onExpire: () => void): void;
  cancelCountdown(): void;
  clearAll(): void;
}

function createDeleteTracker(): DeleteTracker {
  return {
    pendingDeleteId: '',
    pendingDeleteSeconds: 0,
    pendingDeleteTimer: null,
    deletingId: '',
    isPendingDelete(studentId: string) { return this.pendingDeleteId === studentId; },
    isDeletingId(studentId: string) { return this.deletingId === studentId; },
    getPendingDeleteSeconds() { return this.pendingDeleteSeconds; },
    startCountdown(studentId: string, onExpire: () => void) {
      this.pendingDeleteId = studentId;
      this.pendingDeleteSeconds = 3;
      this.pendingDeleteTimer = window.setInterval(() => {
        if (this.pendingDeleteSeconds <= 1) {
          this.cancelCountdown();
          onExpire();
          return;
        }
        this.pendingDeleteSeconds -= 1;
      }, 1000);
    },
    cancelCountdown() {
      if (this.pendingDeleteTimer !== null) {
        window.clearInterval(this.pendingDeleteTimer);
      }
      this.pendingDeleteId = '';
      this.pendingDeleteSeconds = 0;
      this.pendingDeleteTimer = null;
    },
    clearAll() {
      this.cancelCountdown();
      this.deletingId = '';
    },
  };
}

// =============================================================================
// SECTION 4: Filter Tracker
// =============================================================================

interface FilterTracker {
  name: string;
  grade: string;
  getName(): string;
  getGrade(): string;
  hasFilters(): boolean;
  setName(name: string): void;
  setGrade(grade: string): void;
  clear(): void;
}

function createFilterTracker(): FilterTracker {
  return {
    name: '',
    grade: '' as Grade,
    getName() { return this.name.trim(); },
    getGrade() { return this.grade; },
    hasFilters() { return this.name !== '' || this.grade !== ''; },
    setName(name: string) { this.name = name.trim(); },
    setGrade(grade: string) { this.grade = grade; },
    clear() { this.name = ''; this.grade = '' as Grade; },
  };
}

// =============================================================================
// SECTION 5: Global Declarations
// =============================================================================

declare global {
  interface Window {
    tui?: {
      dialog: {
        open(id: string): void;
        close(id: string): void;
      };
    };
  }

  const Alpine: {
    data(name: string, callback: () => ReturnType<typeof studentPage>): void;
  };
}

// =============================================================================
// SECTION 3: Alpine Factory
// =============================================================================
// Factory function to create Alpine.js component with all state and methods.

function studentPage() {
  return {
    // ─────────────────────────────────────────────────────────────
    // Data State
    // ─────────────────────────────────────────────────────────────
    students: [] as StudentResponse[],

    // ─────────────────────────────────────────────────────────────
    // Filter State (delegated to FilterTracker)
    // ─────────────────────────────────────────────────────────────
    filterTracker: createFilterTracker(),

// ─────────────────────────────────────────────────────────────
// Pagination State (delegated to PaginationTracker)
// ─────────────────────────────────────────────────────────────
    pagination: createPaginationTracker(4),

    // ─────────────────────────────────────────────────────────────
    // Delete State (delegated to DeleteTracker)
    // ─────────────────────────────────────────────────────────────
    deleteTracker: createDeleteTracker(),

    // ─────────────────────────────────────────────────────────────
    // Form State
    // ─────────────────────────────────────────────────────────────
    editingId: '',
    form: { ...emptyFormValues },

    // ═══════════════════════════════════════════════════════════════
    // SECTION: UI States
    // ═══════════════════════════════════════════════════════════════
    //
    // These states represent mutually exclusive UI conditions.
    // They control what the user sees and interacts with.
    //
    // State Transitions (varies based on user actions):
    // ┌────────────────┐     ┌─────────────┐     ┌─────────────┐
    // │  ✱Initialized  │────▶│   Loading   │────▶│   Loaded    │
    // │  (not loaded)  │     │ (fetching)  │     │  (success)  │
    // └────────────────┘     └─────────────┘     └─────────────┘
    //                              │                    │
    //                              ▼                    ▼
    //                        ┌───────────┐      ┌──────────────┐
    //                        │  Errored  │      │   Editing    │
    //                        │  (failed) │      │  (dialog open)│
    //                        └───────────┘      └──────────────┘
    //                              │                    │
    //                              ▼                    ▼
    //                        ┌───────────┐      ┌──────────────┐
    //                        │  Loading  │      │   Saving     │
    //                        │  (retry)  │      │ (submitting) │
    //                        └───────────┘      └──────────────┘
    //                                                     │
    //                                                     ▼
    //                                              ┌──────────────┐
    //                                              │   Loaded     │
    //                                              │  (success)   │
    //                                              └──────────────┘
    //
    // States:
    //   - initialized: Page has not made initial fetch yet
    //   - loading: Currently fetching student list (pagination/filter)
    //   - editingId: Edit mode uses ID, create mode has empty ID
    //   - saving: Form is being submitted
    //   - deleteTracker: Row delete state (pending countdown, API delete)
    //   - apiError: Error state for API/operation errors
    //   - formError: Inline form validation errors
    //
    // Mutually Exclusive Groups:
    //   1. Load States: initialized ↔ loading ↔ loaded (via initialized)
    //   2. Edit States: idle (no editingId) ↔ editing (has editingId) ↔ saving
    //   3. Error States: hasError() ↔ isFormError() (can coexist)
    // ═══════════════════════════════════════════════════════════════
    initialized: false,
    loading: false,
    saving: false,
    editLoading: false,
    apiError: '',
    formError: '',

    // ─────────────────────────────────────────────────────────────
    // Computed Properties
    // ─────────────────────────────────────────────────────────────
    get hasFilters() { return this.filterTracker.hasFilters(); },
    get filteredStudents() {
      const query = this.filterTracker.getName().toLowerCase().trim();
      return this.students.filter((s) => {
        const fullName = `${s.first_name} ${s.last_name}`.toLowerCase();
        return (!query || fullName.includes(query)) && (!this.filterTracker.getGrade() || s.grade === this.filterTracker.getGrade());
      });
    },
    get paginatedStudents() { return this.filteredStudents; },
    isInitialized() { return this.initialized; },
    isEditing() { return this.editingId !== ''; },
    isLoading() { return this.loading; },
    isSaving() { return this.saving; },
    isLoadingEdit() { return this.editLoading; },
    isEmpty() { return this.filteredStudents.length === 0; },
    isApiError() { return this.apiError !== ''; },
    isFormError() { return this.formError !== ''; },
    isReady() { return !this.initialized && this.loading; },

    // ─────────────────────────────────────────────────────────────
    // Methods - List
    // ─────────────────────────────────────────────────────────────
    async listStudents() {
      this.loading = true;
      this.apiError = '';
      try {
        const page = this.pagination.getPage();
        const pageSize = this.pagination.getPageSize();
        const params = new URLSearchParams({
          offset: String((page - 1) * pageSize),
          limit: String(pageSize),
        });
        const name = this.filterTracker.getName();
        const grade = this.filterTracker.getGrade();
        if (name) params.set('name', name);
        if (grade) params.set('grade', grade);

        const response = await fetch(`/api/students?${params.toString()}`);
        if (!response.ok) throw new Error('Failed to fetch students');
        const payload = (await response.json()) as StudentListResponse;
        this.students = payload.data ?? [];
        this.pagination.setTotalStudents(payload.count ?? this.students.length);
        this.pagination.setPageFromResponse(payload.offset ?? 0);
        this.initialized = true;
      } catch {
        this.apiError = "Couldn't load students";
      } finally {
        this.loading = false;
      }
    },

    // ─────────────────────────────────────────────────────────────
    // Methods - Filter
    // ─────────────────────────────────────────────────────────────
    applyFilters() {
      this.pagination.resetPage();
      void this.listStudents();
    },
    clearFilters() {
      this.filterTracker.clear();
      this.applyFilters();
    },

    // ─────────────────────────────────────────────────────────────
    // Methods - Pagination
    // ─────────────────────────────────────────────────────────────
    async nextPage() {
      this.pagination.nextPage();
      await this.listStudents();
    },
    async prevPage() {
      this.pagination.prevPage();
      await this.listStudents();
    },

    // ─────────────────────────────────────────────────────────────
    // Methods - Form
    // ─────────────────────────────────────────────────────────────
    resetForm() {
      this.editingId = '';
      this.form = { ...emptyFormValues };
      this.formError = '';
    },
    openCreateModal() {
      this.resetForm();
      window.tui?.dialog.open('student-form');
    },
    async submitForm() {
      this.saving = true;
      const isEdit = this.isEditing();
      const body: StudentRequest = {
        first_name: this.form.firstName,
        last_name: this.form.lastName,
        email: this.form.email,
        age: Number.isFinite(this.form.age) ? this.form.age : 0,
        grade: this.form.grade as Grade,
      };

      try {
        const response = await fetch(isEdit ? `/api/students/${this.editingId}` : '/api/students', {
          method: isEdit ? 'PUT' : 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(body),
        });
        if (!response.ok) {
          const fallback = isEdit ? 'Failed to update student' : 'Failed to create student';
          const payload = (await response.json().catch(() => ({}))) as { error?: string };
          this.formError = payload.error || fallback;
          return;
        }
        window.tui?.dialog.close('student-form');
        this.resetForm();
        this.showToast(
          isEdit ? 'Student updated' : 'Student added',
          isEdit ? 'The student details were saved successfully.' : 'The student record was created successfully.',
          false,
        );
        await this.afterSave(isEdit ? 'update' : 'create');
      } catch (error) {
        const fallback = isEdit ? 'Failed to update student' : 'Failed to create student';
        const message = error instanceof Error ? error.message : fallback;
        this.showToast(isEdit ? 'Update failed' : 'Create failed', message, true);
      } finally {
        this.saving = false;
      }
    },
    async openEditModal(studentId: string) {
      this.resetForm();
      window.tui?.dialog.open('student-form');
      this.editLoading = true;
      try {
        const response = await fetch(`/api/students/${studentId}`);
        if (!response.ok) throw new Error('Failed to fetch student');
        const json = (await response.json()) as { data?: StudentResponse };
        const student = json?.data;
        if (!student) throw new Error('Student not found');

        // Set form values from API
        this.form = {
          firstName: student.first_name,
          lastName: student.last_name,
          email: student.email,
          age: student.age,
          grade: student.grade,
        };
        this.editingId = student.id;
      } catch (error) {
        window.tui?.dialog.close('student-form');
        const message = error instanceof Error ? error.message : 'Failed to fetch student';
        this.showToast('Error', message, true);
      } finally {
        this.editLoading = false;
      }
    },

    // ─────────────────────────────────────────────────────────────
    // Methods - Delete (confirm needs studentPage context)
    // ─────────────────────────────────────────────────────────────
    requestDelete(student: StudentResponse) {
      this.deleteTracker.clearAll();
      this.deleteTracker.startCountdown(student.id, () => void this.confirmPendingDelete());
    },
    async confirmPendingDelete() {
      const id = this.deleteTracker.pendingDeleteId;
      if (!id) return;

      this.deleteTracker.clearAll();
      this.deleteTracker.deletingId = id;

      try {
        const response = await fetch(`/api/students/${id}`, { method: 'DELETE' });
        if (!response.ok) throw new Error('Failed to delete student');
        this.showToast('Student deleted', 'The student record was removed.', false);
        await this.afterSave('delete');
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to delete student';
        this.showToast('Delete failed', message, true);
        this.deleteTracker.deletingId = '';
      }
    },

    // ─────────────────────────────────────────────────────────────
    // Methods - Toast / Feedback
    // ─────────────────────────────────────────────────────────────
    showToast(title: string, message: string, isError: boolean) {
      const toast = cloneToastTemplate(isError);
      if (!toast) return;

      replaceToastText(toast, title, message);
      document.body.appendChild(toast);
    },
    async afterSave(action: StudentMutationAction = 'update') {
      if (action === 'create') {
        const totalAfterCreate = this.pagination.getTotalStudents() + 1;
        this.pagination.setTotalStudents(totalAfterCreate);
        const newPage = Math.max(1, Math.ceil(totalAfterCreate / this.pagination.getPageSize()));
        this.pagination.setPageFromResponse((newPage - 1) * this.pagination.getPageSize());
      }
      await this.listStudents();
    },
    setError(message: string) {
      this.apiError = message;
      this.showToast('Error', message, true);
    },

    // ─────────────────────────────────────────────────────────────
    // Methods - Init
    // ─────────────────────────────────────────────────────────────
    init() {
      void this.listStudents();
    },
  };
}

// =============================================================================
// SECTION 4: Toast Functions
// =============================================================================

function cloneToastTemplate(isError: boolean): HTMLElement | null {
  const templateId = isError ? 'student-error-toast-template' : 'student-success-toast-template';
  const template = document.getElementById(templateId) as HTMLTemplateElement | null;
  if (!template?.content) return null;
  return template.content.cloneNode(true) as HTMLElement;
}

function replaceToastText(toast: HTMLElement, title: string, message: string): void {
  const content = toast.querySelector('span.flex-1');
  if (!content) return;
  const [titleEl, descEl] = content.children;
  if (titleEl) titleEl.textContent = title;
  if (descEl) descEl.textContent = message;
}

// =============================================================================
// SECTION 5: Alpine Registration
// =============================================================================
// Register the studentPage component with Alpine.js

document.addEventListener('alpine:init', () => {
  Alpine.data('studentPage', studentPage);
});

export {};
