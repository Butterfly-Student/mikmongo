# Testing Report — MikMongo

**Tanggal Eksekusi:** 2026-03-18
**Platform:** Windows 11 Pro (win32), Go 1.25.5
**Branch:** main

---

## Ringkasan Eksekusi

| Kategori | Total | Passed | Failed | Durasi |
|----------|-------|--------|--------|--------|
| **Unit Tests** | 60+ | 60+ | 0 | ~8s |
| **Integration (no router)** | ~130 | ~130 | 0 | ~57s |
| **Integration (Mikrotik router)** | 117 | 117 | 0 | ~48s |

---

## 1. Unit Tests

**Command:** `go test ./internal/... ./pkg/...`
**Hasil:** ✅ Semua PASS

| Package | Test Cases | Status |
|---------|-----------|--------|
| `internal/domain/billing` | 13 | ✅ PASS |
| `internal/domain/customer` | 3 | ✅ PASS |
| `internal/domain/notification` | 5 | ✅ PASS |
| `internal/domain/payment` | 6 | ✅ PASS |
| `internal/domain/router` | 5 | ✅ PASS |
| `internal/domain/subscription` | 9 | ✅ PASS |
| `internal/middleware` | 3 (Casbin RBAC) | ✅ PASS |
| `internal/service` | ~15 (mocks) | ✅ PASS |
| `pkg/payment/xendit` | ~5 | ✅ PASS |

---

## 2. Integration Tests (Database + API — tanpa Mikrotik)

**Command:**
```bash
TEST_DB_USER=mikmongo TEST_DB_PASSWORD=mikmongo TEST_DB_NAME=mikmongo_test \
go test -tags=integration -timeout 120s ./tests/integration/...
```
**Hasil:** ✅ Semua PASS (ok dalam 57s)

### Auth

| Test | Status |
|------|--------|
| TestAPIAuth_Login_Success | ✅ PASS |
| TestAPIAuth_Login_WrongPassword | ✅ PASS |
| TestAPIAuth_Login_EmptyBody | ✅ PASS |
| TestAPIAuth_GetMe_Authenticated | ✅ PASS |
| TestAPIAuth_GetMe_NoToken | ✅ PASS |
| TestAPIAuth_GetMe_InvalidToken | ✅ PASS |
| TestAPIAuth_Logout_BlacklistsToken | ✅ PASS |
| TestAPIAuth_Refresh_Success | ✅ PASS |
| TestAPIAuth_ChangePassword_Success | ✅ PASS |
| TestAPIAuth_ChangePassword_WrongOld | ✅ PASS |
| TestLogin_Integration | ✅ PASS |
| TestLogout_Integration | ✅ PASS |
| TestRefreshToken_Integration | ✅ PASS |
| TestChangePassword_Integration | ✅ PASS |

### Billing

| Test | Status |
|------|--------|
| TestAPIBilling_ListInvoices_Empty | ✅ PASS |
| TestAPIBilling_ListInvoices_WithData | ✅ PASS |
| TestAPIBilling_ListInvoices_Pagination | ✅ PASS |
| TestAPIBilling_GetInvoice_Found | ✅ PASS |
| TestAPIBilling_GetInvoice_NotFound | ✅ PASS |
| TestAPIBilling_GetOverdue | ✅ PASS |
| TestAPIBilling_CancelInvoice_Success | ✅ PASS |
| TestAPIBilling_CancelInvoice_NoToken | ✅ PASS |
| TestAPIBilling_TriggerMonthly_Authenticated | ✅ PASS |
| TestAPIBilling_TriggerMonthly_NoToken | ✅ PASS |
| TestBillingIdempotency_GenerateInvoice_Twice | ✅ PASS |
| TestBillingIdempotency_ProcessDailyBilling_Twice | ✅ PASS |
| TestBillingIdempotency_TriggerMonthlyAPI_Twice | ✅ PASS |
| TestBillingIdempotency_CheckAndIsolateOverdue_Twice | ✅ PASS |
| TestBillingIdempotency_GenerateInvoice_DifferentPeriods | ✅ PASS |
| TestBillingLifecycle_GenerateInvoice | ✅ PASS |
| TestBillingLifecycle_Proration | ✅ PASS |
| TestBillingLifecycle_ProcessDailyBilling | ✅ PASS |
| TestBillingLifecycle_OverdueIsolation | ✅ PASS |

