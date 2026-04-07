# Phase 4: Billing & Payments — UI Design Contract

**Generated:** 2026-04-04
**Status:** Ready for planning

## 1. Design System

### Typography
- **Page Titles:** `text-2xl font-bold tracking-tight`
- **Section / Sheet Headers:** `text-lg font-semibold`
- **Body:** `text-sm font-normal text-foreground`
- **Muted / Labels:** `text-sm text-muted-foreground`
- **Badges:** `text-xs font-semibold uppercase tracking-wider`

### Color Tokens
- **Page Background:** `bg-background`
- **Card / Section Background:** `bg-card`
- **Table Header Background:** `bg-muted`
- **Borders:** `border-border`
- **Primary Actions:** `bg-primary text-primary-foreground`
- **Destructive Actions (Reject / Refund):** `bg-destructive text-destructive-foreground`

### Spacing & Layout
- **Page Container:** `p-8` (desktop), `p-4` (mobile)
- **Section Gaps:** `space-y-6`
- **Card Padding:** `p-6`
- **Table Row Height:** `h-14` (fixed for consistency across all data tables)
- **Side Sheet Width:** `w-full sm:max-w-[540px]`

### Status Badge Color Tokens

**Invoice Status:**
| Status | Classes | Indonesian Label |
|--------|---------|-----------------|
| `paid` | `bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300` | Lunas |
| `unpaid` | `bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300` | Belum Lunas |
| `overdue` | `bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300` | Terlambat |

**Payment Status:**
| Status | Classes | Indonesian Label |
|--------|---------|-----------------|
| `confirmed` | `bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300` | Dikonfirmasi |
| `pending` | `bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300` | Menunggu Konfirmasi |
| `rejected` | `bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300` | Ditolak |
| `refunded` | `bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300` | Dikembalikan |

**Cash Entry Status:**
| Status | Classes | Indonesian Label |
|--------|---------|-----------------|
| `approved` | `bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300` | Disetujui |
| `pending` | `bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300` | Menunggu |
| `rejected` | `bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300` | Ditolak |

**Currency formatting:** `Rp 15.230.000` — Indonesian period-as-thousands-separator, "Rp " prefix (no dot after Rp)

---

## 2. Page Inventory

| Route | View Name | Portal |
|-------|-----------|--------|
| `/invoices` | Manajemen Tagihan | Admin |
| `/payments` | Verifikasi Pembayaran | Admin |
| `/cash` | Kas & Dana Kecil | Admin |
| `/customer/invoices` | Tagihan Saya | Customer Portal |
| `/customer/payments` | Riwayat Pembayaran | Customer Portal |
| `/agent/invoices` | Tagihan Klien | Agent Portal |

---

## 3. Component Inventory

### Reused Components (from Phase 2/3)
- `DataTable` — TanStack Table wrapper with pagination
- `DataTableFacetedFilter` — Dropdown faceted filter (used for status, type, method)
- `DataTableColumnHeader` — Sortable column header with sort controls
- `DataTablePagination` — Offset-based pagination footer
- `DataTableToolbar` — Toolbar with search input + filter chips + action buttons
- `ConfirmDialog` / `AlertDialog` — Generic confirmation dialog (Phase 3 pattern)

### New Components

| Component | Location | Purpose |
|-----------|----------|---------|
| `InvoiceDetailSheet` | `features/billing/invoices/components/` | Side sheet showing invoice line items, total, due date, linked payment history |
| `CashRejectDialog` | `features/billing/cash/components/` | Dialog with `Textarea` for cash entry rejection reason |
| `PettyCashCard` | `features/billing/cash/components/` | Balance display card with "Tambah Saldo" (top-up) button |
| `TopUpDialog` | `features/billing/cash/components/` | Modal to top up petty cash fund (amount + notes) |
| `DateRangeFilter` | `features/billing/payments/components/` | Popover with two date inputs (dari / sampai) for payment date range filter |
| `InvoiceGenerationTrigger` | `features/billing/invoices/components/` | Button + `ConfirmDialog` to trigger `POST /api/v1/invoices/trigger-monthly` |
| `PaymentActionMenu` | `features/billing/payments/components/` | Row action dropdown: Konfirmasi / Tolak / Kembalikan Dana / Buka Halaman Pembayaran |

