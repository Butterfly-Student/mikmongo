# Phase 1 Research: Foundation & Auth

**Researched:** 2026-03-30
**Domain:** React 19 + Vite 6 + TanStack Router v1 + Shadcn/UI + Tailwind v4 + Zustand v5
**Confidence:** HIGH (verified against npm registry and official docs)

---

## Summary

This phase scaffolds a React 19 dashboard at `dashboard/` inside the existing Go repo, establishes three separate portal route trees (Admin at `/`, Agent at `/agent`, Customer at `/customer`) using TanStack Router v1 file-based routing, and wires up JWT auth with auto-refresh for each portal.

**Key finding 1:** TanStack Router v1 uses `_` prefix for pathless layout routes — `_admin/route.tsx` creates a layout wrapper with auth guard that does NOT add a path segment. Child files like `_admin/dashboard.tsx` render at `/dashboard`, not `/_admin/dashboard`.

**Key finding 2:** Tailwind v4 is CSS-first. There is NO `tailwind.config.js` anymore. `@import "tailwindcss"` replaces the old 3-line setup. Dark mode is configured via `@custom-variant dark` in the CSS file, not in a config key.

**Key finding 3:** Shadcn/UI now fully supports Tailwind v4. The `npx shadcn@latest init` command auto-detects v4, generates `@theme inline` CSS variables, and uses OKLCH color space instead of HSL. The `new-york` style is now the default.

**Key finding 4:** The Axios refresh token queue pattern is well-established. A `isRefreshing` flag + `failedQueue` array prevents multiple simultaneous refresh calls. Requests arriving during refresh are held in the queue and replayed with the new token.

**Key finding 5:** Zustand v5 `persist` middleware stores tokens in localStorage automatically via `partialize`. The vanilla store pattern (`createStore` from `zustand/vanilla`) allows using getters outside React components — required for Axios interceptors that need tokens without a hook.

**Primary recommendation:** Use file-based routing with `_portal/route.tsx` pathless layout files for each portal's auth guard + layout shell. Pass Zustand auth state into router context via `createRootRouteWithContext` for type-safe `beforeLoad` guards.

---

<phase_requirements>
## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| SETUP-01 | React + Vite + TypeScript in `dashboard/` | Vite 6 + @vitejs/plugin-react 6.0 — verified |
| SETUP-02 | TanStack Router v1 file-based, 3 route trees | `_admin/`, `_agent/`, `_customer/` pathless layout routes |
| SETUP-03 | TanStack Query v5 with global QueryClient | `QueryClient` with `staleTime`, `gcTime` defaults |
| SETUP-04 | Shadcn/UI + Tailwind v4 with design tokens | `npx shadcn@latest init` detects v4 automatically |
| SETUP-05 | Dark/light mode with system pref + manual toggle | `@custom-variant dark` + `localStorage` + inline head script |
| SETUP-06 | Axios with JWT interceptor + refresh + error handling | Queue pattern with `isRefreshing` flag |
| SETUP-07 | Zustand store for auth state per portal | Slice pattern + `persist` middleware with `partialize` |
| SETUP-08 | AppShell, Sidebar, Topbar, mobile nav | Shadcn/UI Sheet for mobile, layout in portal route file |
| SETUP-09 | Mobile-first responsive 375px → 1440px | Tailwind breakpoints: `sm:640 md:768 lg:1024 xl:1280` |
| SETUP-10 | Zod schemas for forms and API responses | Zod v4.3 — verify API response envelope schema |
| AUTH-01 | Admin login via `POST /api/v1/auth/login` | Axios instance for admin API base URL |
| AUTH-02 | Agent login via `POST /agent-portal/v1/auth/login` | Separate Axios instance, separate Zustand slice |
| AUTH-03 | Customer login via `POST /portal/v1/auth/login` | Separate Axios instance, separate Zustand slice |
| AUTH-04 | Auto refresh token before JWT expires | Response interceptor 401 + queue pattern |
| AUTH-05 | Logout clears tokens + redirects to login | `clearTokens` action + TanStack Router `navigate` |
| AUTH-06 | Protected routes redirect to login | `beforeLoad` in `_portal/route.tsx` layout file |
| AUTH-07 | RBAC: superadmin/admin/teknisi | Role check in `beforeLoad`, store role in Zustand |
| AUTH-08 | Cross-portal access isolation | Separate route trees, separate auth stores |
| ADMIN-01 | Overview with summary cards | Phase 1 lays layout shell; data comes from TanStack Query |
| ADMIN-02 | Sidebar with collapsible mobile nav | Shadcn/UI Sheet + Sidebar components |
| ADMIN-03 | Topbar with user info, theme toggle, logout | Part of `_admin/route.tsx` layout shell |
</phase_requirements>

---

## Standard Stack

### Core

| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| react | 19.2.4 | UI framework | Latest stable, concurrent features |
| react-dom | 19.2.4 | DOM renderer | Paired with react |
| vite | 8.0.3 | Build tool | Official Vite 6 — fastest HMR |
| @vitejs/plugin-react | 6.0.1 | React transform + Fast Refresh | Use this, NOT SWC (see Gotchas) |
| typescript | ^5.7 | Type safety | Vite scaffolds this |
| @tanstack/react-router | 1.168.8 | File-based routing, type-safe | Project requirement |
| @tanstack/router-plugin | 1.167.9 | Vite plugin for route generation | Required for file-based routing |
| @tanstack/react-query | 5.95.2 | Server state, caching | Project requirement |
| @tanstack/react-query-devtools | 5.x | Query debug | Dev only |
| zustand | 5.0.12 | Auth/client state | Project requirement |
| axios | 1.14.0 | HTTP client | Project requirement |
| zod | 4.3.6 | Schema validation | Project requirement |
| sonner | 2.0.7 | Toast notifications | Project requirement |
| tailwindcss | 4.2.2 | Utility CSS | Project requirement |
| @tailwindcss/vite | 4.x | Tailwind v4 Vite plugin | Replaces PostCSS pipeline |
| shadcn (CLI) | 4.1.1 | Component scaffolder | `npx shadcn@latest init` |

### Supporting

| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| @tanstack/react-table | 8.21.3 | Data tables | Phase 2+ table pages |
| @tanstack/react-form | 1.28.5 | Form state | Any form with complex validation |
| @tanstack/react-virtual | 3.13.23 | Virtualization | Long lists (subscriptions, invoices) |
| lucide-react | latest | Icons | Ships with Shadcn/UI |
| tw-animate-css | latest | Animations | Replaces `tailwindcss-animate` in v4 |
| date-fns | ^3 | Date formatting | Invoice dates, subscription periods |
| recharts | ^2 | Charts | Overview page metrics (Phase 4) |
| @types/node | latest | Node types for `path` in vite.config | Required for path alias |

### Alternatives Considered

| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| @vitejs/plugin-react | @vitejs/plugin-react-swc | SWC repo archived July 2025; SWC plugin migrated to main repo. Use standard plugin. |
| Axios | ky | Axios has broader ecosystem for interceptor patterns; project requirement |
| localStorage token storage | httpOnly cookie | httpOnly requires backend support (Set-Cookie); localStorage is simpler for pure SPA |
| Zustand slice pattern | Separate create() per portal | Slices allow cross-store access; cleaner single store |

**Installation:**

```bash
cd dashboard
npm create vite@latest . -- --template react-ts
npm install @tanstack/react-router @tanstack/router-plugin @tanstack/react-query @tanstack/react-query-devtools zustand axios zod sonner date-fns
npm install -D tailwindcss @tailwindcss/vite @types/node
npx shadcn@latest init
```

---

## Architecture Patterns

### Recommended Project Structure

```
dashboard/
├── public/
├── src/
│   ├── routes/                      # TanStack Router file-based routes
│   │   ├── __root.tsx               # Root route — createRootRouteWithContext
│   │   ├── index.tsx                # Redirect to /dashboard
│   │   ├── login.tsx                # Admin login page
│   │   ├── _admin/                  # Pathless layout: admin portal
│   │   │   ├── route.tsx            # Auth guard + AppShell layout
│   │   │   ├── dashboard.tsx        # /dashboard
│   │   │   ├── customers/
│   │   │   │   ├── index.tsx        # /customers
│   │   │   │   └── $customerId.tsx  # /customers/:customerId
│   │   │   └── ...
│   │   ├── agent/                   # Path segment: /agent
│   │   │   ├── login.tsx            # /agent/login
│   │   │   └── _agentAuth/          # Pathless layout: agent auth guard
│   │   │       ├── route.tsx        # Agent auth guard + AgentShell layout
│   │   │       └── dashboard.tsx    # /agent/dashboard
│   │   └── customer/                # Path segment: /customer
│   │       ├── login.tsx            # /customer/login
│   │       └── _customerAuth/       # Pathless layout: customer auth guard
│   │           ├── route.tsx        # Customer auth guard + CustomerShell layout
│   │           └── dashboard.tsx    # /customer/dashboard
│   ├── components/
│   │   ├── ui/                      # Shadcn/UI generated components
│   │   ├── layout/
│   │   │   ├── AppShell.tsx         # Admin sidebar + content area
│   │   │   ├── Sidebar.tsx          # Admin sidebar nav
│   │   │   ├── Topbar.tsx           # Top bar with user + theme toggle
│   │   │   ├── AgentShell.tsx       # Agent portal layout
│   │   │   └── CustomerShell.tsx    # Customer portal layout
│   │   └── shared/                  # Cross-portal shared components
│   ├── lib/
│   │   ├── axios/
│   │   │   ├── admin-client.ts      # Axios instance for /api/v1
│   │   │   ├── agent-client.ts      # Axios instance for /agent-portal/v1
│   │   │   └── customer-client.ts   # Axios instance for /portal/v1
│   │   └── utils.ts                 # cn() helper from shadcn
│   ├── store/
│   │   ├── index.ts                 # Combined Zustand store
│   │   ├── slices/
│   │   │   ├── adminAuthSlice.ts    # Admin token + user + role
│   │   │   ├── agentAuthSlice.ts    # Agent token + user
│   │   │   └── customerAuthSlice.ts # Customer token + user
│   │   └── types.ts                 # Store type definitions
│   ├── hooks/
│   │   ├── useAdminAuth.ts          # Admin auth selectors
│   │   ├── useAgentAuth.ts          # Agent auth selectors
│   │   └── useTheme.ts              # Dark/light mode toggle
│   ├── api/
│   │   ├── auth.ts                  # Login/logout/refresh API calls
│   │   └── types.ts                 # API response types + Zod schemas
│   ├── styles/
│   │   └── globals.css              # @import "tailwindcss" + CSS vars + @custom-variant
│   └── main.tsx                     # App entry: QueryClientProvider + RouterProvider
├── index.html                       # Inline theme script in <head>
├── vite.config.ts
├── tsconfig.json
├── tsconfig.app.json
└── components.json                  # Shadcn/UI config (generated)
```

