---
phase: 03-customers-routers-subscriptions
plan: 03
type: execute
wave: 2
depends_on: ["01", "02"]
files_modified: ["website/src/features/subscriptions/data/schema.ts", "website/src/features/subscriptions/data/columns.tsx", "website/src/features/subscriptions/components/subscription-table.tsx", "website/src/features/subscriptions/components/create-subscription-dialog.tsx", "website/src/features/subscriptions/index.tsx", "website/src/routes/_admin/subscriptions.tsx", "website/src/components/ui/confirm-action-dialog.tsx"]
autonomous: true
requirements: ["SUB-01", "SUB-02", "SUB-03", "SUB-04", "SUB-05"]
must_haves:
  truths:
    - "Admin can view all subscriptions on a router"
    - "Admin can assign profiles to customers"
    - "Admin can safely suspend, restore, isolate, or terminate subscriptions via explicit confirmation"
  artifacts:
    - "website/src/features/subscriptions/index.tsx"
    - "website/src/components/ui/confirm-action-dialog.tsx"
  key_links:
    - "website/src/features/subscriptions/index.tsx -> /api/v1/routers/{router_id}/subscriptions"
---

<objective>
Implement Subscription Assignment Flow and Destructive Actions.

Purpose: Connect Customers, Routers, and Bandwidth Profiles together via subscriptions. Provide a safe interface for modifying subscription states.
Output: Global Subscriptions view, assignment dropdowns, and destructive confirmation dialogs.
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
  <name>Task 1: Build Shared Confirmation Dialog</name>
  <files>website/src/components/ui/confirm-action-dialog.tsx</files>
  <action>Create a reusable `ConfirmActionDialog` component based on shadcn `AlertDialog`. It should accept title, description, destructive flag, and confirmation handlers. This fulfills D-03 from CONTEXT.md requiring explicit dialogs for destructive actions.</action>
  <verify>
    <automated>npx tsc --noEmit</automated>
  </verify>
  <done>Confirm action dialog component is reusable and correctly formatted per UI spec.</done>
  <read_first>
    - website/src/components/ui/alert-dialog.tsx
    - .planning/phases/03-customers-routers-subscriptions/03-CONTEXT.md
  </read_first>
</task>

<task type="auto">
  <name>Task 2: Global Subscriptions View</name>
  <files>website/src/features/subscriptions/data/schema.ts, website/src/features/subscriptions/data/columns.tsx, website/src/features/subscriptions/components/subscription-table.tsx, website/src/features/subscriptions/index.tsx, website/src/routes/_admin/subscriptions.tsx</files>
  <action>Create the global Subscriptions view (D-02 from CONTEXT.md). Define Zod schema for SubscriptionResponse. Create a data table that lists assigned subscriptions. Note that API endpoints are router-scoped (`/api/v1/routers/{router_id}/subscriptions`), so the table must consume the active router from the existing router context/ Zustand store to fetch data.</action>
  <verify>
    <automated>npx tsc --noEmit</automated>
  </verify>
  <done>Subscription table renders data accurately tied to the active router.</done>
  <read_first>
    - website/src/features/subscriptions/index.tsx
    - docs/openapi.docs.yml
  </read_first>
</task>

<task type="auto">
  <name>Task 3: Subscription Actions & Assignment</name>
  <files>website/src/features/subscriptions/components/create-subscription-dialog.tsx, website/src/features/subscriptions/data/columns.tsx</files>
  <action>Build `CreateSubscriptionDialog` utilizing `<Select>` components to choose a Customer and a Bandwidth Profile (D-02 constraint). Wire up the row actions for Suspend, Restore, Isolate, and Terminate to conditionally render the new `ConfirmActionDialog` before firing the mutation endpoints.</action>
  <verify>
    <automated>npx tsc --noEmit</automated>
  </verify>
  <done>Create form and explicit row action confirmations are implemented.</done>
  <read_first>
    - website/src/features/subscriptions/components/create-subscription-dialog.tsx
    - .planning/phases/03-customers-routers-subscriptions/03-UI-SPEC.md
  </read_first>
</task>

</tasks>

<verification>
No errors in `npx tsc --noEmit`. Visual validation ensures explicit Dialog popups on suspend actions.
</verification>

<success_criteria>
Subscriptions correctly link routers, profiles, and customers. Explicit confirmations exist.
</success_criteria>

<output>
After completion, create `.planning/phases/03-customers-routers-subscriptions/03-03-SUMMARY.md`
</output>
