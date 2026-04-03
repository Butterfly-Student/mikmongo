---
phase: 03-customers-routers-subscriptions
plan: 01
type: execute
wave: 1
depends_on: []
files_modified: ["website/src/features/routers/data/schema.ts", "website/src/features/routers/data/columns.tsx", "website/src/features/routers/components/router-table.tsx", "website/src/features/routers/components/create-router-dialog.tsx", "website/src/features/routers/index.tsx", "website/src/features/profiles/data/schema.ts", "website/src/features/profiles/data/columns.tsx", "website/src/features/profiles/components/profile-table.tsx", "website/src/features/profiles/index.tsx", "website/src/routes/_admin/routers.tsx"]
autonomous: true
requirements: ["RTR-01", "RTR-02", "RTR-03", "RTR-04", "RTR-05", "RTR-06", "BW-01", "BW-02", "BW-03"]
must_haves:
  truths:
    - "Admin can view, add, and sync MikroTik routers"
    - "Admin can test router connections"
    - "Admin can view and manage bandwidth profiles tied to a router"
  artifacts:
    - "website/src/features/routers/index.tsx"
    - "website/src/features/profiles/index.tsx"
  key_links:
    - "website/src/features/routers/index.tsx -> /api/v1/routers"
---

<objective>
Implement MikroTik Router CRUD, synchronization, connection testing, and Bandwidth Profile CRUD.

Purpose: Admin needs to register their MikroTik routers and define available bandwidth profiles (plans) before any customers can be provisioned.
Output: Router Management UI and Profile Management UI fully wired to the backend API.
</objective>

<execution_context>
@~/.gemini/antigravity/get-shit-done/workflows/execute-plan.md
@~/.gemini/antigravity/get-shit-done/templates/summary.md
</execution_context>

<context>
@.planning/PROJECT.md
@.planning/ROADMAP.md
@.planning/REQUIREMENTS.md
@.planning/phases/03-customers-routers-subscriptions/03-CONTEXT.md
@.planning/phases/03-customers-routers-subscriptions/03-UI-SPEC.md
</context>

<tasks>

<task type="auto">
  <name>Task 1: Define Schemas and Data Layer</name>
  <files>website/src/features/routers/data/schema.ts, website/src/features/profiles/data/schema.ts</files>
  <action>Create Zod schemas matching the OpenAPI definitions for RouterResponse, CreateRouterRequest, ProfileResponse, and CreateProfileRequest. Do not mutate existing state. Follow the pattern established in earlier phases (e.g., users feature).</action>
  <verify>
    <automated>npx tsc --noEmit</automated>
  </verify>
  <done>Zod schemas export exact types for routers and bandwidth profiles.</done>
  <read_first>
    - website/src/features/routers/data/schema.ts
    - website/src/features/profiles/data/schema.ts
    - docs/openapi.docs.yml
  </read_first>
</task>

<task type="auto">
  <name>Task 2: Build Router Management UI</name>
  <files>website/src/features/routers/data/columns.tsx, website/src/features/routers/components/router-table.tsx, website/src/features/routers/components/create-router-dialog.tsx, website/src/features/routers/index.tsx</files>
  <action>Implement the Data Table for Routers using TanStack Table. Include standard search and faceted filters. Implement Create Dialog using react-hook-form + Zod. Add actions for Sync (POST /api/v1/routers/{id}/sync) and Test Connection (POST /api/v1/routers/{id}/test-connection) using blocking loading states (per D-04 context). Style using 03-UI-SPEC.md spacing and colors.</action>
  <verify>
    <automated>npx tsc --noEmit</automated>
  </verify>
  <done>Router list, create, sync, and test connection actions are implemented and wired to TanStack Query hooks hitting /api/v1/routers.</done>
  <read_first>
    - website/src/features/routers/index.tsx
    - .planning/phases/03-customers-routers-subscriptions/03-UI-SPEC.md
  </read_first>
</task>

<task type="auto">
  <name>Task 3: Build Bandwidth Profile UI</name>
  <files>website/src/features/profiles/data/columns.tsx, website/src/features/profiles/components/profile-table.tsx, website/src/features/profiles/index.tsx, website/src/routes/_admin/routers.tsx</files>
  <action>Implement Bandwidth Profile management table. Profiles are nested under routers, so this UI should either accept a routerId or be placed intuitively into the router details view. Add to the TanStack router path `_admin/routers.tsx`. Wire to /api/v1/routers/{router_id}/bandwidth-profiles CRUD endpoints.</action>
  <verify>
    <automated>npx tsc --noEmit</automated>
  </verify>
  <done>Profile list mapped to specific router contexts, complete with create and delete actions.</done>
  <read_first>
    - website/src/features/profiles/index.tsx
    - website/src/routes/_admin/routers.tsx
  </read_first>
</task>

</tasks>

<verification>
Automated checks using `npx tsc --noEmit` and matching visual contracts.
</verification>

<success_criteria>
Router and Bandwidth Profile components are fully implemented and type-safe.
</success_criteria>

<output>
After completion, create `.planning/phases/03-customers-routers-subscriptions/03-01-SUMMARY.md`
</output>