---

## 4. Page Specifications

### Admin: Manajemen Tagihan (`/invoices`)

**Layout:** Full-width page with page header + `DataTableToolbar` + `DataTable`.

**Page Header:**
- Title: `"Manajemen Tagihan"`
- Subtitle: `"Kelola tagihan pelanggan dan buat tagihan bulanan"`
- Right slot: `InvoiceGenerationTrigger` button ("Buat Tagihan Bulanan")

**DataTableToolbar:**
- Search input: placeholder `"Cari nomor tagihan atau pelanggan..."`
- Faceted filter 1: `Status` — options: Lunas, Belum Lunas, Terlambat
- Faceted filter 2: `Jatuh Tempo` — options: Terlambat (Yes), Belum Terlambat (No)
- Column visibility toggle (standard DataTable pattern)

**Table Columns:**
| Column | Key | Sortable | Notes |
|--------|-----|----------|-------|
| No. Tagihan | `invoice_number` | Yes | Monospace font `font-mono text-xs` |
| Pelanggan | `customer.name` | Yes | |
| Tanggal | `issue_date` | Yes | `dd/MM/yyyy` format |
| Jatuh Tempo | `payment_deadline` | Yes | `dd/MM/yyyy`; red text if overdue |
| Total | `total_amount` | Yes | `Rp` format, right-aligned |
| Status | `status` | No | Status badge |
| Aksi | — | No | `...` icon button → row click opens sheet |

**Row interaction:** Clicking anywhere on a row opens `InvoiceDetailSheet`.

**InvoiceDetailSheet:**
- Header: `"Detail Tagihan #[invoice_number]"`
- Sections:
  1. **Info Tagihan:** customer name, period, issue date, due date, type
  2. **Jumlah:** subtotal, tax (if any), total — all Rp formatted, total in `text-xl font-bold`
  3. **Status:** status badge inline
  4. **Riwayat Pembayaran:** list of linked payments (date, method, amount, status badge)
- Footer: Close button ("Tutup")

**Empty state:** `"Tidak ada tagihan ditemukan."` with `FileText` icon, muted text.

**Loading state:** Skeleton rows — 6 rows × columns, pulsing `animate-pulse`.

---

### Admin: Verifikasi Pembayaran (`/payments`)

**Layout:** Full-width page with page header + `DataTableToolbar` + `DataTable`.

**Page Header:**
- Title: `"Riwayat Pembayaran"`
- Subtitle: `"Konfirmasi, tolak, atau kembalikan dana pembayaran pelanggan"`

**DataTableToolbar:**
- Search input: placeholder `"Cari nomor referensi atau tagihan..."`
- Faceted filter 1: `Metode` — Cash, Transfer Bank, E-Wallet, Kartu Kredit, Kartu Debit, QRIS, Gateway
- Faceted filter 2: `Status` — Menunggu Konfirmasi, Dikonfirmasi, Ditolak, Dikembalikan
- `DateRangeFilter` — "Tanggal Pembayaran" label, dari/sampai date inputs

**Table Columns:**
| Column | Key | Sortable | Notes |
|--------|-----|----------|-------|
| No. Referensi | `payment_number` | Yes | Monospace `font-mono text-xs` |
| No. Tagihan | `invoice.invoice_number` | No | |
| Metode | `payment_method` | No | Capitalized label |
| Jumlah | `amount` | Yes | `Rp` format, right-aligned |
| Tanggal | `payment_date` | Yes | `dd/MM/yyyy` |
| Status | `status` | No | Status badge |
| Aksi | — | No | `PaymentActionMenu` dropdown |

**PaymentActionMenu options (conditional by status):**
- `pending` → "Konfirmasi Pembayaran", "Tolak Pembayaran"
- `confirmed` → "Kembalikan Dana"
- `gateway` type → "Buka Halaman Pembayaran" (opens new tab)

**Confirmation Dialogs:**

