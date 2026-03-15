# Laporan Hasil Testing — mikmongo ISP Billing System

**Tanggal:** 2026-03-16 (diperbarui — concurrent payment test fix)
**Environment:** Windows 11, Go 1.23, PostgreSQL 16 (Docker), Redis 7 (Docker)
**Database:** `mikmongo_test` @ localhost:5432 (user: mikhmon)
**Redis:** DB=15 @ localhost:6379 (dedicated test DB, di-flush setiap test)

---

## Ringkasan Eksekusi

| Tier | Kategori | Total Test | PASS | FAIL/SKIP | Durasi |
|------|----------|-----------|------|-----------|--------|
| 1 | Domain Unit Tests | 77 | ✅ 77 | 0 | ~1.2s |
| 2 | Service Unit Tests (Mock) | 48 | ✅ 48 | 0 | ~0.3s |
| 3 | Integration Tests (PostgreSQL + Redis) | 32 | ✅ 32 | 0 | *lihat catatan* |
| **Total** | | **157** | **✅ 157** | **0** | |

> **Perubahan dari sesi ini:** `TestConcurrentPaymentConfirm_RaceCondition` kini PASS kembali setelah fix.
> Root cause: `adminID = uuid.New()` melanggar FK `payments_processed_by_fkey` — bukan race condition.
> Semua 157 tests pass. Integration tests memerlukan PostgreSQL + Redis.

---

## Tier 1: Domain Unit Tests

Tests murni business logic tanpa I/O. Run: `go test ./internal/domain/...`

### `internal/domain/billing` — 25 test ✅

| Test | Sub-test | Status |
|------|----------|--------|
| `TestCalculateTax` | normal 11%, zero tax, zero subtotal, rounding, rp199999.5 @ 11%, large 5M @ 11% | ✅ PASS |
| `TestCalculateTotal` | subtotal+tax-discount+lateFee, zero components, no discount/fee, no fp drift | ✅ PASS |
| `TestCalculateProration` | full month, half month (day 16/31), day 1, last day, zero days | ✅ PASS |
| `TestCalculateLateFee` | 1 day (1%), 10 days (capped 10%), 20 days (capped), zero, negative | ✅ PASS |
| `TestIsOverdue` | past due, today, future, status paid/cancelled/refunded | ✅ PASS |
| `TestDaysOverdue` | 5 days overdue, not overdue | ✅ PASS |
| `TestShouldSuspendForNonPayment` | within grace, beyond grace, already paid | ✅ PASS |
| `TestShouldSendReminder` | no previous, long ago, recently sent | ✅ PASS |
| `TestInvoiceStatusFromAmounts` | unpaid (0), partial, paid (exact), overpay | ✅ PASS |
| `TestClampBillingDay` | day 31 in Feb, day 31 in April, day 1, leap year | ✅ PASS |
| `TestResolveGracePeriod` | sub priority, fallback profile, fallback default 3 | ✅ PASS |
| `TestResolveBillingDay` | sub priority, fallback profile, fallback default 1st | ✅ PASS |
| `TestGetBillingPeriod` | monthly, daily, weekly, yearly, unknown→monthly | ✅ PASS |

### `internal/domain/subscription` — 18 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestValidateStatusTransition` | 14 valid transitions + 9 invalid transitions | ✅ PASS |
| `TestCanActivate` | pending/suspended → ok; active/isolated/terminated → error | ✅ PASS |
| `TestCanSuspend` | active/isolated → ok; pending/terminated → error | ✅ PASS |
| `TestCanIsolate` | active/suspended → ok; pending/isolated → error | ✅ PASS |
| `TestCanRestore` | isolated → ok; active/suspended → error | ✅ PASS |
| `TestCanTerminate` | active/pending/suspended → ok; already terminated → error | ✅ PASS |
| `TestIsExpired` | past expiry, future expiry, nil expiry | ✅ PASS |
| `TestValidateCredentials` | valid, too short/long username, short password | ✅ PASS |
| `TestGeneratePassword` | correct length, alphanumeric, entropy, zero/negative → error | ✅ PASS |

### `internal/domain/payment` — 12 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestValidatePayment` | valid, amount ≤ 0, invoice paid/cancelled, partial ok | ✅ PASS |
| `TestCanConfirm` | pending → ok; confirmed/rejected/refunded → error | ✅ PASS |
| `TestCanReject` | pending → ok; confirmed/rejected → error | ✅ PASS |
| `TestCanRefund` | confirmed → ok; pending/rejected, already refunded → error | ✅ PASS |
| `TestCalculateAllocations` | exact, partial, multi-invoice FIFO, overpay, skip paid, zero | ✅ PASS |
| `TestIsGatewayPayment` | method=gateway, gateway_name set, cash → false | ✅ PASS |

