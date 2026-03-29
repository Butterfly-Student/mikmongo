# Plan Check: Phase 1 — Foundation & Auth

**Verdict:** FAIL

**Plans checked:** 01-01, 01-02, 01-03
**Issues found:** 5 blockers, 4 warnings, 2 notes

---

## Success Criteria Coverage

| Criterion | Covered By | Status |
|-----------|------------|--------|
| `npm run dev` runs at `dashboard/` without error | 01-01 Task 8 + 01-02 Task 4 | BLOCKED — missing `lucide-react` install, missing `@api/auth.ts` |
| Login superadmin/admin/teknisi → redirect to admin dashboard with summary cards | 01-03 Task 3 (`login.tsx`, `_admin/index.tsx`) | BLOCKED — `@/api/auth` file never created; post-login route target inconsistent |
| Login agent → redirect to agent portal | 01-03 Task 3 (`agent/login.tsx`) | BLOCKED — `agentLogin` imported from non-existent `@/api/auth` |
| Login customer → redirect to customer portal | 01-03 Task 3 (`customer/login.tsx`) | BLOCKED — `customerLogin` imported from non-existent `@/api/auth` |
| Unauthenticated access → redirect to appropriate login page | 01-02 Task 4 (route guards) | PARTIAL — admin guard in `_admin.tsx` works; agent and customer guards present; but `_admin.tsx` conflicts with `_admin/route.tsx` (see Gap 1) |
| Teknisi cannot access user management | 01-02 Task 3 (RBAC) + 01-03 Sidebar | COVERED — `hasPermission("teknisi", "users", "read")` returns false; Sidebar filters by role |
| Dark/light mode toggle works, follows system preference by default | 01-03 Task 1 (ThemeProvider) + FOUC script | PARTIAL — ThemeProvider is correct; `useTheme.ts` stub has broken import (see Gap 4) |

---

## Requirements Coverage

| Requirement | Plan | Status |
|-------------|------|--------|
| SETUP-01: React + Vite + TypeScript in `dashboard/` | 01-01 Tasks 1–4 | Covered |
| SETUP-02: TanStack Router v1 file-based, 3 route trees | 01-01 Task 5, 01-02 Task 4 | Covered |
| SETUP-03: TanStack Query v5 with global QueryClient | 01-02 Task 3 (`queryClient.ts`) + 01-01 Task 7 | Covered |
| SETUP-04: Shadcn/UI + Tailwind v4 with design tokens | 01-01 Task 6 | Covered |
| SETUP-05: Dark mode with system pref + toggle | 01-03 Task 1 (ThemeProvider) | Covered |
| SETUP-06: Axios with JWT interceptors + refresh + error handling | 01-02 Task 2 | Covered |
| SETUP-07: Zustand store for auth state | 01-02 Task 1 | Covered |
| SETUP-08: Shared layout: AppShell, Sidebar, Topbar, mobile nav | 01-03 Task 2 | Covered |
| SETUP-09: Mobile-first responsive breakpoints | 01-03 Task 2 (AppShell `lg:` breakpoint) | Covered |
| SETUP-10: Zod schemas for form validation and API responses | 01-02 Task 3 | Covered |
| AUTH-01: Admin/superadmin/teknisi login via `/api/v1/auth/login` | 01-02 Task 4 (`useAdminLogin`), 01-03 Task 3 (`login.tsx`) | BLOCKED — see Gap 2 |
| AUTH-02: Agent login via `/agent-portal/v1/auth/login` | 01-02 Task 4 (`useAgentLogin`), 01-03 Task 3 | BLOCKED — see Gap 2 |
| AUTH-03: Customer login via `/portal/v1/auth/login` | 01-02 Task 4 (`useCustomerLogin`), 01-03 Task 3 | BLOCKED — see Gap 2 |
| AUTH-04: Refresh token runs automatically | 01-02 Task 2 (interceptors) | Covered |
| AUTH-05: Logout clears token and redirects | 01-02 Task 4 (`useAdminLogout`), 01-03 Topbar | Covered |
| AUTH-06: Protected routes redirect to login when unauthenticated | 01-02 Task 4 (route guards) | PARTIAL — see Gap 1 |
| AUTH-07: RBAC — superadmin/admin/teknisi access matrix | 01-02 Task 3 (`rbac.ts`, tests) | Covered |
| AUTH-08: Customer/agent login cannot access admin dashboard | 01-02 Task 4 (separate guard routes) | Covered — separate auth namespaces |
| ADMIN-01: Dashboard overview with summary cards | 01-03 Task 3 (`_admin/index.tsx`) | BLOCKED — see Gap 2 |
| ADMIN-02: Sidebar nav with all menus, collapsible on mobile | 01-03 Task 2 (Sidebar + Sheet) | Covered |
| ADMIN-03: Topbar with user info, theme toggle, logout | 01-03 Task 2 (Topbar) | Covered |