---

### Pattern 1: Vite + TanStack Router Config

**What:** Wire up all Vite plugins in correct order with path alias and API proxy.
**When to use:** Initial project scaffold.

```typescript
// vite.config.ts
// Source: https://tanstack.com/router/v1/docs/framework/react/installation/with-vite
import path from "path"
import { defineConfig } from "vite"
import react from "@vitejs/plugin-react"
import tailwindcss from "@tailwindcss/vite"
import { TanStackRouterVite } from "@tanstack/router-plugin/vite"

export default defineConfig({
  plugins: [
    // CRITICAL: TanStackRouterVite MUST come before react()
    TanStackRouterVite({ target: "react", autoCodeSplitting: true }),
    react(),
    tailwindcss(),
  ],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "/agent-portal": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
      "/portal": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
})
```

---

### Pattern 2: Root Route with Typed Context

**What:** `createRootRouteWithContext` makes auth state available in every route's `beforeLoad`.
**When to use:** Always — required for type-safe auth guards.

```typescript
// src/routes/__root.tsx
// Source: https://tanstack.com/router/v1/docs/framework/react/guide/router-context
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router"
import { QueryClient } from "@tanstack/react-query"

export interface RouterContext {
  adminAuth: {
    isAuthenticated: boolean
    role: string | null
    accessToken: string | null
  }
  agentAuth: {
    isAuthenticated: boolean
    accessToken: string | null
  }
  customerAuth: {
    isAuthenticated: boolean
    accessToken: string | null
  }
  queryClient: QueryClient
}

export const Route = createRootRouteWithContext<RouterContext>()({
  component: () => <Outlet />,
})
```

```typescript
// src/main.tsx
import React from "react"
import ReactDOM from "react-dom/client"
import { RouterProvider, createRouter } from "@tanstack/react-router"
import { QueryClient, QueryClientProvider } from "@tanstack/react-query"
import { routeTree } from "./routeTree.gen"
import { useAdminAuthContext, useAgentAuthContext, useCustomerAuthContext } from "./hooks/useAuth"

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5,   // 5 minutes
      gcTime: 1000 * 60 * 10,      // 10 minutes
      retry: 2,
      refetchOnWindowFocus: false,
    },
  },
})

const router = createRouter({
  routeTree,
  context: {
    adminAuth: undefined!,
    agentAuth: undefined!,
    customerAuth: undefined!,
    queryClient,
  },
})

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router
  }
}

function App() {
  const adminAuth = useAdminAuthContext()
  const agentAuth = useAgentAuthContext()
  const customerAuth = useCustomerAuthContext()

  return (
    <RouterProvider
      router={router}
      context={{ adminAuth, agentAuth, customerAuth, queryClient }}
    />
  )
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <App />
    </QueryClientProvider>
  </React.StrictMode>
)
```

---

### Pattern 3: Pathless Layout Route with Auth Guard

**What:** `_admin/route.tsx` wraps all admin routes with auth check + layout shell — no URL impact.
**When to use:** Per-portal auth guard and layout injection.

```typescript
// src/routes/_admin/route.tsx
// Source: https://tanstack.com/router/v1/docs/framework/react/guide/authenticated-routes
import { createFileRoute, redirect, Outlet } from "@tanstack/react-router"
import { AppShell } from "@/components/layout/AppShell"

export const Route = createFileRoute("/_admin")({
  beforeLoad: ({ context, location }) => {
    if (!context.adminAuth.isAuthenticated) {
      throw redirect({
        to: "/login",
        search: { redirect: location.href },
      })
    }
  },
  component: () => (
    <AppShell>
      <Outlet />
    </AppShell>
  ),
})
```

```typescript
// src/routes/_admin/dashboard.tsx  → renders at URL: /dashboard
import { createFileRoute } from "@tanstack/react-router"

export const Route = createFileRoute("/_admin/dashboard")({
  component: DashboardPage,
})

function DashboardPage() {
  return <div>Dashboard</div>
}
```

