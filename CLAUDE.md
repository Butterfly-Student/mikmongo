<!-- GSD:project-start source:PROJECT.md -->
## Project

**MikMongo ISP Management Dashboard**

A full-featured ISP management dashboard built on the shadcn-admin template, implementing the MikMongo API (80+ endpoints). It serves three portals — Admin, Customer self-service, and Agent self-service — with complete MikroTik router management, real-time monitoring via WebSockets, billing, payments (Midtrans/Xendit), hotspot voucher sales, and business reports.

**Core Value:** Admin can manage their entire ISP operation from one dashboard: customers, routers, subscriptions, billing, and monitor MikroTik devices in real-time.

### Constraints

- **Tech stack**: Must use existing template stack (React 19, TypeScript, TanStack, Tailwind, shadcn/ui, Axios, Zustand, Zod)
- **Template structure**: Must follow the feature-based directory convention; no structural rewrites
- **API contract**: Must match `docs/openapi.docs.yml` exactly — schemas, endpoints, auth schemes
- **Immutability**: Immutable data patterns throughout (no mutation of existing objects)
- **Real-time**: WebSocket endpoints must provide live updates without page refresh
<!-- GSD:project-end -->

<!-- GSD:stack-start source:STACK.md -->
## Technology Stack

Technology stack not yet documented. Will populate after codebase mapping or first phase.
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd:quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd:debug` for investigation and bug fixing
- `/gsd:execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd:profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