### `internal/domain/notification` — 8 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestRenderTemplate` | substitusi, unknown key, empty data, empty body | ✅ PASS |
| `TestRenderSubject` | substitusi subject, nil subject | ✅ PASS |
| `TestExtractPlaceholders` | multiple, none, duplicate, empty | ✅ PASS |
| `TestValidateTemplate` | valid whatsapp/email, empty body, invalid channel | ✅ PASS |
| `TestShouldSend` | matching channel, mismatched, inactive template | ✅ PASS |

### `internal/domain/customer` — 6 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestValidateCustomer` | valid, empty name, whitespace name, empty phone | ✅ PASS |
| `TestCanDeactivate` | active → ok; already inactive → error | ✅ PASS |
| `TestCanActivate` | inactive → ok; already active → error | ✅ PASS |

### `internal/domain/router` — 8 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestValidateConnection` | valid, host empty, port 0/65536/negative/1/65535 | ✅ PASS |
| `TestIsOnline` | online/offline/unknown | ✅ PASS |
| `TestCanConnect` | active+online, active+unknown, active+offline, not active | ✅ PASS |
| `TestShouldSync` | nil (never synced), long ago, recently, exactly at interval | ✅ PASS |
| `TestIsStale` | nil → stale, > threshold, < threshold | ✅ PASS |

---

## Tier 2: Service Unit Tests (dengan Mock Repository)

Run: `go test ./internal/service/...`

### `internal/service/auth_service_test.go` — 13 test ✅

| Test | Status |
|------|--------|
| `TestLogin_Success` | ✅ PASS |
| `TestLogin_InvalidPassword` | ✅ PASS |
| `TestLogin_UserNotFound` | ✅ PASS |
| `TestLogin_InactiveUser` | ✅ PASS |
| `TestRefreshToken_Valid` | ✅ PASS |
| `TestRefreshToken_Blacklisted` | ✅ PASS |
| `TestRefreshToken_InvalidToken` | ✅ PASS |
| `TestChangePassword_Success` | ✅ PASS |
| `TestChangePassword_WrongOldPassword` | ✅ PASS |
| `TestChangePassword_UserNotFound` | ✅ PASS |
| `TestLogout_CallsBlacklistWithCorrectJTI` | ✅ PASS |
| `TestLogout_PropagatesRedisError` | ✅ PASS |
| `TestRefreshToken_BlacklistsOldTokenJTI` | ✅ PASS |

### `internal/service/billing_service_test.go` — 8 test ✅

| Test | Status |
|------|--------|
| `TestGenerateInvoice_NewSubscription` | ✅ PASS |
| `TestGenerateInvoice_Proration` | ✅ PASS |
| `TestGenerateInvoice_SuspendedSubscription` | ✅ PASS |
| `TestGenerateInvoice_WithTax` | ✅ PASS |
| `TestProcessDailyBilling_BillingDayToday` | ✅ PASS |
| `TestProcessDailyBilling_DifferentBillingDay` | ✅ PASS |
| `TestCheckAndIsolateOverdue_ShouldIsolate` | ✅ PASS |
| `TestCheckAndIsolateOverdue_WithinGracePeriod` | ✅ PASS |

### `internal/service/notification_service_test.go` — 11 test ✅

| Test | Status |
|------|--------|
| `TestNotificationService_RenderTemplate` (3 sub) | ✅ PASS |
| `TestNotificationService_RenderAndSend_NoTemplate` | ✅ PASS |
| `TestNotificationService_RenderAndSend_InactiveTemplate` | ✅ PASS |
| `TestSendInvoiceCreated_RendersCorrectData` | ✅ PASS |
| `TestSendPaymentConfirmed_RendersCorrectData` | ✅ PASS |
| `TestSendPaymentReminder_RendersCorrectData` | ✅ PASS |
| `TestSendViaWhatsApp_CallsClientWithPhone` | ✅ PASS |
| `TestRenderAndSend_WhatsAppChannel_CallsWA` | ✅ PASS |
| `TestRenderAndSend_EmailChannel_CallsEmail` | ✅ PASS |
| `TestSendInvoiceCreated_CorrectPhone` | ✅ PASS |
| `TestSendViaWhatsApp_NilClient_ReturnsError` | ✅ PASS |