**RBAC guard — role-aware beforeLoad:**

```typescript
// src/routes/_admin/_superadmin/route.tsx  (nested pathless layout for superadmin-only)
import { createFileRoute, redirect, Outlet } from "@tanstack/react-router"

export const Route = createFileRoute("/_admin/_superadmin")({
  beforeLoad: ({ context }) => {
    if (context.adminAuth.role !== "superadmin") {
      throw redirect({ to: "/dashboard" })
    }
  },
  component: () => <Outlet />,
})
```

---

### Pattern 4: Three Separate Axios Instances

**What:** Each portal gets its own Axios instance with its own baseURL and token source.
**When to use:** Ensures that admin tokens are never sent to the customer portal API and vice versa.

```typescript
// src/lib/axios/admin-client.ts
import axios, { AxiosError, InternalAxiosRequestConfig } from "axios"
import { getAdminAccessToken, getAdminRefreshToken, adminAuthActions } from "@/store"

const adminClient = axios.create({
  baseURL: "/api/v1",
  headers: { "Content-Type": "application/json" },
})

let isRefreshing = false
let failedQueue: Array<{
  resolve: (token: string) => void
  reject: (err: AxiosError) => void
}> = []

const processQueue = (error: AxiosError | null, token: string | null) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) reject(error)
    else resolve(token!)
  })
  failedQueue = []
}

// Request: attach token
adminClient.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = getAdminAccessToken()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

// Response: handle 401, refresh, retry
adminClient.interceptors.response.use(
  (res) => res,
  async (error: AxiosError) => {
    const original = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    if (error.response?.status === 401 && !original._retry) {
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject })
        }).then((token) => {
          original.headers.Authorization = `Bearer ${token}`
          return adminClient(original)
        })
      }

      original._retry = true
      isRefreshing = true

      try {
        const refreshToken = getAdminRefreshToken()
        // Use plain axios (not adminClient) to avoid interceptor loop
        const { data } = await axios.post("/api/v1/auth/refresh", { refresh_token: refreshToken })
        const newToken = data.data.access_token
        adminAuthActions.setTokens(newToken, data.data.refresh_token)
        processQueue(null, newToken)
        original.headers.Authorization = `Bearer ${newToken}`
        return adminClient(original)
      } catch (refreshError) {
        processQueue(refreshError as AxiosError, null)
        adminAuthActions.clearAuth()
        window.location.href = "/login"
        return Promise.reject(refreshError)
      } finally {
        isRefreshing = false
      }
    }

    return Promise.reject(error)
  }
)

export { adminClient }
```

---

### Pattern 5: Zustand Auth Store with Slices

**What:** One Zustand store with three auth slices — one per portal. Vanilla store for use outside React components (Axios interceptors need tokens synchronously).
**When to use:** Auth state management across the entire app.

```typescript
// src/store/slices/adminAuthSlice.ts
// Source: https://zustand.docs.pmnd.rs/guides/slices-pattern
import { StateCreator } from "zustand"
import type { StoreState } from "../types"

export interface AdminUser {
  id: string
  email: string
  role: "superadmin" | "admin" | "teknisi"
  full_name: string
}

export interface AdminAuthSlice {
  adminAccessToken: string | null
  adminRefreshToken: string | null
  adminUser: AdminUser | null
  adminIsAuthenticated: boolean
  adminSetTokens: (access: string, refresh: string) => void
  adminSetUser: (user: AdminUser) => void
  adminClearAuth: () => void
}

export const createAdminAuthSlice: StateCreator<
  StoreState,
  [["zustand/persist", unknown], ["zustand/devtools", never]],
  [],
  AdminAuthSlice
> = (set) => ({
  adminAccessToken: null,
  adminRefreshToken: null,
  adminUser: null,
  adminIsAuthenticated: false,
  adminSetTokens: (access, refresh) =>
    set({ adminAccessToken: access, adminRefreshToken: refresh, adminIsAuthenticated: true }),
  adminSetUser: (user) => set({ adminUser: user }),
  adminClearAuth: () =>
    set({ adminAccessToken: null, adminRefreshToken: null, adminUser: null, adminIsAuthenticated: false }),
})
```

