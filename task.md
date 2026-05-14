# Testing Roadmap & Tasks

Whenever an tests for an module is completed add `[x]` before the module name. Don't change the order of the tasks.

## 🎨 Frontend (apps/client)

### 🏗️ Infrastructure

- [x] **Setup Vitest and React Testing Library**
  - [x] Configure `vitest.config.ts` with `jsdom` environment.
  - [x] Create `vitest.setup.ts` with Next.js router mocks and `jest-dom` extensions.
  - [x] Add `test` script to `package.json`.

### 🧩 Component Testing

- [x] **`SubmitButton.tsx`**
  - [x] Verify children rendering.
  - [x] Ensure loader displays during `isLoading` state.
  - [x] Validate disabled state when `isLoading` or `disabled` prop is true.
- [x] **`AuthFormLogin.tsx`**
  - [x] Test form validation (email format, required fields).
  - [x] Mock API calls and verify success redirect.
  - [x] Verify error message display on failed authentication.
- [x] **`AuthFormSignUp.tsx`**
  - [x] Test registration form validation and submission.
  - [x] Verify loading states and error handling.
- [ ] **`AuthFormReset.tsx` & `AuthFormResetPassword.tsx`**
  - [ ] Test password reset request flow.
  - [ ] Verify password confirmation validation.
- [x] **`FormRender.tsx` (Critical)**
  - [x] Test dynamic rendering of different field types (Input, Textarea, Select, MultiSelect).
  - [x] Verify field validation logic (required fields, regex patterns).
  - [x] Ensure `onSubmit` is called with correct data structure.
- [x] **`FileUpload.tsx`**
  - [x] Validate file type restrictions (image, video, audio).
  - [x] Test file size limit enforcement.
  - [x] Mock upload progress and verify UI updates.
- [x] **`response.schemas.ts` (Validation Logic)**
  - [x] Test `ResponseItemSchema` for all `FieldType` variants (text, email, file, etc.).
  - [x] Verify `isRequired` and `isMultiple` logic for different field types.
  - [x] Validate `FormSubmissionSchema` structural integrity.
- [ ] **`user-response/page.tsx` (Server Component)**
  - [ ] Mock API responses for `getFormForResponse`, `getFieldsForResponse`, and `verifyInvitation`.
  - [ ] Verify correct passing of props to `FormRender`.
- [ ] **`FormInputwrapper.tsx`**
  - [ ] Test error message display.
  - [ ] Verify label and required indicator rendering.
- [ ] **`FormShortInput.tsx`, `FormLongInput.tsx`, `FormInputSelect.tsx`, `FormInputCheckbox.tsx`, `FormStarRating.tsx`**
  - [ ] Unit tests for individual field types.
  - [ ] Verify value propagation to parent form.
- [ ] **`MultiSelect.tsx` (Common Component)**
  - [ ] Test adding/removing options.
  - [ ] Verify search/filtering within options.
- [ ] **`DataTable.tsx` & `Pagination.tsx`**
  - [ ] Test data rendering and empty states.
  - [ ] Verify sorting and pagination triggers.
- [ ] **`Dashboard.tsx`**
  - [ ] Test summary statistics rendering.
  - [ ] Verify chart integration (mocked).
- [ ] **`UpsertForm.tsx`**
  - [ ] Test form creation and editing logic.
  - [ ] Verify field addition/removal/ordering.
- [ ] **`Profile.tsx` & `ProfilePic.tsx`**
  - [ ] Test profile data display.
  - [ ] Verify profile picture upload and update flow.

### 🪝 Hooks & State Management

- [x] **`me.slice.ts` & `auth.slice.ts` (Redux Toolkit)**
  - [x] Test `setMe`, `setMyName`, `setTokens`, and `clearTokens` reducers.
  - [x] Verify initial state and selector outputs.
  - [x] Test cookie synchronization for auth tokens.
- [ ] **`useAuth` Hook**
  - [ ] Test authentication status logic.
  - [ ] Verify redirection for unauthorized users.
  - [ ] Mock session storage/cookies for persistence tests.

### 🔌 API Services

- [ ] **`services/api/responses/routes.ts`**
  - [ ] Mock API requests and verify response payloads for submissions.
- [ ] **`services/api/invitations/routes.ts`**
  - [ ] Test token verification and invitation data retrieval.
- [ ] **`services/api/dashboard/routes.ts`**
  - [ ] Verify dashboard statistics fetching logic.
- [ ] **`services/api/auth/routes.ts`**
  - [ ] Test login, registration, and session management API calls.
- [ ] **`services/api/forms/routes.ts`**
  - [ ] Test form creation, fetching, and updating API integration.
