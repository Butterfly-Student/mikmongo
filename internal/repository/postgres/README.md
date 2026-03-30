# Package postgres

Implementasi repository GORM untuk PostgreSQL. Setiap repository mengimplementasikan interface dari `internal/repository` sehingga service layer tidak bergantung pada detail storage.

## Daftar Repository

| File | Repository | Interface |
|------|-----------|-----------|
| `customer_repo.go` | `customerRepository` | `CustomerRepository` |
| `invoice_repo.go` | `invoiceRepository` | `InvoiceRepository` |
| `invoice_item_repo.go` | `invoiceItemRepository` | `InvoiceItemRepository` |
| `payment_repo.go` | `paymentRepository` | `PaymentRepository` |
| `payment_allocation_repo.go` | `paymentAllocationRepository` | `PaymentAllocationRepository` |
| `subscription_repo.go` | `subscriptionRepository` | `SubscriptionRepository` |
| `bandwidth_profile_repo.go` | `bandwidthProfileRepository` | `BandwidthProfileRepository` |
| `customer_registration_repo.go` | `customerRegistrationRepository` | `CustomerRegistrationRepository` |
| `router_device_repo.go` | `routerDeviceRepository` | `RouterDeviceRepository` |
| `hotspot_sale_repo.go` | `hotspotSaleRepository` | `HotspotSaleRepository` |
| `sales_agent_repo.go` | `salesAgentRepository` | `SalesAgentRepository` |
| `agent_invoice_repo.go` | `agentInvoiceRepository` | `AgentInvoiceRepository` |
| `cash_entry_repo.go` | `cashEntryRepository` | `CashEntryRepository` |
| `petty_cash_fund_repo.go` | `pettyCashFundRepository` | `PettyCashFundRepository` |
| `system_setting_repo.go` | `systemSettingRepository` | `SystemSettingRepository` |
| `message_template_repo.go` | `messageTemplateRepository` | `MessageTemplateRepository` |
| `sequence_counter_repo.go` | `sequenceCounterRepository` | `SequenceCounterRepository` |
| `audit_log_repo.go` | `auditLogRepository` | `AuditLogRepository` |
| `transactor.go` | `transactor` | `Transactor` |
| `registry.go` | `Registry` | — (aggregates all repos) |

## Registry

`Registry` adalah struct yang mengagregasi semua repository. Dibuat sekali di startup via `NewRepository(db)` dan diteruskan ke `cmd/server/main.go` untuk dependency injection ke service layer.

```go
pgRepo := postgres.NewRepository(db)
// pgRepo.CustomerRepo, pgRepo.InvoiceRepo, dst.
```

## hotspot_sale_repo.go

Repository untuk data penjualan voucher hotspot oleh sales agent.

### Tabel: `hotspot_sales`

### Methods

| Method | Deskripsi |
|--------|-----------|
| `Create(ctx, sale)` | Simpan satu record penjualan |
| `CreateBatch(ctx, sales)` | Bulk insert penjualan (digunakan setelah generate voucher) |
| `GetByID(ctx, id)` | Ambil penjualan by UUID |
| `List(ctx, filter, limit, offset)` | List penjualan dengan filter dan pagination |
| `Count(ctx, filter)` | Hitung total record (untuk pagination) |
| `ListByBatchCode(ctx, routerID, batchCode)` | List semua voucher dalam satu batch |
| `DeleteByBatchCode(ctx, routerID, batchCode)` | Hapus semua voucher dalam satu batch |
| `SumByAgentAndPeriod(ctx, agentID, from, to)` | Agregasi total penjualan agent dalam periode (untuk generate agent invoice) |

### Filter (`HotspotSaleFilter`)

| Field | Tipe | Deskripsi |
|-------|------|-----------|
| `RouterID` | `*uuid.UUID` | Filter by router |
| `SalesAgentID` | `*uuid.UUID` | Filter by agent |
| `Profile` | `string` | Filter by nama profil hotspot |
| `BatchCode` | `string` | Filter by kode batch voucher |
| `DateFrom` | `*time.Time` | Filter dari tanggal |
| `DateTo` | `*time.Time` | Filter sampai tanggal |

### Contoh Penggunaan

```go
// List penjualan by router dengan pagination
sales, err := repo.List(ctx, repository.HotspotSaleFilter{
    RouterID: &routerID,
}, 20, 0)

// Agregasi total untuk invoice agent (periode Jan 2024)
count, subtotal, sellingTotal, err := repo.SumByAgentAndPeriod(
    ctx, agentID,
    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
)

// Bulk insert setelah generate voucher mikhmon
err := repo.CreateBatch(ctx, []model.HotspotSale{...})
```

## Konvensi

- Semua query menggunakan `WithContext(ctx)` untuk timeout/cancellation.
- Soft delete menggunakan `deleted_at` (GORM default).
- Filter kosong (zero value) diabaikan — tidak ditambahkan ke WHERE clause.
- `SumByAgentAndPeriod` menggunakan `COALESCE(SUM(...), 0)` untuk menghindari NULL ketika tidak ada data.
