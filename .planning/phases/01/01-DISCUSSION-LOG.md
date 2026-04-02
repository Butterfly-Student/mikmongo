# Phase 1: Auth & API Foundation - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-02
**Phase:** 01-auth-api-foundation
**Areas discussed:** Login page UX

---

## Login Page UX

| Option | Description | Selected |
|--------|-------------|----------|
| Single centered card | Centered card with logo, email/password fields, submit button. Clean, minimal. Existing code already does this. | ✓ |
| Split layout with branding | Split layout — branding/illustration on left, form on right. More visually striking but more work. | |
| Custom design | Something specific in mind. | |

**User's choice:** Single centered card
**Notes:** Consistent with existing code, works for all three portals.

---

| Option | Description | Selected |
|--------|-------------|----------|
| Per-field errors + toast | Show errors below each field and a toast for API errors. Precise feedback. Existing code does this. | ✓ |
| Toast-only errors | Only show a toast at top-right for all errors. Simpler but less precise. | |
| Inline-only errors | Show errors inline below the form submit button. No toast. | |

**User's choice:** Per-field errors + toast
**Notes:** Matches existing code pattern, gives precise feedback per field.

---

| Option | Description | Selected |
|--------|-------------|----------|
| Spinner + disabled form | Button shows spinner + 'Signing in...' text, entire form disabled during request. Standard pattern. Existing code does this. | ✓ |
| Spinner on button only | Button shows spinner only, form stays interactive. Unusual. | |
| Full-page loading overlay | Overlay with full-page loader. Blocks interaction entirely. | |

**User's choice:** Spinner + disabled form
**Notes:** Standard pattern, matches existing code.

---

| Option | Description | Selected |
|--------|-------------|----------|
| Indonesian | Use Indonesian text for labels, errors, and buttons. Existing code already uses Indonesian. | ✓ |
| English | Use English text throughout. | |
| Mixed language | Labels/buttons in English, errors in Indonesian. Unusual. | |

**User's choice:** Indonesian
**Notes:** Existing code already uses Indonesian strings. Consistent.

---

| Option | Description | Selected |
|--------|-------------|----------|
| M logo + portal title | M logo + 'MikMongo' title + portal subtitle. Clean and simple. Existing code pattern. | ✓ |
| Custom logo image | Custom logo image (user would provide asset). | |
| No logo | No logo — just card title and subtitle. Minimalist. | |

**User's choice:** M logo + portal title
**Notes:** Matches existing code. Portal subtitle changes per portal.

---

| Option | Description | Selected |
|--------|-------------|----------|
| Include remember/forgot | Include 'Remember me' checkbox and 'Forgot password?' link. Standard login UX. | ✓ |
| Minimal — just email + password | Just email, password, submit. Minimal. | |
| Remember me only | 'Remember me' toggle only, no forgot password link. | |

**User's choice:** Include remember/forgot
**Notes:** UI completeness even without backend forgot-password endpoint.

---

| Option | Description | Selected |
|--------|-------------|----------|
| Redirect to dashboard | After login, go to dashboard. Simple and expected. | ✓ |
| Redirect to original URL | Redirect back to URL the user was trying to access. Better for deep-linking. | |
| Smart redirect | Use original URL if protected route redirect, otherwise dashboard. | |

**User's choice:** Redirect to dashboard
**Notes:** Simplest, consistent across portals. No deep-link preservation.

---

## Claude's Discretion

- Token storage mechanism details (localStorage key name, storage structure)
- Exact refresh timing (proactive vs reactive-only)
- Logout behavior edge cases (multiple tabs, network failures)
- Change password page layout and flow
- Customer and agent portal login page visual differences from admin

## Deferred Ideas

- "Remember me" — checkbox present but no backend support for extended sessions
- "Forgot password?" — link present but no endpoint; could show toast or static page
- Proactive token refresh — deferred to implementation planning