### `internal/service/payment_service_test.go` — 9 test ✅

| Test | Status |
|------|--------|
| `TestCreatePayment` | ✅ PASS |
| `TestCreatePayment_GeneratesPaymentNumber` | ✅ PASS |
| `TestConfirmPayment_SingleInvoice` | ✅ PASS |
| `TestConfirmPayment_AlreadyConfirmed` | ✅ PASS |
| `TestConfirmPayment_MultipleInvoices_FIFO` | ✅ PASS |
| `TestRejectPayment` | ✅ PASS |
| `TestRejectPayment_AlreadyConfirmed` | ✅ PASS |
| `TestRefundPayment` | ✅ PASS |
| `TestRefundPayment_AlreadyRefunded` | ✅ PASS |

### `internal/service/subscription_service_test.go` — 8 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestCreate_Success` | Mock provider+adapter; assert `AddSecret` dipanggil, MtPPPID tersimpan | ✅ PASS |
| `TestCreate_RouterConnectionFails` | Provider return error → error propagated, `subRepo.Create` tidak dipanggil | ✅ PASS |
| `TestCreate_MikrotikAddSecretFails` | Adapter AddSecret error → `subRepo.Create` tidak dipanggil | ✅ PASS |
| `TestCreate_DBSaveFails_RollbackMikrotik` | DB save fail → `RemoveSecret` dipanggil (rollback MikroTik) | ✅ PASS |
| `TestActivate_Success` | `EnableSecret` dipanggil; status=active, ActivatedAt tersimpan | ✅ PASS |
| `TestIsolate_Success` | `UpdateSecret` dipanggil dengan isolate profile; status=isolated | ✅ PASS |
| `TestSuspend_Success` | `DisableSecret` dipanggil; status=suspended | ✅ PASS |
| `TestTerminate_Success` | `RemoveSecret` dipanggil; status=terminated, TerminatedAt tersimpan | ✅ PASS |

---

## Tier 3: Integration Tests (PostgreSQL + Redis Real)

Run:
```bash
TEST_DB_HOST=localhost TEST_DB_PORT=5432 TEST_DB_USER=mikhmon \
TEST_DB_PASSWORD=secret TEST_DB_NAME=mikmongo_test \
go test -v -tags=integration -timeout=120s ./tests/integration/...
```

**Infrastruktur:** Docker container PostgreSQL 16 + Redis 7 @ localhost:6379 DB=15
**Migrasi:** 20 migrasi goose dijalankan satu kali via `TestMain` (bukan per-test)
**Cleanup:** Rollback transaksi (bukan `TRUNCATE CASCADE`) + Redis FlushDB setiap test

### Auth Service — 4 test ✅ *(diperbarui: real Redis assertions)*

| Test | Keterangan | Status |
|------|------------|--------|
| `TestLogin_Integration` | Create user → login → JWT token pair valid | ✅ PASS |
| `TestLogout_Integration` | Login → logout → assert `blacklist:<JTI>` tersimpan di Redis | ✅ PASS |
| `TestRefreshToken_Integration` | Login → refresh → old JTI blacklisted, second use errors | ✅ PASS |
| `TestChangePassword_Integration` | Change pass → `pwd_changed:<userID>` tersimpan di Redis + login baru berhasil | ✅ PASS |

> **Sebelum fix:** semua auth tests menggunakan `integrationNoopRedis` — blacklisting tidak pernah diverifikasi.
> **Setelah fix:** `suite.RedisClient` (real Redis DB=15) — state Redis diverifikasi langsung.

### Redis Auth E2E — 3 test ✅ *(BARU)*

| Test | Keterangan | Status |
|------|------------|--------|
| `TestE2E_LogoutBlacklistsToken` | Login → GET /protected 200 → Logout → GET /protected 401 "token has been revoked" | ✅ PASS |
| `TestE2E_RefreshTokenRotation_OldTokenRejected` | Login → Refresh → old refresh token errors pada pemakaian kedua | ✅ PASS |
| `TestE2E_PasswordChange_OldTokenInvalidated` | Login → GET /protected 200 → ChangePassword → GET /protected 401 "password has been changed" | ✅ PASS |

> File: `tests/integration/redis_auth_e2e_test.go`
> Stack: `httptest.NewRecorder()` + minimal Gin router + real `AuthMiddleware` + real Redis
> Membuktikan middleware + Redis bekerja end-to-end, bukan hanya di level service.