---

## Gaps Found

---

### Gap 1 — `_admin.tsx` vs `_admin/route.tsx` Route File Conflict (BLOCKER)

**Risk:** HIGH

**Description:**

Plan 01-01 creates `src/routes/_admin/route.tsx` (directory form of the pathless layout). Plan 01-02 creates `src/routes/_admin.tsx` (file form). Plan 01-03 then creates `src/routes/_admin/route.tsx` again, overwriting what 01-01 made.

TanStack Router file-based routing treats `_admin.tsx` and `_admin/route.tsx` as two different registrations of the same pathless layout route `/_admin`. The Vite plugin will throw a duplicate route error or silently use one over the other — behavior is undefined and likely breaks compilation.

Additionally, Plan 01-02 notes: "_admin.tsx currently renders a plain `<Outlet />`. Plan 01-03 will replace this with `<AppShell>`" — but Plan 01-03 only creates `_admin/route.tsx` (directory form), never replacing the `_admin.tsx` file created by Plan 01-02. After all three plans execute, both `_admin.tsx` and `_admin/route.tsx` exist simultaneously.

**Fix:**

Pick one form and use it consistently across all three plans:
- Preferred: directory form `src/routes/_admin/route.tsx` (already used by agent and customer portals).
- Plan 01-01: create `_admin/route.tsx` stub (already correct).
- Plan 01-02: create `_admin/route.tsx` (not `_admin.tsx`) for the auth guard.
- Plan 01-03: replaces `_admin/route.tsx` with AppShell-wrapped version (already correct).
- Delete the `_admin.tsx` from Plan 01-02.

---

### Gap 2 — `src/api/auth.ts` Never Created (BLOCKER)

**Risk:** HIGH

**Description:**

Plan 01-03 creates three login pages (`login.tsx`, `agent/login.tsx`, `customer/login.tsx`). Each imports functions from `@/api/auth`:

```
import { adminLogin } from "@/api/auth"   // login.tsx
import { agentLogin } from "@/api/auth"   // agent/login.tsx
import { customerLogin } from "@/api/auth" // customer/login.tsx
```

No plan in Phase 1 creates `src/api/auth.ts`. The file is not listed in any plan's "Files Created/Modified" section. This will cause a TypeScript compilation error (`Cannot find module '@/api/auth'`) and a runtime import error — the dev server will fail to start.

Note: Plan 01-02 provides `useAdminLogin`, `useAgentLogin`, `useCustomerLogin` hooks that call the Axios clients directly. Plan 01-03's login pages bypass those hooks and call raw `adminLogin()` functions instead — but those functions don't exist anywhere.

**Fix:**

Add a task to Plan 01-02 or Plan 01-03 that creates `src/api/auth.ts` with the three login functions:

```typescript
// src/api/auth.ts
import axios from "axios"
import { AdminLoginResponseSchema, AgentLoginResponseSchema, CustomerLoginResponseSchema } from "@/lib/schemas/auth"

export async function adminLogin(email: string, password: string) {
  const { data } = await axios.post("/api/v1/auth/login", { email, password })
  return AdminLoginResponseSchema.parse(data)
}
export async function agentLogin(email: string, password: string) { ... }
export async function customerLogin(email: string, password: string) { ... }
```