| Action | Dialog Title | Dialog Description | Button Label | Button Variant |
|--------|-------------|-------------------|-------------|----------------|
| Konfirmasi | "Konfirmasi Pembayaran?" | "Pembayaran ini akan ditandai sebagai dikonfirmasi. Tindakan ini tidak dapat dibatalkan." | "Konfirmasi" | `default` |
| Tolak | "Tolak Pembayaran?" | "Masukkan alasan penolakan" + `Textarea` (required) | "Tolak" | `destructive` |
| Kembalikan Dana | "Kembalikan Dana?" | "Dana sebesar [Rp amount] akan dikembalikan. Ini tidak dapat dibatalkan." | "Kembalikan Dana" | `destructive` |

**Empty state:** `"Tidak ada pembayaran ditemukan."` with receipt icon.

---

### Admin: Kas & Dana Kecil (`/cash`)

**Layout:** Page header → `PettyCashCard` section → Cash Entries table section, separated by `space-y-6`.

**Page Header:**
- Title: `"Kas & Dana Kecil"`
- Subtitle: `"Kelola entri kas dan saldo dana kecil"`
- Right slot: "Tambah Entri Kas" button → opens `CreateCashEntryDialog`

**PettyCashCard:**
- Full-width card, `bg-card border rounded-lg p-6`
- Left: label `"Saldo Dana Kecil Saat Ini"` (`text-sm text-muted-foreground`), balance amount (`text-3xl font-bold`)
- Right: "Tambah Saldo" button → opens `TopUpDialog`
- If no petty cash fund exists: `"Dana kecil belum dikonfigurasi. Buat dana kecil baru."` with create action

**Cash Entries DataTableToolbar:**
- Search: placeholder `"Cari deskripsi atau sumber..."`
- Faceted filter 1: `Status` — Menunggu, Disetujui, Ditolak
- Faceted filter 2: `Tipe` — Masuk (income), Keluar (expense)

**Cash Entry Table Columns:**
| Column | Key | Sortable | Notes |
|--------|-----|----------|-------|
| Tanggal | `date` | Yes | `dd/MM/yyyy` |
| Tipe | `type` | No | Badge: Masuk (green) / Keluar (red) |
| Sumber | `source` | No | |
| Deskripsi | `description` | No | Truncate at 40 chars |
| Jumlah | `amount` | Yes | `Rp` format |
| Status | `status` | No | Status badge |
| Aksi | — | No | Inline: ✓ Approve button (green) + ✗ Reject button (red) for pending; none for others |

**Inline Approve:** Single icon button (`CheckIcon`, `text-green-600`) — no dialog. On click: call `POST /cash-entries/{id}/approve` immediately, optimistic update or refetch.

**CashRejectDialog:**
- Title: `"Tolak Entri Kas"`
- Description: `"Berikan alasan penolakan untuk entri ini."`
- `Textarea` field, label: `"Alasan"`, placeholder: `"Masukkan alasan penolakan..."`, required
- Buttons: "Batal" (ghost) + "Tolak" (destructive)

**Empty state:** `"Tidak ada entri kas ditemukan."` with banknote icon.

---

### Customer Portal: Tagihan Saya (`/customer/invoices`)

**Layout:** Full-width table (same pattern as admin invoices, simplified).

**Page Header:** Title `"Tagihan Saya"`, subtitle `"Lihat dan bayar tagihan langganan Anda"`

**DataTableToolbar:**
- Search: `"Cari nomor tagihan..."`
- Faceted filter: `Status` — Lunas, Belum Lunas, Terlambat

**Table Columns:** No. Tagihan, Periode, Jatuh Tempo, Total, Status, Aksi

**Row Actions:**
- "Detail" → opens `InvoiceDetailSheet` (same component, read-only from customer perspective)
- If `status === 'unpaid' || status === 'overdue'`: "Bayar Sekarang" button → calls `POST /portal/v1/payments/{id}/pay` → opens `payment_url` in new tab → toast `"Halaman pembayaran dibuka di tab baru"`

**Empty state:** `"Belum ada tagihan."`

---

### Customer Portal: Riwayat Pembayaran (`/customer/payments`)

**Layout:** Read-only table. No action buttons.

**Table Columns:** No. Referensi, No. Tagihan, Metode, Jumlah, Tanggal, Status

