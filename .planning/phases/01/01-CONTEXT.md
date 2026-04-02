# Phase 1: Auth & API Foundation - Context

**Gathered:** 2026-04-02
**Status:** Ready for planning

<domain>
## Phase Boundary

Admin, customer, and agent can each log into their respective portal using custom JWT authentication, and all API communication goes through properly configured Axios instances with token refresh. Covers: login/logout/change-password pages, protected route guards, Zustand auth store, Axios interceptors with silent token refresh, three-portal route structure.

</domain>

<decisions>
## Implementation Decisions

### Login Page UX
- **D-01:** Single centered card layout for all three portals (admin, customer, agent). Card contains logo, portal title, email/password fields, and submit button.
- **D-02:** M logo + "MikMongo" title + portal subtitle (e.g., "Admin Portal — sign in to continue", "Customer Portal — sign in to continue", "Agent Portal — sign in to continue").
- **D-03:** Per-field validation errors (below each input) + toast for API-level errors (wrong credentials). No toast for field validation.
- **D-04:** Loading state: spinner icon inside submit button + "Signing in..." text + entire form disabled during request.
- **D-05:** Login page text in Indonesian (labels, buttons, error messages, toast messages). E.g., "Email tidak valid", "Login berhasil", "Email atau password salah".
- **D-06:** Login forms include "Remember me" checkbox and "Forgot password?" link (even though API lacks forgot-password endpoint — link present for UI completeness).
- **D-07:** After successful login, always redirect to the portal's dashboard (not the original URL).

### Existing Code Strategy
- **D-08:** Build on existing code from git HEAD (store, hooks, Axios clients, login pages, route guards, schemas). Restore `dashboard/` directory from git and iterate on what's there. Do not rewrite from scratch.

### Claude's Discretion
- Token storage mechanism details (localStorage key name, storage structure)
- Exact refresh timing (when to proactively refresh vs reactive-only)
- Logout behavior edge cases (multiple tabs, network failures)
- Change password page layout and flow (no specific UX decisions made)
- Customer and agent portal login page visual differences from admin (layout same, branding adapted)

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### API Contract
- `docs/openapi.docs.yml` — Authoritative API contract. Auth endpoints at lines 2462-2580 (/api/v1/auth/login, /refresh, /change-password, /logout, /me). Customer portal login at line 3850, agent portal login at line 4094. Security schemes (BearerAuth, PortalAuth, AgentPortalAuth) at lines 91-106.

### UI Design Contract
- `.planning/phases/01/01-UI-SPEC.md` — Visual and interaction contract for auth pages. Defines spacing, typography, color system, copywriting contract, component usage rules. All login pages must follow this spec.

### Project Constraints
- `.planning/PROJECT.md` — Tech stack constraints (React 19, TypeScript, TanStack, Tailwind, shadcn/ui, Axios, Zustand, Zod). Immutable data patterns required. API contract must match OpenAPI spec exactly.

### Requirements
- `.planning/REQUIREMENTS.md` — Phase 1 requirements: AUTH-01 through AUTH-09, NAV-04. Each has specific acceptance criteria.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets (in git HEAD — dashboard/ directory)
- **Zustand store** (`dashboard/src/store/`): Three auth slices (admin, customer, agent) with persist middleware, hydration gate (`isHydrated`), vanilla getters for Axios interceptors. Partialize excludes actions. Well-structured.
- **Axios clients** (`dashboard/src/lib/axios/`): Three separate clients (admin, customer, agent) with request interceptors (attach Bearer token) and response interceptors (401 → refresh token → retry queue). Uses plain axios for refresh calls to avoid interceptor loops.
- **Auth hooks** (`dashboard/src/hooks/`): `useAdminLogin`, `useAdminLogout`, `useAdminUser`, `useAdminRole` with TanStack Query mutations. Context-facing hooks (`useAuth.ts`) return RouterContext-compatible shapes.
- **Zod schemas** (`dashboard/src/lib/schemas/auth.ts`): ApiResponseSchema, LoginRequest/Response schemas for all three portals, form validation schemas.
- **Route structure** (`dashboard/src/routes/`): `__root.tsx` with typed RouterContext, `_admin/route.tsx` with beforeLoad guard, separate login pages for admin/customer/agent.
- **UI components** (`dashboard/src/components/ui/`): button, input, label, card — all shadcn/ui primitives ready to use.
- **API functions** (`dashboard/src/api/auth.ts`): Raw login functions using plain axios (no interceptors). Separate endpoints for each portal.

### Established Patterns
- **Feature-based structure**: `data/schema.ts`, `data/[entity].ts`, `components/`, `index.tsx` convention
- **Immutability**: New objects via set(), never mutate existing state
- **Auth context injection**: Router context populated from Zustand in `main.tsx`, gated by `isHydrated`
- **Error handling**: Toast notifications (sonner) for user-facing errors, per-field validation for form inputs

### Integration Points
- Route guards use `context.adminAuth.isAuthenticated` from RouterContext
- Axios interceptors read tokens from `useStore.getState()` (vanilla getters)
- Login mutations call `adminSetTokens` + `adminSetUser` then navigate
- Token refresh updates store via `adminAuthActions.setTokens()`

### Known Issues in Existing Code
- Refresh endpoint schema mismatch: OpenAPI returns `{ token, refresh_token }` but existing code reads `{ access_token, refresh_token }`. Must fix.
- `dashboard/src/api/types.ts` is empty (placeholder `export {}`). Needs content.

</code_context>

<specifics>
## Specific Ideas

- Login page uses the "M" branded logo inside a rounded square — matches existing code pattern
- Indonesian UI text throughout (e.g., "Login berhasil", "Email atau password salah", "Email tidak valid")
- Include "Remember me" and "Forgot password?" in login forms for UI completeness even without backend forgot-password support
- Post-login always redirects to dashboard (not deep-link preserved URL)

</specifics>

<deferred>
## Deferred Ideas

- "Remember me" functionality has no backend support yet — checkbox present but behavior is cosmetic for now
- "Forgot password?" link has no endpoint — could redirect to a static page or show a toast saying "Contact admin" for now
- Proactive token refresh (refresh before expiry) vs reactive-only (refresh on 401) — deferred to implementation

</deferred>

---
*Phase: 01-auth-api-foundation*
*Context gathered: 2026-04-02*