Alternatively, rewrite Plan 01-03 login pages to use `useAdminLogin()` / `useAgentLogin()` / `useCustomerLogin()` hooks from Plan 01-02 (which already exist and call the correct Axios clients with proper error handling).

---

### Gap 3 — `lucide-react` and `@tanstack/router-devtools` Never Installed (BLOCKER)

**Risk:** HIGH

**Description:**

Plan 01-01 Task 2 installs all runtime and dev dependencies. `lucide-react` is not in the install list. Plan 01-02 and Plan 01-03 import from `lucide-react` in over a dozen files:
- `Sidebar.tsx`: `LayoutDashboard`, `Users`, `Router`, `FileText`, `CreditCard`, etc.
- `Topbar.tsx`: `Menu`, `Moon`, `Sun`, `Monitor`, `LogOut`, `User`
- `login.tsx` files: `Loader2`
- `_admin/index.tsx`: `Users`, `Wifi`, `BadgeDollarSign`, `AlertTriangle`

Missing `lucide-react` causes compile-time failures for all of the above files.

Additionally, Plan 01-02's `__root.tsx` imports `TanStackRouterDevtools` from `@tanstack/router-devtools`, which is also not in the Task 2 install list.

**Fix:**

Add to Plan 01-01 Task 2 install command:
```bash
npm install lucide-react
npm install -D @tanstack/router-devtools
```

Or add a new install task at the start of Plan 01-02.

---

### Gap 4 — `useTheme.ts` Stub Has Broken Import in Plan 01-03 (BLOCKER)

**Risk:** HIGH

**Description:**

Plan 01-01 Task 5 creates `src/hooks/useTheme.ts` as a self-contained hook (reads/writes localStorage directly).

Plan 01-03 Task 1 replaces `src/hooks/useTheme.ts` with:

```typescript
import { useContext } from "react"
import { ThemeContext } from "@/components/providers/ThemeProvider"

export { useTheme } from "@/components/providers/ThemeProvider"
```

This file:
1. Imports `useContext` from React but never calls it — TypeScript `noUnusedLocals: true` will error.
2. Imports `ThemeContext` from ThemeProvider — but `ThemeContext` is declared as a `const` inside ThemeProvider.tsx and is NOT exported from that file. This is a TypeScript error: `Module has no exported member 'ThemeContext'`.

The `export { useTheme }` re-export line is valid since `useTheme` is exported from ThemeProvider. But the two broken imports above will cause compilation failure.

**Fix:**

Replace the broken `useTheme.ts` with:
```typescript
// src/hooks/useTheme.ts
export { useTheme } from "@/components/providers/ThemeProvider"
```

Remove the unused `useContext` import and the unexported `ThemeContext` import.

---

### Gap 5 — Route Target Inconsistency: `/` vs `/dashboard` After Admin Login (BLOCKER)

**Risk:** HIGH

**Description:**

Three different post-login navigation targets exist for the admin portal, creating a contradiction:

1. Plan 01-01 `src/routes/index.tsx`: redirects `"/"` → `"/dashboard"` via `throw redirect({ to: "/dashboard" })`.
2. Plan 01-02 `useAdminLogin` hook: calls `navigate({ to: "/dashboard" })` on success.
3. Plan 01-03 `login.tsx`: calls `navigate({ to: "/" })` on success.
4. Plan 01-03 `_admin/index.tsx`: uses `createFileRoute("/_admin/")` which renders at URL `"/"`.
5. Plan 01-01 `_admin/dashboard.tsx`: uses `createFileRoute("/_admin/dashboard")` which renders at URL `"/dashboard"`.

The conflicts:
- Plan 01-01 creates a route at `/dashboard` with a heading "Dashboard".
- Plan 01-03 creates a route at `/` (via `/_admin/`) with the overview page and stat cards.
- The root `index.tsx` redirects `/` → `/dashboard`, which means the overview page at `/_admin/` may never render — the user hits `/`, gets redirected to `/dashboard`, which is the Plan 01-01 stub (not the overview).
- Plan 01-02's `useAdminLogin` navigates to `/dashboard` (stub); Plan 01-03's `login.tsx` navigates to `/` (which redirects back to `/dashboard` anyway).