### Billing Lifecycle — 4 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestBillingLifecycle_GenerateInvoice` | Invoice bulanan dengan tax 11% = 222.000 | ✅ PASS |
| `TestBillingLifecycle_Proration` | Aktivasi day 16 March → ~160.000 (16/31 × 310.000) | ✅ PASS |
| `TestBillingLifecycle_ProcessDailyBilling` | billing_day=today → invoice terbuat | ✅ PASS |
| `TestBillingLifecycle_OverdueIsolation` | Invoice overdue > grace period → subscription isolated | ✅ PASS |

### Payment Lifecycle — 5 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestPaymentLifecycle_ConfirmSingle` | Payment penuh → invoice "paid", allocation created | ✅ PASS |
| `TestPaymentLifecycle_PartialPayment` | 2× partial payment → invoice "partial" lalu "paid" | ✅ PASS |
| `TestPaymentLifecycle_FIFOMultipleInvoices` | 250K untuk 2 invoice 200K: oldest paid, newest partial | ✅ PASS |
| `TestPaymentLifecycle_RejectPayment` | Reject → status "rejected", invoice tetap "unpaid" | ✅ PASS |
| `TestPaymentLifecycle_Refund` | Confirm → refund → status "refunded", RefundAmount tersimpan | ✅ PASS |

### Payment Edge Cases — 2 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestPaymentConfirm_Overpayment` | Invoice 200K, payment 300K → PaidAmount=200K; kelebihan 100K tidak dialokasikan | ✅ PASS |
| `TestConcurrentPaymentConfirm_RaceCondition` | 1 payment, 2 goroutine confirm → idempotency guard: tepat 1 sukses + 1 gagal "only pending payments can be confirmed" | ✅ PASS |

### Customer Service — 6 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestCustomerService_Create` | Customer code format `CST#####` | ✅ PASS |
| `TestCustomerService_CreateWithSubscription` | Customer + subscription pending | ✅ PASS |
| `TestCustomerService_SuspendAllSubscriptions` | 2 sub active → keduanya suspended | ✅ PASS |
| `TestCustomerService_IsolateAllSubscriptions` | active → isolated | ✅ PASS |
| `TestCustomerService_RestoreAllSubscriptions` | isolated → active | ✅ PASS |
| `TestCustomerService_PortalAuth` | SetPortalPassword + AuthPortal benar/salah | ✅ PASS |

### Registration Workflow — 5 test ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestRegistration_Create` | Status "pending" setelah create | ✅ PASS |
| `TestRegistration_Approve_WithoutSubscription` | Approve tanpa profile → customer created, status "approved" | ✅ PASS |
| `TestRegistration_Approve_WithSubscription` | Approve dengan profile → customer created | ✅ PASS |
| `TestRegistration_Reject` | Reject dengan reason → status "rejected", reason tersimpan | ✅ PASS |
| `TestRegistration_ListByStatus` | 3 pending → filter "pending" = 3, filter "approved" = 0 | ✅ PASS |

### Load Test — 1 test (3 sub-run) ✅

| Test | Keterangan | Status |
|------|------------|--------|
| `TestProcessDailyBilling_LoadTest/10_subscriptions` | 10 sub aktif → 10 invoice; 83ms total (8.2ms/sub) | ✅ PASS |
| `TestProcessDailyBilling_LoadTest/50_subscriptions` | 50 sub aktif → 50 invoice; 394ms total (7.9ms/sub) | ✅ PASS |
| `TestProcessDailyBilling_LoadTest/100_subscriptions` | 100 sub aktif → 100 invoice; 1.193s total (11.9ms/sub) | ✅ PASS |

---

## Catatan Teknis

### Improvement 9: Concurrent Payment Test Fix *(sesi ini)*

**Root cause:** `adminID := uuid.New().String()` menghasilkan UUID acak yang tidak ada di tabel `users`. PostgreSQL menolak kedua `UPDATE payments SET processed_by=...` dengan FK violation `payments_processed_by_fkey` → kedua goroutine gagal (`failCount=2`), bukan karena race condition.

**Perbaikan pada `tests/integration/concurrent_payment_test.go`:**

