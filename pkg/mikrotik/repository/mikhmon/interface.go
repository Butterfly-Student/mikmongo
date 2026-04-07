package mikhmon

import (
	"context"

	"github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

// VoucherRepository defines the interface for voucher management
type VoucherRepository interface {
	GenerateBatch(ctx context.Context, req *mikhmon.VoucherGenerateRequest) (*mikhmon.VoucherBatch, error)
	GetAllVouchers(ctx context.Context) ([]*mikhmon.Voucher, error)
	GetVouchersByComment(ctx context.Context, comment string) ([]*mikhmon.Voucher, error)
	GetVouchersByCode(ctx context.Context, code string) ([]*mikhmon.Voucher, error)
	RemoveVoucherBatch(ctx context.Context, comment string) error
}

// ProfileRepository defines the interface for Mikhmon profile management
type ProfileRepository interface {
	CreateProfile(ctx context.Context, req *mikhmon.ProfileRequest) error
	UpdateProfile(ctx context.Context, id string, req *mikhmon.ProfileRequest) error
	GenerateOnLoginScript(data *mikhmon.OnLoginScriptData) string
}

// ReportRepository defines the interface for sales report management
type ReportRepository interface {
	AddReport(ctx context.Context, req *mikhmon.SalesReportRequest) error
	GetReportsByOwner(ctx context.Context, owner string) ([]*mikhmon.SalesReport, error)
	GetReportsByDay(ctx context.Context, day string) ([]*mikhmon.SalesReport, error)
	GetReportSummary(ctx context.Context, filter *mikhmon.ReportFilter) (*mikhmon.ReportSummary, error)
}

// GeneratorRepository defines the interface for user/voucher generation
type GeneratorRepository interface {
	GenerateUsername(config *mikhmon.GeneratorConfig) string
	GeneratePassword(config *mikhmon.GeneratorConfig) string
	GeneratePair(mode string, config *mikhmon.GeneratorConfig) *mikhmon.GeneratorResult
}

// ExpireRepository defines the interface for expire monitor management
type ExpireRepository interface {
	SetupExpireMonitor(ctx context.Context) error
	DisableExpireMonitor(ctx context.Context) error
	IsExpireMonitorEnabled(ctx context.Context) (bool, error)
	GenerateExpireMonitorScript() string
}

// Repository is the aggregator interface for all Mikhmon repositories
type Repository interface {
	Voucher() VoucherRepository
	Profile() ProfileRepository
	Report() ReportRepository
	Generator() GeneratorRepository
	Expire() ExpireRepository
}