**Empty state:** `"Belum ada riwayat pembayaran."`

---

### Agent Portal: Tagihan Klien (`/agent/invoices`)

**Layout:** Simple data table.

**Page Header:** Title `"Tagihan Klien"`, subtitle `"Ajukan permintaan pembayaran untuk tagihan pelanggan"`

**Table Columns:** No. Tagihan, Pelanggan, Total, Status, Aksi

**Row Actions:** "Ajukan Pembayaran" button (visible for unpaid invoices) → calls agent payment request API → toast on success.

**Empty state:** `"Tidak ada tagihan klien ditemukan."`

---

## 5. Interaction Contracts

| Interaction | Trigger | API Call | Success | Error |
|-------------|---------|----------|---------|-------|
| Open invoice detail | Click table row | `GET /api/v1/invoices/{id}` | Sheet opens with data | Toast: `"Gagal memuat detail tagihan"` |
| Trigger monthly generation | "Buat Tagihan Bulanan" → confirm | `POST /api/v1/invoices/trigger-monthly` | Toast: `"Tagihan bulanan berhasil dibuat"` | Toast: `"Gagal membuat tagihan: [error]"` |
| Confirm payment | "Konfirmasi" → dialog confirm | `POST /api/v1/payments/{id}/confirm` | Toast: `"Pembayaran berhasil dikonfirmasi"` + row status updates | Toast: `"Gagal mengkonfirmasi pembayaran"` |
| Reject payment | "Tolak" → dialog with reason | `POST /api/v1/payments/{id}/reject` | Toast: `"Pembayaran ditolak"` + row updates | Toast: `"Gagal menolak pembayaran"` |
| Refund payment | "Kembalikan Dana" → dialog confirm | `POST /api/v1/payments/{id}/refund` | Toast: `"Refund berhasil diproses"` + row updates | Toast: `"Gagal memproses refund"` |
| Open gateway URL | "Buka Halaman Pembayaran" | — | `window.open(payment_url, '_blank')` + Toast: `"Halaman pembayaran dibuka di tab baru"` | — |
| Initiate gateway (customer) | "Bayar Sekarang" | `POST /portal/v1/payments/{id}/pay` | `window.open(payment_url, '_blank')` + Toast: `"Halaman pembayaran dibuka di tab baru"` | Toast: `"Gagal memuat halaman pembayaran"` |
| Approve cash entry | Inline ✓ button | `POST /api/v1/cash-entries/{id}/approve` | Toast: `"Entri kas disetujui"` + row status updates | Toast: `"Gagal menyetujui entri kas"` |
| Reject cash entry | Inline ✗ → dialog with reason | `POST /api/v1/cash-entries/{id}/reject` | Toast: `"Entri kas ditolak"` + row updates | Toast: `"Gagal menolak entri kas"` |
| Top up petty cash | "Tambah Saldo" → dialog | `POST /api/v1/petty-cash/{id}/topup` | Toast: `"Saldo berhasil ditambahkan"` + balance updates | Toast: `"Gagal menambahkan saldo"` |

---

## 6. Copywriting (Bahasa Indonesia)

### Page Titles & Subtitles
- Invoices admin: `"Manajemen Tagihan"` / `"Kelola tagihan pelanggan dan buat tagihan bulanan"`
- Payments admin: `"Riwayat Pembayaran"` / `"Konfirmasi, tolak, atau kembalikan dana pembayaran pelanggan"`
- Cash admin: `"Kas & Dana Kecil"` / `"Kelola entri kas dan saldo dana kecil"`
- Customer invoices: `"Tagihan Saya"` / `"Lihat dan bayar tagihan langganan Anda"`
- Customer payments: `"Riwayat Pembayaran Saya"` / `"Lihat semua riwayat pembayaran Anda"`
- Agent invoices: `"Tagihan Klien"` / `"Ajukan permintaan pembayaran untuk tagihan pelanggan"`