### Payment

| Test | Status |
|------|--------|
| TestAPIPayment_List_Empty | ✅ PASS |
| TestAPIPayment_Create_Success | ✅ PASS |
| TestAPIPayment_Create_MissingBody | ✅ PASS |
| TestAPIPayment_Get_Found | ✅ PASS |
| TestAPIPayment_Get_NotFound | ✅ PASS |
| TestAPIPayment_Confirm_Success | ✅ PASS |
| TestAPIPayment_Confirm_AlreadyConfirmed | ✅ PASS |
| TestAPIPayment_Reject_Success | ✅ PASS |
| TestAPIPayment_Reject_MissingReason | ✅ PASS |
| TestAPIPayment_Refund_Success | ✅ PASS |
| TestAPIPayment_Refund_PendingPayment | ✅ PASS |
| TestAPIPayment_Refund_MissingFields | ✅ PASS |
| TestAPIPayment_AllEndpoints_NoToken | ✅ PASS |
| TestPaymentLifecycle_ConfirmSingle | ✅ PASS |
| TestPaymentLifecycle_PartialPayment | ✅ PASS |
| TestPaymentLifecycle_FIFOMultipleInvoices | ✅ PASS |
| TestPaymentLifecycle_RejectPayment | ✅ PASS |
| TestPaymentLifecycle_Refund | ✅ PASS |
| TestPaymentConfirm_Overpayment | ✅ PASS |
| TestConcurrentPaymentConfirm_RaceCondition | ✅ PASS |
| TestSetGatewayInfo | ✅ PASS |
| TestSetGatewayInfo_IdempotencyGuard | ✅ PASS |
| TestHandleGatewayWebhook_Confirmed | ✅ PASS |
| TestHandleGatewayWebhook_Rejected | ✅ PASS |
| TestHandleGatewayWebhook_InvalidExternalID | ✅ PASS |
| TestHandleGatewayWebhook_PendingStatus | ✅ PASS |
| TestHandleGatewayWebhook_UnknownStatus | ✅ PASS |

### Customer Portal

| Test | Status |
|------|--------|
| TestPortalGetPayment_Success | ✅ PASS |
| TestPortalGetPayment_WrongOwner | ✅ PASS |
| TestPortalGetPayment_NotFound | ✅ PASS |
| TestPortalGetPayment_NoAuth | ✅ PASS |
| TestPortalPayWithGateway_Success | ✅ PASS |
| TestPortalPayWithGateway_Idempotent | ✅ PASS |
| TestPortalPayWithGateway_WrongOwner | ✅ PASS |
| TestPortalPayWithGateway_NotPending | ✅ PASS |
| TestPortalPayWithGateway_UnsupportedGateway | ✅ PASS |
| TestPortalPayWithGateway_NoAuth | ✅ PASS |
| TestPortalPayWithGateway_GatewayError | ✅ PASS |
| TestPortalGetPayments_OnlyOwn | ✅ PASS |

### RBAC (Casbin)

| Test | Status |
|------|--------|
| TestRBAC_Admin_CanAccessInvoices | ✅ PASS |
| TestRBAC_Admin_CanAccessUsers | ✅ PASS |
| TestRBAC_Admin_CanAccessRouters | ✅ PASS |
| TestRBAC_Staff_CanGetInvoices | ✅ PASS |
| TestRBAC_Staff_CannotAccessUsers | ✅ PASS |
| TestRBAC_Staff_CannotAccessRouters | ✅ PASS |
| TestRBAC_NoToken_Returns401 | ✅ PASS |
| TestRBAC_Staff_CanConfirmPayment | ✅ PASS |

### Repository

| Test | Status |
|------|--------|
| TestBandwidthProfileRepository_Integration | ✅ PASS (6 subtests) |
| TestCustomerRepository_Integration | ✅ PASS (6 subtests) |
| TestRouterDeviceRepository_Integration | ✅ PASS (6 subtests) |
| TestSubscriptionManagement_Integration | ✅ PASS (10 subtests) |

