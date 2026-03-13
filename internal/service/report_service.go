package service

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// SubscriptionStats holds subscription counts by status
type SubscriptionStats struct {
	Active    int64 `json:"active"`
	Isolated  int64 `json:"isolated"`
	Suspended int64 `json:"suspended"`
	Pending   int64 `json:"pending"`
	Total     int64 `json:"total"`
}

// ReportSummary holds summary statistics for a given period
type ReportSummary struct {
	PeriodStart     time.Time         `json:"period_start"`
	PeriodEnd       time.Time         `json:"period_end"`
	TotalRevenue    float64           `json:"total_revenue"`
	TotalInvoiced   float64           `json:"total_invoiced"`
	TotalInvoices   int64             `json:"total_invoices"`
	PaidInvoices    int64             `json:"paid_invoices"`
	UnpaidInvoices  int64             `json:"unpaid_invoices"`
	OverdueInvoices int64             `json:"overdue_invoices"`
	TotalPayments   int64             `json:"total_payments"`
	TotalCustomers  int64             `json:"total_customers"`
	ActiveCustomers int64             `json:"active_customers"`
	NewCustomers    int64             `json:"new_customers"`
	ActiveSubs      int64             `json:"active_subscriptions"`
	IsolatedSubs    int64             `json:"isolated_subscriptions"`
	SuspendedSubs   int64             `json:"suspended_subscriptions"`
	Subscriptions   SubscriptionStats `json:"subscriptions"`
}

// SubscriptionReportItem represents a single subscription in reports
type SubscriptionReportItem struct {
	ID           string     `json:"id"`
	CustomerCode string     `json:"customer_code"`
	CustomerName string     `json:"customer_name"`
	Username     string     `json:"username"`
	PlanName     string     `json:"plan_name"`
	Status       string     `json:"status"`
	StaticIP     *string    `json:"static_ip,omitempty"`
	ActivatedAt  *time.Time `json:"activated_at,omitempty"`
	BillingDay   *int       `json:"billing_day,omitempty"`
	AutoIsolate  bool       `json:"auto_isolate"`
	MonthlyPrice float64    `json:"monthly_price"`
}

// ReportService provides reporting functionality
type ReportService struct {
	db *gorm.DB
}

// NewReportService creates a new report service
func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{db: db}
}

// countSubsByStatus counts subscriptions for a given status
func (s *ReportService) countSubsByStatus(ctx context.Context, status string) int64 {
	var count int64
	s.db.WithContext(ctx).Table("subscriptions").
		Where("status = ? AND deleted_at IS NULL", status).
		Count(&count)
	return count
}

// GetSummary returns a summary report for the given date range
func (s *ReportService) GetSummary(ctx context.Context, from, to time.Time) (*ReportSummary, error) {
	summary := &ReportSummary{
		PeriodStart: from,
		PeriodEnd:   to,
	}

	// Revenue (confirmed payments in period)
	s.db.WithContext(ctx).
		Table("payments").
		Where("status = ? AND payment_date BETWEEN ? AND ? AND deleted_at IS NULL", "confirmed", from, to).
		Select("COALESCE(COUNT(*), 0) as total_payments, COALESCE(SUM(amount), 0) as total_revenue").
		Row().Scan(&summary.TotalPayments, &summary.TotalRevenue)

	// Invoices in period
	s.db.WithContext(ctx).
		Table("invoices").
		Where("issue_date BETWEEN ? AND ? AND deleted_at IS NULL", from, to).
		Select("COALESCE(COUNT(*), 0) as total, COALESCE(SUM(total_amount), 0) as amount").
		Row().Scan(&summary.TotalInvoices, &summary.TotalInvoiced)

	// Paid invoices
	s.db.WithContext(ctx).
		Table("invoices").
		Where("status = ? AND issue_date BETWEEN ? AND ? AND deleted_at IS NULL", "paid", from, to).
		Count(&summary.PaidInvoices)

	// Unpaid invoices
	s.db.WithContext(ctx).
		Table("invoices").
		Where("status IN (?, ?) AND issue_date BETWEEN ? AND ? AND deleted_at IS NULL", "unpaid", "partial", from, to).
		Count(&summary.UnpaidInvoices)

	// Overdue invoices
	s.db.WithContext(ctx).
		Table("invoices").
		Where("status IN (?, ?, ?) AND due_date < ? AND deleted_at IS NULL", "unpaid", "partial", "overdue", time.Now()).
		Count(&summary.OverdueInvoices)

	// Customer stats
	s.db.WithContext(ctx).Table("customers").Where("deleted_at IS NULL").Count(&summary.TotalCustomers)
	s.db.WithContext(ctx).Table("customers").Where("status = ? AND deleted_at IS NULL", "active").Count(&summary.ActiveCustomers)
	s.db.WithContext(ctx).Table("customers").Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", from, to).Count(&summary.NewCustomers)

	// Subscription breakdown (PPPoE only)
	summary.Subscriptions.Active = s.countSubsByStatus(ctx, "active")
	summary.Subscriptions.Isolated = s.countSubsByStatus(ctx, "isolated")
	summary.Subscriptions.Suspended = s.countSubsByStatus(ctx, "suspended")
	summary.Subscriptions.Pending = s.countSubsByStatus(ctx, "pending")
	summary.Subscriptions.Total = summary.Subscriptions.Active + summary.Subscriptions.Isolated + summary.Subscriptions.Suspended + summary.Subscriptions.Pending

	// Backward-compatible aggregates
	summary.ActiveSubs = summary.Subscriptions.Active
	summary.IsolatedSubs = summary.Subscriptions.Isolated
	summary.SuspendedSubs = summary.Subscriptions.Suspended

	return summary, nil
}

// GetSubscriptionReport returns a filtered, paginated subscription report
func (s *ReportService) GetSubscriptionReport(ctx context.Context, from, to time.Time, limit, offset int) ([]SubscriptionReportItem, int64, error) {
	var items []SubscriptionReportItem
	var count int64

	query := s.db.WithContext(ctx).
		Table("subscriptions s").
		Select(`s.id, c.customer_code, c.full_name as customer_name,
                s.username, bp.name as plan_name,
                s.status, s.static_ip, s.activated_at,
                s.billing_day, s.auto_isolate, bp.price_monthly as monthly_price`).
		Joins("JOIN customers c ON c.id = s.customer_id").
		Joins("JOIN bandwidth_profiles bp ON bp.id = s.plan_id").
		Where("s.deleted_at IS NULL AND c.deleted_at IS NULL")

	if !from.IsZero() {
		query = query.Where("s.created_at >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("s.created_at <= ?", to)
	}

	// Count
	countQuery := query.Session(&gorm.Session{})
	countQuery.Count(&count)

	// Fetch with pagination
	query.Order("s.created_at DESC").Limit(limit).Offset(offset).Scan(&items)

	return items, count, nil
}