The intended behavior is that after login the user should see the overview with stat cards. This only works if the overview lives at `/dashboard` (not `/`), OR the `index.tsx` redirect is removed, OR the routing is reorganized.

**Fix:**

Either:
- Option A: Remove `src/routes/index.tsx` (the redirect), and have Plan 01-03's `_admin/index.tsx` be the canonical landing page at `/`. All post-login navigation goes to `"/"`. The Plan 01-01 `_admin/dashboard.tsx` stub is removed or repurposed.
- Option B: Move Plan 01-03's overview page to `_admin/dashboard.tsx` (URL `/dashboard`). Remove Plan 01-01's stub `_admin/dashboard.tsx` in Plan 01-02 or 01-03. Keep `index.tsx` redirect. Both login hooks navigate to `"/dashboard"`.

Option B is the simpler fix — it preserves Plan 01-01's redirect pattern and Plan 01-02's `navigate({ to: "/dashboard" })`.

---

### Gap 6 — `__root.tsx` Rewritten 3 Times with Incompatible Shapes (Warning)

**Risk:** MEDIUM

**Description:**

`src/routes/__root.tsx` is created/replaced in all three plans with different component bodies:

- **Plan 01-01**: Component is `() => <Outlet />`, includes Toaster and ReactQueryDevtools inline in `main.tsx`.
- **Plan 01-02**: Component is `() => (<><Outlet />{DEV && <TanStackRouterDevtools />}{DEV && <ReactQueryDevtools />}</>)`. Toaster is moved to `main.tsx`. No ThemeProvider.
- **Plan 01-03**: Component wraps everything in `<ThemeProvider>`, adds Toaster inside root (not main.tsx). Imports ReactQueryDevtools (Plan 01-02 version also imports it).

The final version from Plan 01-03 is correct. However, each overwrite is a complete file replacement — if execution stops partway through, an inconsistent version of `__root.tsx` is left behind. More critically, the Plan 01-02 version of `__root.tsx` adds a `TanStackRouterDevtools` import from `@tanstack/router-devtools` (not installed — see Gap 3), causing a compile failure even at the Plan 01-02 stage.

**Fix:**

In Plan 01-02, skip the `TanStackRouterDevtools` import in `__root.tsx` since it requires a package not yet installed. Add the devtools only in Plan 01-03 after `@tanstack/router-devtools` is installed.

---

### Gap 7 — Agent/Customer Post-Login Destination: `/agent/dashboard` Does Not Exist as a Named Route (Warning)

**Risk:** MEDIUM

**Description:**

