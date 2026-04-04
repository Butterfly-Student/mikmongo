---
status: partial
phase: 03-customers-routers-subscriptions
source: [03-VERIFICATION.md]
started: 2026-04-04T13:00:00Z
updated: 2026-04-04T13:00:00Z
---

## Current Test

[awaiting human testing]

## Tests

### 1. Full Sidebar Navigation Flow
expected: All three nav items (Customers, Routers, Subscriptions) are clickable and route to the correct pages without errors or blank screens
result: [pending]

### 2. Router Edit Round-Trip
expected: EditRouterDialog opens pre-populated with current values; saving updates the table; empty password field means password is not changed
result: [pending]

### 3. Customer Portal Subscriptions at /customer/subscriptions
expected: Subscription cards render with status badges, IP, profile, expiry; WifiOff empty state when no subscriptions; requires logged-in customer session
result: [pending]

### 4. Registration Approval with Dependent Profile Select
expected: Router dropdown appears; profiles load only after a router is selected; approval wires profile to customer
result: [pending]

## Summary

total: 4
passed: 0
issues: 0
pending: 4
skipped: 0
blocked: 0

## Gaps