```typescript
// src/store/index.ts
import { create } from "zustand"
import { persist, devtools } from "zustand/middleware"
import { createAdminAuthSlice, AdminAuthSlice } from "./slices/adminAuthSlice"
import { createAgentAuthSlice, AgentAuthSlice } from "./slices/agentAuthSlice"
import { createCustomerAuthSlice, CustomerAuthSlice } from "./slices/customerAuthSlice"

export type StoreState = AdminAuthSlice & AgentAuthSlice & CustomerAuthSlice

export const useStore = create<StoreState>()(
  devtools(
    persist(
      (...a) => ({
        ...createAdminAuthSlice(...a),
        ...createAgentAuthSlice(...a),
        ...createCustomerAuthSlice(...a),
      }),
      {
        name: "mikmongo-auth",
        partialize: (state) => ({
          // Only persist tokens — not actions
          adminAccessToken: state.adminAccessToken,
          adminRefreshToken: state.adminRefreshToken,
          adminUser: state.adminUser,
          adminIsAuthenticated: state.adminIsAuthenticated,
          agentAccessToken: state.agentAccessToken,
          agentUser: state.agentUser,
          agentIsAuthenticated: state.agentIsAuthenticated,
          customerAccessToken: state.customerAccessToken,
          customerUser: state.customerUser,
          customerIsAuthenticated: state.customerIsAuthenticated,
        }),
      }
    ),
    { name: "MikMongo Store", enabled: import.meta.env.DEV }
  )
)

// Vanilla getters for use outside React (Axios interceptors)
export const getAdminAccessToken = () => useStore.getState().adminAccessToken
export const getAdminRefreshToken = () => useStore.getState().adminRefreshToken
export const adminAuthActions = {
  setTokens: (a: string, r: string) => useStore.getState().adminSetTokens(a, r),
  clearAuth: () => useStore.getState().adminClearAuth(),
}
// Repeat getAgentAccessToken, getCustomerAccessToken, etc.
```

---

### Pattern 6: Tailwind v4 + Dark Mode + Shadcn/UI CSS

**What:** Full CSS setup replacing tailwind.config.js.
**When to use:** Initial `src/styles/globals.css` setup.

```css
/* src/styles/globals.css */
/* Source: https://ui.shadcn.com/docs/tailwind-v4 + https://tailwindcss.com/docs/dark-mode */

@import "tailwindcss";
@import "tw-animate-css";

/* Enable class-based dark mode for manual toggle support */
@custom-variant dark (&:where(.dark, .dark *));

/* Shadcn/UI design tokens — OKLCH (generated by `npx shadcn@latest init`) */
@layer base {
  :root {
    --background: oklch(1 0 0);
    --foreground: oklch(0.145 0 0);
    --card: oklch(1 0 0);
    --card-foreground: oklch(0.145 0 0);
    --primary: oklch(0.205 0 0);
    --primary-foreground: oklch(0.985 0 0);
    --muted: oklch(0.97 0 0);
    --muted-foreground: oklch(0.556 0 0);
    --border: oklch(0.922 0 0);
    --radius: 0.625rem;
    /* ...full token set generated by shadcn init... */
  }

  .dark {
    --background: oklch(0.145 0 0);
    --foreground: oklch(0.985 0 0);
    --primary: oklch(0.922 0 0);
    --primary-foreground: oklch(0.205 0 0);
    --muted: oklch(0.269 0 0);
    --muted-foreground: oklch(0.708 0 0);
    --border: oklch(1 0 0 / 10%);
    /* ...full dark token set... */
  }
}

@theme inline {
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-foreground);
  --color-border: var(--border);
  --radius-sm: calc(var(--radius) - 4px);
  --radius-md: calc(var(--radius) - 2px);
  --radius-lg: var(--radius);
  /* ...full mapping... */
}

@layer base {
  * { @apply border-border; }
  body { @apply bg-background text-foreground; }
}
```

**Dark mode inline script in `index.html` head (prevents flash of unstyled content):**

```html
<!-- index.html -->
<head>
  <script>
    (function () {
      const stored = localStorage.getItem("theme")
      const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches
      if (stored === "dark" || (!stored && prefersDark)) {
        document.documentElement.classList.add("dark")
      }
    })()
  </script>
</head>
```

**React `useTheme` hook:**

```typescript
// src/hooks/useTheme.ts
import { useEffect, useState } from "react"

type Theme = "dark" | "light" | "system"

export function useTheme() {
  const [theme, setTheme] = useState<Theme>(() => {
    return (localStorage.getItem("theme") as Theme) ?? "system"
  })

  useEffect(() => {
    const root = document.documentElement
    const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches

    if (theme === "dark" || (theme === "system" && prefersDark)) {
      root.classList.add("dark")
    } else {
      root.classList.remove("dark")
    }

    if (theme === "system") {
      localStorage.removeItem("theme")
    } else {
      localStorage.setItem("theme", theme)
    }
  }, [theme])

  return { theme, setTheme }
}
```

---

### Pattern 7: Zod API Response Envelope

**What:** Zod schema for the backend's standard response format.
**When to use:** Wrap all API calls — validate at the boundary.

```typescript
// src/api/types.ts
import { z } from "zod"

// Backend standard envelope: { success, data, error?, meta? }
export const ApiResponseSchema = <T extends z.ZodTypeAny>(dataSchema: T) =>
  z.object({
    success: z.boolean(),
    data: dataSchema,
    error: z.string().optional(),
    meta: z
      .object({
        total: z.number(),
        limit: z.number(),
        offset: z.number(),
      })
      .optional(),
  })

// Auth login response
export const LoginResponseSchema = ApiResponseSchema(
  z.object({
    access_token: z.string(),
    refresh_token: z.string(),
    expires_in: z.number(),
    user: z.object({
      id: z.string(),
      email: z.string(),
      role: z.enum(["superadmin", "admin", "teknisi"]),
      full_name: z.string(),
    }),
  })
)

export type LoginResponse = z.infer<typeof LoginResponseSchema>
```