Plan 01-03 agent login navigates to `{ to: "/agent/dashboard" }` on success. Plan 01-01 creates `src/routes/agent/_agentAuth/dashboard.tsx` which registers the route `"/agent/_agentAuth/dashboard"`. TanStack Router file-based routing derives the URL from the directory structure: `agent/_agentAuth/dashboard.tsx` renders at path `/agent/dashboard` (the `_agentAuth` pathless layout doesn't add a URL segment). This should work correctly.

However, Plan 01-02 `useAgentLogin` hook also navigates to `"/agent/dashboard"`. This means if the hook is used instead of the Plan 01-03 login page's inline mutation, the behavior is consistent. But if TypeScript strict mode checks the route string against the generated routeTree, an unrecognized `to` value may raise a type error depending on how the route is registered.

Verify that `createFileRoute("/agent/_agentAuth/dashboard")` (Plan 01-01) produces a URL route of `/agent/dashboard` in TanStack Router's route tree — this is the expected behavior for pathless layouts.

**Fix:** No code change needed, but confirm via the generated `routeTree.gen.ts` after Plan 01-01 execution that `/agent/dashboard` is a valid typed route.

---

### Gap 8 — `main.tsx` Duplication: Plan 01-01 and Plan 01-02 Both Write Full `main.tsx` (Warning)

**Risk:** LOW

**Description:**

Plan 01-01 Task 7 writes a complete `src/main.tsx` with stub context. Plan 01-02 Task 4 replaces it entirely with a new `main.tsx` that uses `useAdminAuthContext`, `queryClient` from `./lib/queryClient`, and the hydration gate. This is intentional and documented.

The risk is that Plan 01-01's `main.tsx` imports `./styles/globals.css` — but Plan 01-02's replacement does NOT import this CSS file. If Plan 01-02's version of `main.tsx` forgets to include `import "./styles/globals.css"`, the entire Tailwind CSS and OKLCH design tokens will be missing from the final build.

Looking at Plan 01-02's `main.tsx` code (lines 1418–1480): it does NOT include `import "./styles/globals.css"`. The stylesheet import is dropped.

**Fix:**

Add `import "./styles/globals.css"` to Plan 01-02's `main.tsx` replacement, before the `routeTree.gen` import.

---

## Verdict Explanation

**FAIL** — 5 blockers prevent the phase goal from being achieved.

The phase goal is: "Project scaffold running at `dashboard/`, all 3 portals can login and access protected routes with RBAC."

The plans collectively cover the right ground — the architecture is correct, RBAC is complete, Zustand store is well-designed, Axios refresh queue is correct, and the layout components are thorough. However, five issues will prevent `npm run dev` from starting without error and prevent any portal from completing a login:

1. **Gap 1** (file conflict): `_admin.tsx` and `_admin/route.tsx` both exist → TanStack Router duplicate route error or undefined behavior.
2. **Gap 2** (missing file): `@/api/auth` is imported by all three login pages but never created → TypeScript compile error, import failure at runtime.
3. **Gap 3** (missing packages): `lucide-react` used by Sidebar, Topbar, all login pages, and overview — never installed → build fails.
4. **Gap 4** (broken import): `useTheme.ts` as written in Plan 01-03 imports an unexported symbol and an unused one → TypeScript errors with `noUnusedLocals: true`.
5. **Gap 5** (routing conflict): Route `/` redirects to `/dashboard`, but the overview page with stat cards is at `/_admin/` (URL `/`) — after login, users land on the `/dashboard` stub with no stat cards.

Gaps 6 and 8 are warnings that may produce subtle bugs (invisible CSS, devtools import crash before install). Gap 7 is informational.

## Required Fixes Before Execution

1. **Plan 01-02**: Replace `_admin.tsx` creation with `_admin/route.tsx` (or document that Plan 01-03 will overwrite `_admin.tsx` with the directory form, and add a cleanup step to delete `_admin.tsx`).

2. **Plan 01-02 or 01-03**: Add task to create `src/api/auth.ts` with `adminLogin`, `agentLogin`, `customerLogin` functions, OR rewrite Plan 01-03 login pages to use hooks from `useAdminAuth.ts`, `useAgentAuth.ts`, `useCustomerAuth.ts` (these already exist in Plan 01-02).

3. **Plan 01-01 Task 2**: Add `lucide-react` to runtime dependencies. Add `@tanstack/router-devtools` to dev dependencies.

4. **Plan 01-03 Task 1**: Fix `useTheme.ts` replacement — remove the broken `useContext` and `ThemeContext` imports, leaving only `export { useTheme } from "@/components/providers/ThemeProvider"`.

5. **Choose and document one routing convention**: Either (A) remove `src/routes/index.tsx` and keep the overview at `/_admin/` (URL `/`), or (B) move the overview to `_admin/dashboard.tsx` (URL `/dashboard`) and remove the stub created in Plan 01-01. Update all `navigate({ to: ... })` calls in login hooks and pages to use the chosen target consistently.

6. **Plan 01-02 `main.tsx`**: Add `import "./styles/globals.css"` to the replacement file.

7. **Plan 01-02 `__root.tsx`**: Remove `TanStackRouterDevtools` import (package not yet installed at this stage). Add it only in Plan 01-03 after installing `@tanstack/router-devtools`.