- [ ] **`services/api/users/routes.ts`**
  - [ ] Test user profile fetching and updating API integration.

---

## ⚙️ Backend (apps/server)

### 🧪 Service Layer (Business Logic)

- [x] **`users.service.go`**
  - [x] Test `GetUserDetailsById` with mocked repository.
  - [x] Test `CreateUser` with parameter mapping.
  - [ ] Test `UpdateUserProfilePic` including file validation logic.
- [x] **`forms.service.go`**
  - [x] Test `CreateForm` with nested field structure.
  - [x] Test `GetFormByID` handling 404 cases.
  - [x] Test `DeleteForm` and associated cleanup logic.
- [x] **`response.service.go`**
  - [x] Test form response validation and processing.
  - [x] Verify response data mapping and database insertion.
- [x] **`invitations.services.go`**
  - [x] Test unique link generation and tracking.
  - [x] Verify invitation token validity logic.
  - [x] Test `DeleteInvitation` logic.
- [x] **`dashboard.service.go`**
  - [x] Test data aggregation for dashboard statistics.
- [x] **`jwt.services.go`**
  - [x] Test token generation, signing, and verification.
- [x] **`auth.services.go`** & **`google_auth.services.go`**
  - [x] Test authentication flow and OAuth integration.
  - [x] Test `RequestResetPassword` logic.
- [x] **`mails.service.go`** (Skipped as per user request)
- [x] **`minio.services.go`** (Skipped as per user request)

### 🌐 Handler Layer (API Endpoints)

- [x] **`auth.handler.go`**
  - [x] Test `Login` endpoint with valid/invalid credentials (using Gin test context).
  - [x] Verify JWT cookie setting and response body.
  - [x] Test `GoogleAuth` callback handling.
- [x] **`forms.handler.go`**
  - [x] Test `POST /forms` for authorization and payload validation.
  - [x] Test `GET /forms/:id` for correct response mapping.
- [x] **`response.handler.go`**
  - [x] Test `POST /response` (CreateSubmission) endpoint.
  - [x] Test `GET /response/:formId` (GetSubmissions) endpoint.
- [x] **`invitations.handler.go`**
  - [x] Test `DELETE /invitations/:id` (DeleteInvitation) endpoint.
- [x] **`dashboard.handler.go`**
  - [x] Test `GET /dashboard` (GetDashboardData) endpoint.
- [ ] **`file_upload.handler.go`**
  - [ ] Test file upload endpoints, including validation and storage integration.
- [x] **`users.handler.go`**
  - [x] Test user profile retrieval and update endpoints.

### 🎮 Controller Layer

- [ ] **Auth, Form, Response, Dashboard, Users, Invitations Controllers**
  - [ ] Test request binding and service calls.
  - [ ] Verify error handling and response formatting.
  - [ ] Skip if only boilerplate logic is present.

### 🛡️ Middleware

- [x] **Auth, Response Middlewares**
  - [x] Test token validation and context enrichment.
  - [x] Verify access control logic.
  - [ ] Skip if no custom logic is implemented.

### 🏗️ Infrastructure & Background Tasks

- [x] **`cache.go`**
  - [x] Test cache operations (Get, Set, Delete).
- [ ] **`scheduler/` (Invitation, User, Form)**
  - [ ] Test job scheduling logic and execution triggers.
- [/] **`workers/` (Form, User, Invitations)**
  - [x] Test background processing logic and error retries.
- [ ] **`webhook/` (Invitations)**
  - [ ] Test webhook payload handling and delivery logic.

### 🛠️ Utilities & Serialisers

- [x] **`serialisers/` (Form Fields, Users)**
  - [x] Test data transformation and masking logic.
- [x] **`utils/` (Helpers, Mapper, Google Auth, File Validator)**
  - [x] Test utility functions (validation, mapping, etc.).
  - [ ] Note: `resend.go` (Mail utility) skipped as per user request.
- [x] **`services/templates`**
  - [x] Test template rendering and variable injection.

### 💾 Repository Layer (Database Operations)

- [x] **`users.repository.go`**
  - [x] Use `sqlmock` or a test DB to verify SQLC generated queries.
  - [x] Test `SoftDeleteUser` and verify record persistence.
- [x] **`forms.repository.go`**
  - [x] Implement unit tests for CRUD operations (GetFormByID, SoftDeleteForm).
- [x] **`response.repository.go`**
  - [x] Implement unit tests for submission data access (CreateSubmission, GetSubmissionCount).
- [x] **`invitations.repository.go`**
  - [x] Test invitation record creation and querying.
- [x] **`dashboard.repository.go`**
  - [x] Verify data aggregation queries for dashboard metrics.