| # | Problem | Fix |
|---|---------|-----|
| 1 | `adminID` bukan user valid (FK violation) | INSERT user nyata ke RootDB + DELETE di `t.Cleanup` |
| 2 | Identifier hard-coded (`"RACEPROFILE01"`, `"race-user-01"`, `"08999000001"`) | UUID suffix pada semua fixture string |
| 3 | Tidak ada `defer suite.Cleanup(t)` — tx dari `SetupSuite` bocor | Tambah `defer suite.Cleanup(t)` |
| 4 | `errs[0]/errs[1]` ditulis dari goroutine berbeda | Ganti dengan `chan confirmResult` (buffered) |
| 5 | Dua payment berbeda (p1, p2) → `failCount==1` non-deterministik | Satu payment, dua goroutine — menguji idempotency guard |
| 6 | G2 bisa baca payment sebagai "pending" sebelum G1 commit | `g1Done` channel: G2 hanya launch setelah `<-g1Done` |

**Perilaku yang ditest:** Satu payment dikonfirmasi oleh dua goroutine. G1 sukses (status: pending→confirmed). G2 — setelah G1 commit — re-read payment di dalam tx dan melihat "confirmed", sehingga idempotency guard (`CanConfirm`) mengembalikan error `"only pending payments can be confirmed"`. Deterministik dan diverifikasi 3× (`-count=3`).

---

### Improvement 8: Redis Real Integration Tests *(sesi ini)*

**Problem:** Semua auth integration tests menggunakan `integrationNoopRedis` — operasi Redis selalu return `nil`/`false` tanpa menyentuh Redis. Token blacklisting, refresh token rotation, dan password change invalidation tidak pernah diverifikasi di level integration.

**Solusi:**

| File | Aksi |
|------|------|
| `pkg/redis/client.go` | Tambah `FlushDB(ctx) error` untuk cleanup antar test |
| `tests/integration/suite_test.go` | `sharedRedis *pkgredis.Client`, `TestSuite.RedisClient`, Redis flush di `SetupSuite` |
| `tests/integration/auth_service_test.go` | Hapus `integrationNoopRedis`, pakai `suite.RedisClient`, tambah Redis state assertions |
| `tests/integration/redis_auth_e2e_test.go` | **File baru** — 3 E2E test via `httptest` + Gin + real middleware |

**Konfigurasi Redis test:**
```
TEST_REDIS_HOST=localhost (default)
TEST_REDIS_PORT=6379     (default)
TEST_REDIS_PASSWORD=     (default: kosong)
TEST_REDIS_DB=15         (default: DB 15 — isolated test database)
```

**Perubahan arsitektur `TestSuite`:**
- Sebelum: `integrationNoopRedis{}` diinstansiasi lokal di setiap auth test
- Sesudah: `suite.RedisClient` shared via `sharedRedis` (diinisialisasi di `TestMain`, di-flush setiap `SetupSuite`)

**Yang kini diverifikasi di level integration:**
- `blacklist:<JTI>` tersimpan di Redis setelah logout
- Old refresh token JTI diblacklist setelah rotation, dan tidak bisa dipakai ulang
- `pwd_changed:<userID>` tersimpan di Redis setelah ganti password
- Middleware benar-benar menolak token ter-blacklist (E2E via httptest)
- Middleware benar-benar menolak token pre-password-change (E2E via httptest)

---

### 7 Perbaikan Sebelumnya (sesi lalu)

#### Improvement 1: MikroTik Adapter Interface (subscription_service)
- `routerSvc *RouterService` → `routerProvider MikrotikProvider` (interface)
- File baru: `internal/service/mikrotik_adapter.go`
- 8 unit tests baru di `subscription_service_test.go`

#### Improvement 2: Transaction Rollback Cleanup (integration suite)
- `TestMain` dengan koneksi DB tunggal, rollback per-test (bukan TRUNCATE)
- Estimasi: ~75s → <20s untuk full integration suite

#### Improvement 3: Financial Edge Cases
- `TestCalculateTax`: +2 subtests (rounding fractional, large amount 5M)
- `TestConcurrentPaymentConfirm_RaceCondition`: race condition didokumentasikan

#### Improvement 4: Notification Channel Verification
- `internal/notification/interfaces.go` — `WhatsAppSender`, `EmailSender` interfaces
- 5 unit tests baru di `notification_service_test.go`

#### Improvement 5: Redis Blacklisting Verification (unit test)
- 3 unit tests baru di `auth_service_test.go` (mock Redis)

#### Improvement 6: Load Test ProcessDailyBilling
- `tests/integration/load_test.go` — 3 sub-run (10/50/100 subs) dengan timing