---

### Anti-Patterns to Avoid

- **Plugin order wrong in vite.config.ts:** `TanStackRouterVite` MUST come before `react()`. Reversed order breaks route tree generation.
- **Tailwind v4 with `tailwind.config.js`:** Do not create a config file. v4 is CSS-first. Any `tailwind.config.js` in the project root will cause v4 to fall back to v3 mode.
- **`darkMode: 'class'` in tailwind.config:** This key does not exist in v4. Use `@custom-variant dark` in CSS.
- **Axios interceptor loop:** The refresh call inside the interceptor MUST use plain `axios.post(...)`, not the `adminClient` instance — otherwise the interceptor will intercept its own refresh call, causing infinite loops.
- **Multiple refresh calls:** Without the `isRefreshing` flag + queue, simultaneous 401s will each try to refresh, invalidating each other's tokens.
- **Shadcn/UI v3 setup with v4 Tailwind:** Running `npx shadcn@latest init` on a project with `@tailwindcss/vite` (v4) auto-detects v4 and generates the correct `@theme inline` CSS. Do NOT follow v3 setup docs.
- **`useStore` inside Axios interceptors:** Axios interceptors run outside React. Use the vanilla getter pattern `useStore.getState().adminAccessToken` instead of a hook.
- **Route file naming confusion:** `_admin/route.tsx` is the layout. `_admin/dashboard.tsx` is a child at `/dashboard`. If you accidentally name a file `_admin.tsx` (flat, not a directory), it becomes a different route type.

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Token refresh + queue | Custom fetch wrapper with retry | Axios interceptors (Pattern 4 above) | Race condition handling is non-trivial; queue + `isRefreshing` pattern is battle-tested |
| Form validation | Manual input checks | Zod + TanStack Form | Nested error paths, async validation, array fields are complex |
| Data tables with pagination | Custom table | TanStack Table v8 | Virtualization, sorting, filtering, column pinning — all edge cases |
| State persistence | `localStorage.setItem` calls scattered | Zustand `persist` middleware | Handles hydration timing, version migration, partial state |
| Toast notifications | Custom toast system | Sonner | Positioning, stacking, dismissal, accessibility, promise toasts |
| Date formatting | Manual `new Date().toLocaleDateString()` | date-fns | Locale, timezone, relative time, DST — all footguns |
| Icon SVGs | Copy-paste SVGs | lucide-react (bundled with Shadcn) | Tree-shaking, consistent sizing, accessibility |
| Dark mode with FOUC | CSS-only or React state | Inline script in `<head>` (Pattern 6 above) | Script runs before paint, eliminates flash |

**Key insight:** Most custom solutions in this domain require handling edge cases (race conditions, timezone issues, accessibility) that dedicated libraries already solve.

---

## Common Pitfalls

### Pitfall 1: FOUC (Flash of Unstyled Content) on Dark Mode
**What goes wrong:** Page renders in light mode for ~100ms then switches to dark — visible flash.
**Why it happens:** React state initializes after first paint; `useEffect` runs too late.
**How to avoid:** Add the inline `<script>` in `index.html` `<head>` BEFORE any stylesheet link (Pattern 6). The script runs synchronously before paint.
**Warning signs:** Visible white flash on page load when OS is set to dark mode.

### Pitfall 2: TanStack Router Plugin Order
**What goes wrong:** Routes are not generated; build fails or `routeTree.gen.ts` is not updated.
**Why it happens:** `@tanstack/router-plugin/vite` must run before `@vitejs/plugin-react` to scan and generate the route tree.
**How to avoid:** Always place `TanStackRouterVite(...)` as the FIRST entry in `plugins: []`.
**Warning signs:** `routeTree.gen.ts` not updating on file save; TypeScript errors about missing route types.

### Pitfall 3: Pathless Layout Route — Wrong File Placement
**What goes wrong:** `_admin/dashboard.tsx` creates route at `/_admin/dashboard` instead of `/dashboard`.
**Why it happens:** If the directory itself has a path segment (no `_` prefix), children inherit it.
**How to avoid:** The directory must be named `_admin/` (underscore prefix) for it to be pathless. The `route.tsx` inside defines the layout. Children of `_admin/` render at their own path without the `_admin` prefix.
**Warning signs:** URL has `/_admin/` in it during navigation.

### Pitfall 4: Axios Instance Intercepting Its Own Refresh Call
**What goes wrong:** The refresh token endpoint returns 401 (expired refresh token), which triggers another refresh attempt — infinite loop.
**Why it happens:** Using `adminClient` for the refresh call means the response interceptor intercepts the refresh response.
**How to avoid:** Use plain `import axios from "axios"` for the refresh POST, not the `adminClient` instance (Pattern 4). The plain axios instance has no interceptors.
**Warning signs:** Network tab shows repeated calls to `/auth/refresh` in quick succession.

