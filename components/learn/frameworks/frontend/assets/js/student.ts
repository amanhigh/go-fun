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
// SECTION 2: Global Declarations
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
    totalStudents: 0,

    // ─────────────────────────────────────────────────────────────
    // Filter State
    // ─────────────────────────────────────────────────────────────
    name: '',
    grade: '' as Grade,

    // ─────────────────────────────────────────────────────────────
    // Pagination State
    // ─────────────────────────────────────────────────────────────
    page: 1,
    pageSize: 4,

    // ─────────────────────────────────────────────────────────────
    // Delete State
    // ─────────────────────────────────────────────────────────────
    pendingDeleteSeconds: 0,
    pendingDeleteTimer: null as number | null,

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
    //   - editing: Edit/Create dialog is open
    //   - saving: Form is being submitted
    //   - pendingDeleteId: Row has pending delete countdown
    //   - errorMessage: Error state for API/operation errors
    //   - formError: Inline form validation errors
    //
    // Mutually Exclusive Groups:
    //   1. Load States: initialized ↔ loading ↔ loaded (via initialized)
    //   2. Edit States: idle ↔ editing ↔ saving
    //   3. Error States: hasError() ↔ isFormError() (can coexist)
    // ═══════════════════════════════════════════════════════════════
    initialized: false,
    loading: false,
    saving: false,
    editing: false,
    pendingDeleteId: '',
    errorMessage: '',
    formError: '',

    // ─────────────────────────────────────────────────────────────
    // Computed Properties
    // ─────────────────────────────────────────────────────────────
    get currentPage() { return this.page; },
    get hasFilters() { return this.name !== '' || this.grade !== ''; },
    get filteredStudents() {
      const query = this.name.toLowerCase().trim();
      return this.students.filter((s) => {
        const fullName = `${s.first_name} ${s.last_name}`.toLowerCase();
        return (!query || fullName.includes(query)) && (!this.grade || s.grade === this.grade);
      });
    },
    get totalPages() { return Math.max(1, Math.ceil(this.totalStudents / this.pageSize)); },
    get paginatedStudents() { return this.filteredStudents; },
    isInitialized() { return this.initialized; },
    isEditing() { return this.editing; },
    isLoading() { return this.loading; },
    isSaving() { return this.saving; },
    isEmpty() { return this.filteredStudents.length === 0; },
    isErrored() { return this.errorMessage !== ''; },
    isFormError() { return this.formError !== ''; },
    isPendingDelete(student: StudentResponse) { return this.pendingDeleteId === student.id; },
    hasNext() { return this.currentPage < this.totalPages; },
    hasPrev() { return this.currentPage > 1; },

    // ─────────────────────────────────────────────────────────────
    // Methods - List
    // ─────────────────────────────────────────────────────────────
    async listStudents() {
      this.loading = true;
      this.errorMessage = '';
      try {
        const offset = (this.page - 1) * this.pageSize;
        const params = new URLSearchParams({
          offset: String(offset),
          limit: String(this.pageSize),
        });
        if (this.name.trim() !== '') params.set('name', this.name.trim());
        if (this.grade !== '') params.set('grade', this.grade);

        const response = await fetch(`/api/students?${params.toString()}`);
        if (!response.ok) throw new Error('Failed to fetch students');
        const payload = (await response.json()) as StudentListResponse;
        this.students = payload.data ?? [];
        this.totalStudents = payload.count ?? this.students.length;
        this.pageSize = payload.limit ?? this.pageSize;
        this.page = Math.floor((payload.offset ?? offset) / this.pageSize) + 1;
        this.initialized = true;
      } catch {
        this.errorMessage = "Couldn't load students";
      } finally {
        this.loading = false;
      }
    },

    // ─────────────────────────────────────────────────────────────
    // Methods - Filter
    // ─────────────────────────────────────────────────────────────
    onGradeFilterChange(event: Event) {
      const target = event.target as HTMLInputElement | HTMLSelectElement | null;
      this.grade = (target?.value ?? '') as Grade;
      this.page = 1;
      void this.listStudents();
    },
    clearFilters() {
      this.name = '';
      this.grade = '';
      this.page = 1;
      void this.listStudents();
    },

    // ─────────────────────────────────────────────────────────────
    // Methods - Pagination
    // ─────────────────────────────────────────────────────────────
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.page += 1;
        void this.listStudents();
      }
    },
    prevPage() {
      if (this.currentPage > 1) {
        this.page -= 1;
        void this.listStudents();
      }
    },

    // ─────────────────────────────────────────────────────────────
    // Methods - Form
    // ─────────────────────────────────────────────────────────────
    resetForm() {
      this.editingId = '';
      this.editing = false;
      this.form = { ...emptyFormValues };
      this.formError = '';
    },
    async submitForm() {
      this.saving = true;
      const isEdit = this.editingId !== '';
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
      try {
        // Reset form first
        this.resetForm();

        const response = await fetch(`/api/students/${studentId}`);
        if (!response.ok) throw new Error('Failed to fetch student');
        const json = (await response.json()) as { data?: StudentResponse };
        const student = json?.data;
        if (!student) throw new Error('Student not found');

        // Open dialog
        window.tui?.dialog.open('student-form');

        // Set form values from API
        this.form = {
          firstName: student.first_name,
          lastName: student.last_name,
          email: student.email,
          age: student.age,
          grade: student.grade,
        };
        this.editingId = student.id;
        this.editing = true;
      } catch (error) {
        const message = error instanceof Error ? error.message : 'Failed to fetch student';
        this.showToast('Error', message, true);
      }
    },

    // ─────────────────────────────────────────────────────────────
    // Methods - Delete
    // ─────────────────────────────────────────────────────────────
    requestDelete(student: StudentResponse) {
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
      const nextPage = action === 'create'
        ? Math.max(1, Math.ceil((this.totalStudents + 1) / this.pageSize))
        : this.page;
      this.page = nextPage;
      await this.listStudents();
    },
    setError(message: string) {
      this.errorMessage = message;
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
// SECTION 4: Helper Functions
// =============================================================================

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

// =============================================================================
// SECTION 5: Alpine Registration
// =============================================================================
// Register the studentPage component with Alpine.js

document.addEventListener('alpine:init', () => {
  Alpine.data('studentPage', studentPage);
});

export {};