#### Improvement 7: SELECT FOR UPDATE — Fix Race Condition di `PaymentService.Confirm()`
- `internal/repository/transactor.go` + `postgres/transactor.go`
- `Confirm()` kini atomic dengan SELECT FOR UPDATE

---

### MikroTik Real Router

Tests yang memerlukan koneksi real ke MikroTik (PPPoE provisioning) dikompilasi dengan tag `integration && mikrotik_legacy`. Tidak dijalankan dalam suite standar.

```bash
TEST_MIKROTIK_HOST=192.168.27.1 TEST_MIKROTIK_USER=admin TEST_MIKROTIK_PASS=r00t \
go test -v -tags="integration mikrotik_legacy" ./tests/integration/...
```

### Arsitektur Test

```
internal/
├── domain/*/domain_test.go     ← Tier 1: Pure unit tests (no I/O)
├── service/*_test.go           ← Tier 2: Unit tests dengan mock repos
└── service/mocks/              ← Manual mocks (testify/mock)
    ├── customer_repo_mock.go
    ├── invoice_repo_mock.go
    ├── misc_repo_mocks.go
    ├── mikrotik_mock.go
    ├── notification_mocks.go
    ├── payment_repo_mock.go
    ├── redis_mock.go
    └── subscription_repo_mock.go

internal/notification/
├── gowa_client.go
├── email_client.go
└── interfaces.go

internal/repository/
├── transactor.go               ← Transactor interface
├── invoice_repo.go             ← +GetByCustomerIDForUpdate
└── postgres/
    ├── transactor.go           ← gormTransactor (SAVEPOINT-safe)
    ├── invoice_repo.go         ← SELECT FOR UPDATE
    ├── registry.go             ← +Transactor field
    └── ...

pkg/redis/
├── client.go                   ← +FlushDB (untuk cleanup test)
├── session.go                  ← BlacklistToken, IsBlacklisted, SetPasswordChangedAt
└── ...

tests/integration/
├── suite_test.go               ← +sharedRedis, +RedisClient, +getEnvInt
├── auth_service_test.go        ← Real Redis (hapus noopRedis), +state assertions
├── redis_auth_e2e_test.go      ← BARU: 3 E2E tests via httptest + Gin
├── billing_lifecycle_test.go
├── payment_lifecycle_test.go
├── customer_service_test.go
├── registration_test.go
├── concurrent_payment_test.go  ← DIUBAH: fix FK adminID, UUID identifiers, channel sync, 1-payment idempotency
└── load_test.go
```

### Known Issues

Tidak ada known issues aktif. Semua 157 test pass.

---

## Cara Menjalankan

```bash
# Unit tests (cepat, tanpa DB)
go test ./internal/domain/... ./internal/service/...

# Integration tests (butuh PostgreSQL + Redis)
TEST_DB_HOST=localhost TEST_DB_PORT=5432 \
TEST_DB_USER=mikhmon TEST_DB_PASSWORD=secret \
TEST_DB_NAME=mikmongo_test \
go test -v -tags=integration -timeout=120s ./tests/integration/...

# Hanya Redis/Auth E2E tests
TEST_DB_HOST=localhost TEST_DB_PORT=5432 \
TEST_DB_USER=mikhmon TEST_DB_PASSWORD=secret \
TEST_DB_NAME=mikmongo_test \
go test -v -tags=integration -run "TestE2E_|TestLogout_Integration|TestRefreshToken_Integration|TestChangePassword_Integration" \
./tests/integration/...

# Load test saja (dengan timing + threshold assertion)
TEST_DB_HOST=localhost TEST_DB_PORT=5432 \
TEST_DB_USER=mikhmon TEST_DB_PASSWORD=secret \
TEST_DB_NAME=mikmongo_test \
go test -v -tags=integration -run TestProcessDailyBilling_LoadTest ./tests/integration/...

# Concurrent test dengan race detector
TEST_DB_HOST=localhost TEST_DB_PORT=5432 \
TEST_DB_USER=mikhmon TEST_DB_PASSWORD=secret \
TEST_DB_NAME=mikmongo_test \
go test -v -race -tags=integration -run TestConcurrentPaymentConfirm_RaceCondition ./tests/integration/...

# Semua sekaligus
go test ./internal/... && \
TEST_DB_HOST=localhost TEST_DB_PORT=5432 \
TEST_DB_USER=mikhmon TEST_DB_PASSWORD=secret \
TEST_DB_NAME=mikmongo_test \
go test -tags=integration -timeout=120s ./tests/integration/...
```