### Table Column Headers
| Key | Indonesian |
|-----|-----------|
| invoice_number | No. Tagihan |
| payment_number | No. Referensi |
| customer.name | Pelanggan |
| issue_date | Tanggal Tagihan |
| payment_deadline | Jatuh Tempo |
| payment_date | Tanggal Pembayaran |
| total_amount / amount | Jumlah |
| status | Status |
| payment_method | Metode |
| source | Sumber |
| description | Deskripsi |
| type | Tipe |

### Button Labels
| Action | Label |
|--------|-------|
| Create monthly invoices | `"Buat Tagihan Bulanan"` |
| Add cash entry | `"Tambah Entri Kas"` |
| Top up petty cash | `"Tambah Saldo"` |
| Confirm payment | `"Konfirmasi"` |
| Reject payment | `"Tolak"` |
| Refund payment | `"Kembalikan Dana"` |
| Open gateway URL | `"Buka Halaman Pembayaran"` |
| Customer pay now | `"Bayar Sekarang"` |
| Agent request payment | `"Ajukan Pembayaran"` |
| Approve cash entry | `"Setujui"` (tooltip) |
| Reject cash entry | `"Tolak"` (tooltip) |
| Close sheet | `"Tutup"` |
| Cancel (dialogs) | `"Batal"` |

### Dialog Titles & Descriptions
| Dialog | Title | Description |
|--------|-------|-------------|
| Monthly generation confirm | `"Buat Tagihan Bulanan?"` | `"Tindakan ini akan membuat tagihan untuk semua pelanggan aktif bulan ini. Lanjutkan?"` |
| Confirm payment | `"Konfirmasi Pembayaran?"` | `"Pembayaran ini akan ditandai sebagai dikonfirmasi. Tindakan ini tidak dapat dibatalkan."` |
| Reject payment | `"Tolak Pembayaran?"` | `"Berikan alasan penolakan agar pelanggan dapat mengirimkan ulang bukti yang benar."` |
| Refund payment | `"Kembalikan Dana?"` | `"Dana sebesar [amount] akan dikembalikan kepada pelanggan. Ini tidak dapat dibatalkan."` |
| Reject cash entry | `"Tolak Entri Kas"` | `"Berikan alasan penolakan untuk entri ini."` |
| Top up petty cash | `"Tambah Saldo Dana Kecil"` | `"Masukkan jumlah dan keterangan penambahan saldo."` |

### Empty States
| View | Message |
|------|---------|
| Admin invoices | `"Tidak ada tagihan ditemukan."` |
| Admin payments | `"Tidak ada pembayaran ditemukan."` |
| Admin cash entries | `"Tidak ada entri kas ditemukan."` |
| Customer invoices | `"Belum ada tagihan."` |
| Customer payments | `"Belum ada riwayat pembayaran."` |
| Agent invoices | `"Tidak ada tagihan klien ditemukan."` |

### Toast Messages
| Event | Toast |
|-------|-------|
| Monthly invoices generated | `"Tagihan bulanan berhasil dibuat"` |
| Payment confirmed | `"Pembayaran berhasil dikonfirmasi"` |
| Payment rejected | `"Pembayaran ditolak"` |
| Payment refunded | `"Refund berhasil diproses"` |
| Gateway URL opened | `"Halaman pembayaran dibuka di tab baru"` |
| Cash entry approved | `"Entri kas disetujui"` |
| Cash entry rejected | `"Entri kas ditolak"` |
| Petty cash topped up | `"Saldo berhasil ditambahkan"` |
| Generic error | `"[Action] gagal. Silakan coba lagi."` |

---

## 7. Responsive Behavior

- **Tables:** Horizontal scroll enabled on `sm` and below. Priority columns visible at all widths: No. Tagihan/Referensi, Jumlah, Status, Aksi. Secondary columns (dates, method) visible from `md` up.
- **Side Sheets:** `w-full` on `< sm` (full viewport width). `max-w-[540px]` on `sm+`.
- **DataTableToolbar:** Filters collapse into a "Filter" button on `< md` that opens a Sheet with full filter options.
- **PettyCashCard:** On `< md`, balance and top-up button stack vertically.
- **Page padding:** `p-4` on mobile (`< md`), `p-8` on desktop.

---

*Phase: 04-billing-payments*
*UI-SPEC generated: 2026-04-04*
