---
phase: 03-customers-routers-subscriptions
plan: 02
type: execute
wave: 1
depends_on: []
files_modified: ["website/src/features/customers/data/schema.ts", "website/src/features/customers/data/columns.tsx", "website/src/features/customers/components/customer-table.tsx", "website/src/features/customers/components/create-customer-dialog.tsx", "website/src/features/customers/index.tsx", "website/src/routes/_admin/customers.tsx"]
autonomous: true
requirements: ["CUST-01", "CUST-02", "CUST-03", "CUST-04", "CUST-05", "CUST-06", "CUST-07"]
must_haves:
  truths:
    - "Admin can view customer list and pipeline"
    - "Admin can create new customers"
    - "Admin can approve or reject pending customer registrations"
  artifacts:
    - "website/src/features/customers/index.tsx"
    - "website/src/routes/_admin/customers.tsx"
  key_links:
    - "website/src/features/customers/index.tsx -> /api/v1/customers"
---

<objective>
Implement Customer CRUD and Registration Pipeline.

Purpose: Admin needs to manage the customer lifecycle, viewing complete lists of customers and handling the flow of pending registrations (approve/reject).
Output: Customer Data Table with status pipelines and forms.
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
  <name>Task 1: Define Customer Schema and Data Layer</name>
  <files>website/src/features/customers/data/schema.ts</files>
  <action>Create Zod schemas matching the OpenAPI definitions for CustomerResponse, CreateCustomerRequest. Be sure to include fields relevant to activation states (pending/approved/rejected).</action>
  <verify>
    <automated>npx tsc --noEmit</automated>
  </verify>
  <done>Zod schema matches /api/v1/customers constraints.</done>
  <read_first>
    - website/src/features/customers/data/schema.ts
    - docs/openapi.docs.yml
  </read_first>
</task>

<task type="auto">
  <name>Task 2: Build Customer Table and Registration Pipeline UX</name>
  <files>website/src/features/customers/data/columns.tsx, website/src/features/customers/components/customer-table.tsx</files>
  <action>Implement the Customer data table using TanStack Table (re-use standard table layout per UI-SPEC.md and D-01). Add "Approve" and "Reject" buttons directly into the row actions for customers in "pending" state. Add activate/deactivate account actions.</action>
  <verify>
    <automated>npx tsc --noEmit</automated>
  </verify>
  <done>Customer columns and table components handle the registration pipeline correctly.</done>
  <read_first>
    - website/src/features/customers/components/customer-table.tsx
    - .planning/phases/03-customers-routers-subscriptions/03-CONTEXT.md
  </read_first>
</task>

<task type="auto">
  <name>Task 3: Build Customer Views and Routing</name>
  <files>website/src/features/customers/components/create-customer-dialog.tsx, website/src/features/customers/index.tsx, website/src/routes/_admin/customers.tsx</files>
  <action>Assemble the main Customers page. Build the Create Customer dialog. Wire up hooks to /api/v1/customers endpoints for listing and creating. Ensure `_admin/customers.tsx` is properly structured within TanStack Router.</action>
  <verify>
    <automated>npx tsc --noEmit</automated>
  </verify>
  <done>Customers index exports a fully working module, added to routing.</done>
  <read_first>
    - website/src/features/customers/index.tsx
    - website/src/routes/_admin/customers.tsx
  </read_first>
</task>

</tasks>

<verification>
Automated checks using `npx tsc --noEmit`. Visual checks ensure standard list layout and status buttons as per D-01.
</verification>

<success_criteria>
Customer lifecycle management is implemented correctly, handling pending and active customers.
</success_criteria>

<output>
After completion, create `.planning/phases/03-customers-routers-subscriptions/03-02-SUMMARY.md`
</output>
