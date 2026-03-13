// Package billing contains billing domain logic
package billing

import (
	"math"
	"mikmongo/internal/model"
	"time"
)

// Domain represents billing business logic
type Domain struct{}

// NewDomain creates a new billing domain
func NewDomain() *Domain {
	return &Domain{}
}

// CalculateTax returns subtotal * taxRate rounded to 2 decimal places
func (d *Domain) CalculateTax(subtotal, taxRate float64) float64 {
	return math.Round(subtotal*taxRate*100) / 100
}

// CalculateTotal returns subtotal + tax - discount + lateFee
func (d *Domain) CalculateTotal(subtotal, tax, discount, lateFee float64) float64 {
	return math.Round((subtotal+tax-discount+lateFee)*100) / 100
}

// CalculateProration returns monthlyPrice × (billedDays / totalDaysInMonth)
func (d *Domain) CalculateProration(monthlyPrice float64, totalDays, billedDays int) float64 {
	if totalDays <= 0 {
		return 0
	}
	return math.Round(monthlyPrice*float64(billedDays)/float64(totalDays)*100) / 100
}

// CalculateLateFee returns 1% per day overdue, capped at 10% of totalAmount
func (d *Domain) CalculateLateFee(totalAmount float64, daysOverdue int) float64 {
	if daysOverdue <= 0 {
		return 0
	}
	fee := totalAmount * 0.01 * float64(daysOverdue)
	max := totalAmount * 0.10
	if fee > max {
		fee = max
	}
	return math.Round(fee*100) / 100
}

// IsOverdue returns true if now > DueDate and status is not paid/cancelled/refunded
func (d *Domain) IsOverdue(inv *model.Invoice, now time.Time) bool {
	switch inv.Status {
	case "paid", "cancelled", "refunded":
		return false
	}
	return now.After(inv.DueDate)
}

// DaysOverdue returns the number of days past the DueDate (0 if not overdue)
func (d *Domain) DaysOverdue(inv *model.Invoice, now time.Time) int {
	if !d.IsOverdue(inv, now) {
		return 0
	}
	return int(now.Sub(inv.DueDate).Hours() / 24)
}

// ShouldSuspendForNonPayment returns true if invoice is overdue and now > DueDate + graceDays
func (d *Domain) ShouldSuspendForNonPayment(inv *model.Invoice, now time.Time, graceDays int) bool {
	if !d.IsOverdue(inv, now) {
		return false
	}
	suspendAt := inv.DueDate.AddDate(0, 0, graceDays)
	return now.After(suspendAt)
}

// ShouldSendReminder returns true if enough time has passed since the last reminder
func (d *Domain) ShouldSendReminder(inv *model.Invoice, now time.Time, intervalDays int) bool {
	if inv.LastReminderSent == nil {
		return true
	}
	nextReminder := inv.LastReminderSent.AddDate(0, 0, intervalDays)
	return now.After(nextReminder) || now.Equal(nextReminder)
}

// InvoiceStatusFromAmounts determines invoice status from paid vs total amounts
func (d *Domain) InvoiceStatusFromAmounts(total, paid float64) string {
	if paid <= 0 {
		return "unpaid"
	}
	if paid >= total {
		return "paid"
	}
	return "partial"
}

// ClampBillingDay returns billingDay clamped to the last day of the given month.
// Handles Feb (28/29), months with 30 days, etc.
func (d *Domain) ClampBillingDay(billingDay int, year int, month time.Month) int {
	lastDay := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
	if billingDay > lastDay {
		return lastDay
	}
	return billingDay
}

// ResolveGracePeriod returns the effective grace period using priority chain:
// subscription override > profile default > fallback 3
func (d *Domain) ResolveGracePeriod(subGrace *int, profileGrace int) int {
	if subGrace != nil && *subGrace >= 0 {
		return *subGrace
	}
	if profileGrace > 0 {
		return profileGrace
	}
	return 3
}

// ResolveBillingDay returns the effective billing day using priority chain:
// subscription override > profile default > fallback 1
func (d *Domain) ResolveBillingDay(subDay *int, profileDay *int) int {
	if subDay != nil && *subDay > 0 {
		return *subDay
	}
	if profileDay != nil && *profileDay > 0 {
		return *profileDay
	}
	return 1
}

// GetBillingPeriod calculates start and end date from a given time and billing cycle
func (d *Domain) GetBillingPeriod(from time.Time, cycle string) (start, end time.Time) {
	start = from
	switch cycle {
	case "daily":
		end = from.AddDate(0, 0, 1)
	case "weekly":
		end = from.AddDate(0, 0, 7)
	case "yearly":
		end = from.AddDate(1, 0, 0)
	default: // monthly
		end = from.AddDate(0, 1, 0)
	}
	return start, end
}