### Pitfall 5: Zustand v5 Persist on Store Creation
**What goes wrong:** Persisted state is not loaded on first render — `accessToken` is `null` even though localStorage has a value.
**Why it happens:** Zustand v5 changed persist middleware behavior: items are no longer stored at creation time, and rehydration is async.
**How to avoid:** Use `onRehydrateStorage` callback to set `isHydrated: true` and gate rendering behind that flag. Or use `useStore.persist.hasHydrated()`.
**Warning signs:** Auth state is null on page reload even with valid localStorage values; user is redirected to login on every refresh.

### Pitfall 6: Shadcn/UI `components.json` baseUrl Mismatch
**What goes wrong:** `npx shadcn add button` adds components to wrong path, or imports fail.
**Why it happens:** `components.json` `aliases.components` must match the `@/` alias in `tsconfig`.
**How to avoid:** Run `npx shadcn@latest init` AFTER configuring `tsconfig.json` `@/*` alias. Let the CLI detect the correct paths.
**Warning signs:** TypeScript import errors on generated Shadcn components.

### Pitfall 7: React 19 + `@vitejs/plugin-react-swc`
**What goes wrong:** Build warnings or incompatibilities with React 19.
**Why it happens:** The `vite-plugin-react-swc` repository was archived in July 2025. The SWC functionality migrated into the main `@vitejs/plugin-react` package.
**How to avoid:** Use `@vitejs/plugin-react` (v6.0.1) — it now includes SWC-based transforms. Do NOT add the old `@vitejs/plugin-react-swc` package separately.
**Warning signs:** Dependency resolution warnings; stale security advisories.

---

## Code Examples

### Login mutation with TanStack Query + Zod validation

```typescript
// src/api/auth.ts
import { adminClient } from "@/lib/axios/admin-client"
import { LoginResponseSchema } from "./types"
import type { LoginResponse } from "./types"

export async function adminLogin(email: string, password: string): Promise<LoginResponse> {
  const { data } = await adminClient.post("/auth/login", { email, password })
  return LoginResponseSchema.parse(data)
}
```

```typescript
// Usage in login page
import { useMutation } from "@tanstack/react-query"
import { useNavigate } from "@tanstack/react-router"
import { adminLogin } from "@/api/auth"
import { useStore } from "@/store"
import { toast } from "sonner"

function AdminLoginPage() {
  const navigate = useNavigate()
  const { adminSetTokens, adminSetUser } = useStore()

  const mutation = useMutation({
    mutationFn: ({ email, password }: { email: string; password: string }) =>
      adminLogin(email, password),
    onSuccess: (data) => {
      adminSetTokens(data.data.access_token, data.data.refresh_token)
      adminSetUser(data.data.user)
      navigate({ to: "/dashboard" })
      toast.success("Login berhasil")
    },
    onError: (err) => {
      toast.error("Email atau password salah")
    },
  })

  // ...form JSX
}
```

### Auth context hook (for router context)

```typescript
// src/hooks/useAuth.ts
import { useStore } from "@/store"

export function useAdminAuthContext() {
  const isAuthenticated = useStore((s) => s.adminIsAuthenticated)
  const role = useStore((s) => s.adminUser?.role ?? null)
  const accessToken = useStore((s) => s.adminAccessToken)
  return { isAuthenticated, role, accessToken }
}

export function useAgentAuthContext() {
  const isAuthenticated = useStore((s) => s.agentIsAuthenticated)
  const accessToken = useStore((s) => s.agentAccessToken)
  return { isAuthenticated, accessToken }
}

export function useCustomerAuthContext() {
  const isAuthenticated = useStore((s) => s.customerIsAuthenticated)
  const accessToken = useStore((s) => s.customerAccessToken)
  return { isAuthenticated, accessToken }
}
```

### Shadcn/UI `cn` utility

