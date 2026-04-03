# Phase 2: Layout, Dashboard & Users - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-03
**Phase:** 02-layout-dashboard-users
**Areas discussed:** Router store & persistence, Template cleanup scope

---

## Router store & persistence

### Where should the active router state live?

| Option | Description | Selected |
|--------|-------------|----------|
| Separate store | New `router-store.ts` alongside auth-store. Cleaner separation of concerns. | ✓ |
| Slice in auth store | Add RouterSlice to existing auth-store. All persisted state in one place. | |

**User's choice:** Separate store (Recommended)
**Notes:** Router selection is a different domain than auth. Separate store keeps concerns clean.

### What should the dashboard show when no router is selected?

| Option | Description | Selected |
|--------|-------------|----------|
| Widgets work, ping shows '-- ms' | All widgets load from /reports/summary (doesn't need router). Ping shows "-- ms". Router health shows empty state. | ✓ |
| Block dashboard, show prompt | Show "Select a router to get started" instead of widgets. | |

**User's choice:** Widgets work, ping shows '-- ms' (Recommended)
**Notes:** Dashboard data (customers, subscriptions, revenue) doesn't depend on a specific router.

### Should the selected router sync from the API on page load?

| Option | Description | Selected |
|--------|-------------|----------|
| Sync from API on load | Fetch from `GET /routers/selected` and update store. Keeps frontend in sync. | ✓ |
| localStorage only | Only use localStorage value. Simpler but can get out of sync. | |

**User's choice:** Sync from API on load (Recommended)
**Notes:** Prevents stale state if router was deselected in another admin session.

### What happens when admin switches routers?

| Option | Description | Selected |
|--------|-------------|----------|
| Immediate refresh | Widgets update immediately. Ping reconnects WebSocket. KPIs refresh. | ✓ |
| Loading transition | Show loading skeleton while switching routers. | |

**User's choice:** Immediate refresh (Recommended)
**Notes:** Cleaner UX — no artificial delay.

---

## Template cleanup scope

### What to do with template demo pages?

| Option | Description | Selected |
|--------|-------------|----------|
| Keep files, hide from nav | Leave files but remove from sidebar nav. Less risky. | ✓ |
| Remove template pages now | Delete route files and components for apps, chats, tasks, help-center. | |

**User's choice:** Keep files, hide from nav (Recommended)
**Notes:** Can clean up in a future dedicated cleanup phase.

### Should sidebar show future nav items?

| Option | Description | Selected |
|--------|-------------|----------|
| Disabled items with 'Coming soon' | All ISP groups shown. Unimplemented items greyed out with tooltip. | ✓ |
| Only show active groups | Only groups with active routes. Sidebar grows as phases complete. | |

**User's choice:** Disabled items with 'Coming soon' (Recommended)
**Notes:** Admin can see what's planned even if not yet implemented.

---

## Skipped Gray Areas (not selected for discussion)

- **User edit & missing API** — User chose not to discuss. Claude's discretion at implementation time.
- **Activity feed data source** — User chose not to discuss. Claude's discretion at implementation time.

## Claude's Discretion

- User edit functionality (no PUT endpoint — skip, workaround, or defer)
- Activity feed data source (combined data vs placeholder)
- Profile dropdown update to use auth store
- Router store exact structure
- WebSocket reconnection strategy
- Loading skeleton designs

## Deferred Ideas

- Template demo page file deletion — keep for now, clean up later
- Profile dropdown auth store integration — cosmetic