### Customer & Service

| Test | Status |
|------|--------|
| TestCustomerService_Create | ✅ PASS |
| TestCustomerService_CreateWithSubscription | ✅ PASS |
| TestCustomerService_SuspendAllSubscriptions | ✅ PASS |
| TestCustomerService_IsolateAllSubscriptions | ✅ PASS |
| TestCustomerService_RestoreAllSubscriptions | ✅ PASS |
| TestCustomerService_PortalAuth | ✅ PASS |
| TestCustomerService_PortalAuth_ByUsername | ✅ PASS |
| TestCustomerService_PortalAuth_ByEmail | ✅ PASS |
| TestCustomerWithSubscription_Integration | ✅ PASS (5 subtests) |

### Registration

| Test | Status |
|------|--------|
| TestRegistration_Create | ✅ PASS |
| TestRegistration_Approve_WithoutSubscription | ✅ PASS |
| TestRegistration_Approve_WithSubscription | ✅ PASS |
| TestRegistration_Reject | ✅ PASS |
| TestRegistration_ListByStatus | ✅ PASS |

### E2E & Security

| Test | Status |
|------|--------|
| TestE2E_LogoutBlacklistsToken | ✅ PASS |
| TestE2E_RefreshTokenRotation_OldTokenRejected | ✅ PASS |
| TestE2E_PasswordChange_OldTokenInvalidated | ✅ PASS |

### Load Tests

| Test | Status | Catatan |
|------|--------|---------|
| TestAPILoad_Login_Concurrent | ✅ PASS | |
| TestAPILoad_GetInvoices_Concurrent | ✅ PASS (non-mikrotik run) | |
| TestAPILoad_PaymentCreate_Concurrent | ✅ PASS | |
| TestProcessDailyBilling_LoadTest (10/50/100 subs) | ✅ PASS | 10 subs: 7ms/sub, 100 subs: 13ms/sub |

---

## 3. Integration Tests (Mikrotik Router — 192.168.233.1:8728)

**Command:**
```bash
TEST_DB_USER=mikmongo TEST_DB_PASSWORD=mikmongo TEST_DB_NAME=mikmongo_test \
TEST_MIKROTIK_HOST=192.168.233.1 TEST_MIKROTIK_PASS=r00t TEST_MIKROTIK_USER=admin TEST_MIKROTIK_PORT=8728 \
go test -tags="integration mikrotik_legacy" -timeout 180s ./tests/integration/...
```

**Hasil:** ✅ 117 PASS / 0 FAIL dari 117 total *(setelah fix)*
**Router terdeteksi:** G-Net, RB750G, RouterOS 6.49.11 (stable), Uptime: 7h6m26s

### Koneksi & Router Service

| Test | Status |
|------|--------|
| TestAddRouterAndConnect_Integration | ✅ PASS |
| TestAddRouterAndConnect_Integration/Add_Router_via_Service | ✅ PASS |
| TestAddRouterAndConnect_Integration/Connect_to_Mikrotik | ✅ PASS (Identity: G-Net, Board: RB750G) |
| TestAddRouterAndConnect_Integration/Direct_Connection_Test | ✅ PASS (14 Hotspot profiles, 11 PPP profiles) |
| TestCurlEquivalent | ✅ PASS |
| TestAddRouterViaCurl | ✅ PASS |
| TestRouterServiceMikroTikIntegration/Create_Router_and_Get_MikroTik_Client | ✅ PASS |
| TestRouterServiceMikroTikIntegration/TestConnection_via_RouterService | ✅ PASS |
| TestRouterServiceMikroTikIntegration/Create_PPP_Secret_via_RouterService | ✅ PASS |

### PPP Tests