```typescript
// src/lib/utils.ts  (generated by shadcn init)
import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `tailwind.config.js` with `darkMode: 'class'` | `@custom-variant dark` in CSS | Tailwind v4 (2025) | No JS config file needed |
| `@tailwind base/components/utilities` | `@import "tailwindcss"` | Tailwind v4 | Single import replaces three directives |
| HSL color values in tokens | OKLCH color values | Shadcn/UI 2025 refresh | Better perceptual uniformity in dark mode |
| `tailwindcss-animate` | `tw-animate-css` | Shadcn/UI v4 migration | Old package not maintained for v4 |
| `toast` from shadcn/ui | `sonner` | Shadcn/UI deprecation 2025 | `toast` component officially deprecated |
| `React.forwardRef` in Shadcn components | Standard function + `React.ComponentProps` | React 19 / Shadcn refresh | React 19 supports ref as prop directly |
| `@vitejs/plugin-react-swc` (separate) | `@vitejs/plugin-react` v6 (unified) | July 2025 | SWC repo archived; unified into main plugin |
| `cacheTime` in TanStack Query | `gcTime` | TanStack Query v5 | Renamed for clarity |
| `useQuery({ queryKey, queryFn, cacheTime })` | `useQuery({ queryKey, queryFn, gcTime })` | TanStack Query v5 | Breaking rename |

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Node.js | Build tooling | Yes | v22.22.1 | — |
| npm | Package management | Yes | 10.9.4 | — |
| pnpm | Faster installs | No | — | Use npm (project uses npm) |
| Go backend | API proxy target | Assumed running on :8080 | — | Dev can run `go run` separately |

**Missing dependencies with no fallback:**
- None that block execution. The Go backend must be running for API calls to work, but Vite dev server and frontend can scaffold/build independently.

---

## Open Questions

1. **Admin base path — root `/` vs `/admin`**
   - What we know: REQUIREMENTS.md says Admin at `/`, Agent at `/agent`, Customer at `/customer`
   - What's unclear: The `_admin` pathless layout means admin routes live at `/dashboard`, `/customers`, etc. — top-level. This is correct per requirements but means `/` is NOT a marketing landing page; it immediately redirects to `/dashboard` or `/login`.
   - Recommendation: Create `src/routes/index.tsx` that redirects to `/dashboard` (if authenticated) or `/login` (if not). Confirm this is the intended behavior.

2. **Refresh token storage**
   - What we know: Project uses localStorage via Zustand persist (acceptable for ISP internal dashboard)
   - What's unclear: The backend sets `refresh_token` in the response body, not as httpOnly cookie. If backend adds httpOnly cookie support in future, interceptor pattern changes significantly.
   - Recommendation: Use localStorage for now. Document that migrating to httpOnly cookies requires backend change + removing Zustand persist for refresh token.

3. **Multiple browser tabs + token refresh**
   - What we know: The queue pattern handles multiple concurrent requests in ONE tab
   - What's unclear: Two tabs simultaneously detecting a 401 will each call refresh — the second call may fail with an already-used refresh token
   - Recommendation: Add a `storage` event listener on `adminAccessToken` key to sync new tokens across tabs. The leader-election pattern (`broadcast-channel` library) solves this fully if it becomes an issue.

---

## Sources

### Primary (HIGH confidence)
- [TanStack Router v1 Authenticated Routes](https://tanstack.com/router/v1/docs/framework/react/guide/authenticated-routes) — beforeLoad, redirect, context patterns
- [TanStack Router Vite Installation](https://tanstack.com/router/v1/docs/framework/react/installation/with-vite) — plugin setup, autoCodeSplitting
- [Shadcn/UI Tailwind v4 Guide](https://ui.shadcn.com/docs/tailwind-v4) — CSS variables, OKLCH, migration notes
- [Shadcn/UI Vite Installation](https://ui.shadcn.com/docs/installation/vite) — complete vite.config.ts, tsconfig setup
- [Tailwind CSS v4 Dark Mode](https://tailwindcss.com/docs/dark-mode) — @custom-variant, class strategy, localStorage pattern
- npm registry — verified all package versions as of 2026-03-30

### Secondary (MEDIUM confidence)
- [DEV: Custom Layout for Route Group in TanStack Router](https://dev.to/xb16/custom-layout-for-specific-route-group-in-tanstack-router-solution-2ndp) — pathless layout route structure confirmed
- [leonardomontini.dev: TanStack Router auth guard](https://leonardomontini.dev/tanstack-router-guard/) — beforeLoad + context patterns
- [Medium: Token Refresh with Axios Interceptors](https://medium.com/@velja/token-refresh-with-axios-interceptors-for-a-seamless-authentication-experience-854b06064bde) — queue pattern TypeScript implementation
- [doichevkostia.dev: Authentication store with Zustand](https://doichevkostia.dev/blog/authentication-store-with-zustand/) — vanilla store + persist + hydration
- [Medium: Dark/light mode Tailwind v4](https://medium.com/@moazm942/how-to-create-dark-light-mode-toggle-with-tailwind-v4-400eb2a9cf39) — inline script + React hook

### Tertiary (LOW confidence)
- GitHub Issue #252 vite-plugin-react-swc React 19 — SWC repo archived July 2025 (needs validation of @vitejs/plugin-react v6 feature parity)

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all versions verified against npm registry 2026-03-30
- Architecture: HIGH — patterns sourced from official TanStack Router docs (file-based routing, beforeLoad, router context)
- Tailwind v4 + Shadcn/UI setup: HIGH — official shadcn docs confirmed v4 support
- Axios interceptor pattern: HIGH — well-established pattern across multiple authoritative sources
- Zustand v5 slice + persist: HIGH — official zustand docs + confirmed v5.0.12 changelog fix for persist
- Pitfalls: MEDIUM — most verified, FOUC and plugin order confirmed; cross-tab refresh token (LOW, edge case)

**Research date:** 2026-03-30
**Valid until:** 2026-04-30 (TanStack Router and shadcn/ui move fast; re-verify before implementing if delayed)