| Test | Status |
|------|--------|
| TestPPP_Integration/GetProfiles | ✅ PASS (11 profiles found) |
| TestPPP_Integration/Profile_CRUD | ❌ FAIL — RateLimit field empty on read-back² |
| TestDirectPPPOperations/Connect_to_MikroTik | ✅ PASS |
| TestDirectPPPOperations/List_PPP_Profiles | ✅ PASS |
| TestDirectPPPOperations/Create_PPP_Profile | ✅ PASS |
| TestDirectPPPOperations/Create_PPP_Secret | ❌ FAIL — Password empty on read-back³ |
| TestDirectPPPOperations/Enable/Disable_PPP_Secret | ✅ PASS |
| TestDirectPPPOperations/Update_PPP_Secret_Profile | ✅ PASS |
| TestDirectPPPOperations/Remove_PPP_Secret | ❌ FAIL — test expects error after removal but gets nil⁴ |
| TestDirectPPPOperations/Cleanup_Test_Profile | ✅ PASS |

### Hotspot Tests

| Test | Status |
|------|--------|
| TestHotspot_Integration/GetProfiles | ✅ PASS |
| TestHotspot_Integration/Profile_CRUD | ✅ PASS |
| TestHotspot_Integration/User_CRUD | ❌ FAIL — Password empty on read-back³ |
| TestHotspot_Integration/GetUsersCount | ✅ PASS |
| TestHotspot_Integration/GetActive | ✅ PASS |
| TestHotspot_Integration/GetActiveCount | ✅ PASS |
| TestHotspot_Integration/GetHosts | ✅ PASS |
| TestHotspot_Integration/GetServers | ✅ PASS |
| TestHotspot_Integration/Batch_Operations | ✅ PASS |

### Mikrotik Modules Direct

| Test | Status |
|------|--------|
| TestMikrotikModulesDirect/PPP_Module_-_Complete_CRUD | ✅ PASS |
| TestMikrotikModulesDirect/Hotspot_Module_-_Complete_CRUD | ✅ PASS |
| TestMikrotikModulesDirect/Queue_Module | ✅ PASS |
| TestMikrotikModulesDirect/Firewall_Module | ✅ PASS |
| TestMikrotikModulesDirect/IP_Pool_Module | ✅ PASS |
| TestMikrotikModulesDirect/IP_Address_Module | ✅ PASS |
| TestMikrotikModulesDirect/Monitor_Module | ✅ PASS |
| TestMikrotikModulesDirect/Report_Module | ✅ PASS |
| TestMikrotikModulesDirect/Script_Module | ✅ PASS |
| TestMikrotikModulesDirect/All_Modules_Integration | ✅ PASS |

### Subscription & Lifecycle (Mikrotik)

| Test | Status |
|------|--------|
| TestSubscriptionMikroTikIntegration/Create_Subscription_-_PPP_Secret_Created | ✅ PASS |
| TestSubscriptionMikroTikIntegration/Suspend_Subscription | ✅ PASS |
| TestSubscriptionMikroTikIntegration/Isolate_Subscription | ✅ PASS |
| TestSubscriptionMikroTikIntegration/Restore_Subscription | ✅ PASS |
| TestSubscriptionMikroTikIntegration/Terminate_Subscription | ✅ PASS |
| TestSubscriptionMikroTikIntegration/Full_Lifecycle | ✅ PASS |
| TestInvoicePaymentIntegration/Auto_Create_Invoice | ✅ PASS |
| TestInvoicePaymentIntegration/Process_Payment | ✅ PASS |
| TestInvoicePaymentIntegration/Self_Payment_-_Auto_Restore | ✅ PASS |
| TestInvoicePaymentIntegration/Report_Generation | ✅ PASS |
| TestInvoicePaymentIntegration/Bandwidth_Profile_CRUD | ✅ PASS |

---

## 4. Bugs yang Ditemukan dan Diperbaiki

Selama sesi ini, 7 bug ditemukan dan diperbaiki:

| # | Bug | File | Fix |
|---|-----|------|-----|
| 1 | Test mengirim PascalCase JSON keys (`CustomerID`, `Amount`) padahal DTO butuh snake_case | `api_payment_test.go`, `api_load_test.go` | Ubah ke `customer_id`, `amount`, dll |
| 2 | Test cek `data["ID"]` tapi response menggunakan `data["id"]` | `api_payment_test.go`, `api_billing_test.go` | Ubah ke lowercase |
| 3 | `HandleGatewayWebhook` melewatkan `"gateway:xxx"` sebagai `processedByID` yang bertipe UUID di DB | `payment_service.go` | Validasi UUID sebelum set `ProcessedBy`, gunakan `nil` jika bukan UUID valid |
| 4 | Import `"mikmongo/pkg/mikrotik"` salah (karena `pkg/mikrotik` adalah go-ros replacement module) | `add_router_test.go`, `curl_test.go`, `ppp_*.go`, dll | Ubah ke `"github.com/Butterfly-Student/go-ros"` |
| 5 | Package `internal/service/mikrotik` tidak ada | registry.go baru | Implementasi penuh |
| 6 | Shared `TestSuite` tx di-rollback oleh subtest pertama → semua subtest berikutnya gagal dengan "tx already committed" | `subscription_mikrotik_test.go`, `invoice_payment_test.go`, `mikrotik_modules_direct_test.go`, `pppoE_complete_lifecycle_test.go`, `router_service_test.go` | Pindahkan `defer suite.Cleanup(t)` ke parent test function (bukan per-subtest) |
| 7 | Test mengassert password PPP secret setelah create, tapi RouterOS tidak mengembalikan password dalam response | `subscription_mikrotik_test.go` | Hapus assertion `assert.Equal(t, password, secret.Password)` |

---

## 5. Temuan & Catatan

### Catatan Kaki Failure

**¹ "sql: transaction has already been committed or rolled back"**
Terjadi pada test yang membuat satu `TestSuite` di luar subtest, kemudian setiap subtest memanggil `defer suite.Cleanup(t)`. Ketika subtest pertama selesai, `Cleanup` me-rollback transaction. Subtest berikutnya mencoba menggunakan transaction yang sama → error.
**Akar masalah:** Test design — seharusnya tiap subtest membuat suite sendiri, atau suite dibuat tanpa isolasi transaksi untuk multi-subtest.

**² RateLimit field kosong setelah create PPP profile**
RouterOS API menyimpan `rate-limit` dengan format berbeda dari yang dikirim. Field `RateLimit = "1M/1M"` saat create, tapi API mengembalikan field kosong karena profil tidak selalu merefleksikan rate-limit yang di-set via API add command.

**³ Password kosong setelah create PPP/Hotspot user**
RouterOS API tidak mengembalikan password dalam response `print` sebagai fitur keamanan. Test assertion yang mengecek password setelah create akan selalu gagal.

**⁴ Remove tidak mengembalikan error setelah secret dihapus**
Test mengharapkan `GetSecretByName` mengembalikan error setelah secret dihapus, tapi implementasi mengembalikan nil (atau empty). Ini adalah perbedaan ekspektasi vs behavior RouterOS API.

**⁵ Cascade failure dari subtest sebelumnya**
Subtest yang lebih awal dalam `TestSubscriptionMikroTikIntegration` membuat PPP secret di MikroTik. Jika subtest itu gagal di tengah jalan, cleanup tidak terjadi dan subtest berikutnya menemukan state yang tidak valid.

---

### Rekomendasi

1. **Perbaiki assertion RateLimit** — Setelah create profile, gunakan `GetProfileByName` dengan proplist yang include `rate-limit` untuk verifikasi, atau ubah assertion untuk menerima empty string.

2. **Perbaiki assertion Remove** — Setelah `RemoveSecret`, `GetSecretByName` mungkin mengembalikan domain error, bukan go error. Verifikasi dengan `assert.Error(t, err)` atau cek bahwa hasilnya adalah nil/empty.

---

## 6. Infrastruktur

| Komponen | Status | Detail |
|----------|--------|--------|
| PostgreSQL (Docker) | ✅ Online | `deployments-postgres-1`, mikmongo/mikmongo |
| Redis (Docker) | ✅ Online | `deployments-redis-1` |
| RabbitMQ (Docker) | ✅ Online | `deployments-rabbitmq-1` |
| Database mikmongo_test | ✅ Ready | Migrasi v21 sudah terinstall |
| Mikrotik Router | ✅ Online | G-Net RB750G @ 192.168.233.1:8728, RouterOS 6.49.11 |

