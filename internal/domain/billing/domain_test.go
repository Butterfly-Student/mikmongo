package billing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mikmongo/internal/model"
)

func TestCalculateTax(t *testing.T) {
	d := NewDomain()

	t.Run("normal 11%", func(t *testing.T) {
		got := d.CalculateTax(100_000, 0.11)
		assert.Equal(t, 11_000.0, got)
	})

	t.Run("zero tax rate", func(t *testing.T) {
		got := d.CalculateTax(100_000, 0)
		assert.Equal(t, 0.0, got)
	})

	t.Run("zero subtotal", func(t *testing.T) {
		got := d.CalculateTax(0, 0.11)
		assert.Equal(t, 0.0, got)
	})

	t.Run("rounding to 2 decimals", func(t *testing.T) {
		// 99999 * 0.11 = 10999.89
		got := d.CalculateTax(99_999, 0.11)
		assert.Equal(t, 10999.89, got)
	})

	t.Run("fractional result rp199999.5 @ 11%", func(t *testing.T) {
		// 199999.5 * 0.11 = 21999.945 → rounded to 2 decimals
		got := d.CalculateTax(199_999.5, 0.11)
		assert.InDelta(t, 21999.95, got, 0.01)
	})

	t.Run("large amount 5M @ 11%", func(t *testing.T) {
		got := d.CalculateTax(5_000_000, 0.11)
		assert.Equal(t, 550_000.0, got)
	})
}

func TestCalculateTotal(t *testing.T) {
	d := NewDomain()

	t.Run("subtotal + tax - discount + lateFee", func(t *testing.T) {
		got := d.CalculateTotal(100_000, 11_000, 5_000, 2_000)
		assert.Equal(t, 108_000.0, got)
	})

	t.Run("zero components", func(t *testing.T) {
		got := d.CalculateTotal(0, 0, 0, 0)
		assert.Equal(t, 0.0, got)
	})

	t.Run("no discount no late fee", func(t *testing.T) {
		got := d.CalculateTotal(100_000, 11_000, 0, 0)
		assert.Equal(t, 111_000.0, got)
	})

	t.Run("no fp drift with fractional subtotal", func(t *testing.T) {
		// 99999.99 + 10999.999 = 110999.989 → rounded to 2 decimals
		got := d.CalculateTotal(99_999.99, 10_999.999, 0, 0)
		assert.InDelta(t, 110_999.99, got, 0.01)
	})
}

func TestCalculateProration(t *testing.T) {
	d := NewDomain()

	t.Run("full month", func(t *testing.T) {
		got := d.CalculateProration(300_000, 31, 31)
		assert.Equal(t, 300_000.0, got)
	})

	t.Run("half month activation (day 16 of 31)", func(t *testing.T) {
		got := d.CalculateProration(300_000, 31, 16)
		assert.InDelta(t, 154838.71, got, 0.01)
	})

	t.Run("day 1 activation (full month billed)", func(t *testing.T) {
		got := d.CalculateProration(300_000, 30, 30)
		assert.Equal(t, 300_000.0, got)
	})

	t.Run("last day activation (nearly zero)", func(t *testing.T) {
		got := d.CalculateProration(300_000, 31, 1)
		assert.InDelta(t, 9677.42, got, 0.01)
	})

	t.Run("zero totalDays", func(t *testing.T) {
		got := d.CalculateProration(300_000, 0, 15)
		assert.Equal(t, 0.0, got)
	})
}

func TestCalculateLateFee(t *testing.T) {
	d := NewDomain()

	t.Run("1 day overdue → 1%", func(t *testing.T) {
		got := d.CalculateLateFee(100_000, 1)
		assert.Equal(t, 1_000.0, got)
	})

	t.Run("10 days overdue → 10% (cap)", func(t *testing.T) {
		got := d.CalculateLateFee(100_000, 10)
		assert.Equal(t, 10_000.0, got)
	})

	t.Run("20 days overdue → still 10% (capped)", func(t *testing.T) {
		got := d.CalculateLateFee(100_000, 20)
		assert.Equal(t, 10_000.0, got)
	})

	t.Run("zero days overdue", func(t *testing.T) {
		got := d.CalculateLateFee(100_000, 0)
		assert.Equal(t, 0.0, got)
	})

	t.Run("negative days overdue", func(t *testing.T) {
		got := d.CalculateLateFee(100_000, -1)
		assert.Equal(t, 0.0, got)
	})
}

func TestIsOverdue(t *testing.T) {
	d := NewDomain()
	past := time.Now().AddDate(0, 0, -5)
	future := time.Now().AddDate(0, 0, 5)
	now := time.Now()

	t.Run("due date in past → overdue", func(t *testing.T) {
		inv := &model.Invoice{Status: "unpaid", DueDate: past}
		assert.True(t, d.IsOverdue(inv, now))
	})

	t.Run("due date today → not overdue", func(t *testing.T) {
		inv := &model.Invoice{Status: "unpaid", DueDate: now}
		assert.False(t, d.IsOverdue(inv, now))
	})

	t.Run("due date in future → not overdue", func(t *testing.T) {
		inv := &model.Invoice{Status: "unpaid", DueDate: future}
		assert.False(t, d.IsOverdue(inv, now))
	})

	t.Run("status paid → not overdue even if past due", func(t *testing.T) {
		inv := &model.Invoice{Status: "paid", DueDate: past}
		assert.False(t, d.IsOverdue(inv, now))
	})

	t.Run("status cancelled → not overdue", func(t *testing.T) {
		inv := &model.Invoice{Status: "cancelled", DueDate: past}
		assert.False(t, d.IsOverdue(inv, now))
	})

	t.Run("status refunded → not overdue", func(t *testing.T) {
		inv := &model.Invoice{Status: "refunded", DueDate: past}
		assert.False(t, d.IsOverdue(inv, now))
	})
}

func TestDaysOverdue(t *testing.T) {
	d := NewDomain()
	now := time.Now()

	t.Run("5 days overdue", func(t *testing.T) {
		dueDate := now.AddDate(0, 0, -5)
		inv := &model.Invoice{Status: "unpaid", DueDate: dueDate}
		assert.Equal(t, 5, d.DaysOverdue(inv, now))
	})

	t.Run("not overdue → 0", func(t *testing.T) {
		dueDate := now.AddDate(0, 0, 5)
		inv := &model.Invoice{Status: "unpaid", DueDate: dueDate}
		assert.Equal(t, 0, d.DaysOverdue(inv, now))
	})
}

func TestShouldSuspendForNonPayment(t *testing.T) {
	d := NewDomain()
	now := time.Now()

	t.Run("within grace period → false", func(t *testing.T) {
		dueDate := now.AddDate(0, 0, -2)
		inv := &model.Invoice{Status: "unpaid", DueDate: dueDate}
		assert.False(t, d.ShouldSuspendForNonPayment(inv, now, 3))
	})

	t.Run("beyond grace period → true", func(t *testing.T) {
		dueDate := now.AddDate(0, 0, -10)
		inv := &model.Invoice{Status: "unpaid", DueDate: dueDate}
		assert.True(t, d.ShouldSuspendForNonPayment(inv, now, 3))
	})

	t.Run("already paid → false", func(t *testing.T) {
		dueDate := now.AddDate(0, 0, -10)
		inv := &model.Invoice{Status: "paid", DueDate: dueDate}
		assert.False(t, d.ShouldSuspendForNonPayment(inv, now, 3))
	})
}

func TestShouldSendReminder(t *testing.T) {
	d := NewDomain()
	now := time.Now()

	t.Run("no previous reminder → true", func(t *testing.T) {
		inv := &model.Invoice{LastReminderSent: nil}
		assert.True(t, d.ShouldSendReminder(inv, now, 3))
	})

	t.Run("last reminder more than interval ago → true", func(t *testing.T) {
		lastSent := now.AddDate(0, 0, -4)
		inv := &model.Invoice{LastReminderSent: &lastSent}
		assert.True(t, d.ShouldSendReminder(inv, now, 3))
	})

	t.Run("last reminder sent recently (within interval) → false", func(t *testing.T) {
		lastSent := now.AddDate(0, 0, -1)
		inv := &model.Invoice{LastReminderSent: &lastSent}
		assert.False(t, d.ShouldSendReminder(inv, now, 3))
	})
}

func TestInvoiceStatusFromAmounts(t *testing.T) {
	d := NewDomain()

	t.Run("paid == 0 → unpaid", func(t *testing.T) {
		assert.Equal(t, "unpaid", d.InvoiceStatusFromAmounts(100_000, 0))
	})

	t.Run("paid < total → partial", func(t *testing.T) {
		assert.Equal(t, "partial", d.InvoiceStatusFromAmounts(100_000, 50_000))
	})

	t.Run("paid == total → paid", func(t *testing.T) {
		assert.Equal(t, "paid", d.InvoiceStatusFromAmounts(100_000, 100_000))
	})

	t.Run("paid > total → paid (overpay treated as paid)", func(t *testing.T) {
		assert.Equal(t, "paid", d.InvoiceStatusFromAmounts(100_000, 150_000))
	})
}

func TestClampBillingDay(t *testing.T) {
	d := NewDomain()

	t.Run("day 31 in February (28 days) → 28", func(t *testing.T) {
		got := d.ClampBillingDay(31, 2023, time.February)
		assert.Equal(t, 28, got)
	})

	t.Run("day 31 in April (30 days) → 30", func(t *testing.T) {
		got := d.ClampBillingDay(31, 2023, time.April)
		assert.Equal(t, 30, got)
	})

	t.Run("day 1 always stays 1", func(t *testing.T) {
		got := d.ClampBillingDay(1, 2023, time.February)
		assert.Equal(t, 1, got)
	})

	t.Run("day 28 in March (31 days) → 28", func(t *testing.T) {
		got := d.ClampBillingDay(28, 2023, time.March)
		assert.Equal(t, 28, got)
	})

	t.Run("leap year Feb 29 with day 31 → 29", func(t *testing.T) {
		got := d.ClampBillingDay(31, 2024, time.February)
		assert.Equal(t, 29, got)
	})
}

func TestResolveGracePeriod(t *testing.T) {
	d := NewDomain()

	t.Run("subscription value takes priority", func(t *testing.T) {
		subGrace := 5
		got := d.ResolveGracePeriod(&subGrace, 7)
		assert.Equal(t, 5, got)
	})

	t.Run("subscription override of 0 is respected", func(t *testing.T) {
		subGrace := 0
		got := d.ResolveGracePeriod(&subGrace, 7)
		assert.Equal(t, 0, got)
	})

	t.Run("falls back to profile value", func(t *testing.T) {
		got := d.ResolveGracePeriod(nil, 7)
		assert.Equal(t, 7, got)
	})

	t.Run("falls back to default (3 days)", func(t *testing.T) {
		got := d.ResolveGracePeriod(nil, 0)
		assert.Equal(t, 3, got)
	})
}

func TestResolveBillingDay(t *testing.T) {
	d := NewDomain()

	t.Run("subscription value takes priority", func(t *testing.T) {
		subDay := 15
		profileDay := 1
		got := d.ResolveBillingDay(&subDay, &profileDay)
		assert.Equal(t, 15, got)
	})

	t.Run("falls back to profile value", func(t *testing.T) {
		profileDay := 20
		got := d.ResolveBillingDay(nil, &profileDay)
		assert.Equal(t, 20, got)
	})

	t.Run("falls back to default (1st)", func(t *testing.T) {
		got := d.ResolveBillingDay(nil, nil)
		assert.Equal(t, 1, got)
	})

	t.Run("subscription day 0 treated as unset, falls back to profile", func(t *testing.T) {
		subDay := 0
		profileDay := 10
		got := d.ResolveBillingDay(&subDay, &profileDay)
		assert.Equal(t, 10, got)
	})
}

func TestGetBillingPeriod(t *testing.T) {
	d := NewDomain()
	base := time.Date(2024, time.January, 15, 0, 0, 0, 0, time.UTC)

	t.Run("monthly cycle: correct start/end dates", func(t *testing.T) {
		start, end := d.GetBillingPeriod(base, "monthly")
		assert.Equal(t, base, start)
		assert.Equal(t, base.AddDate(0, 1, 0), end)
	})

	t.Run("daily cycle", func(t *testing.T) {
		start, end := d.GetBillingPeriod(base, "daily")
		assert.Equal(t, base, start)
		assert.Equal(t, base.AddDate(0, 0, 1), end)
	})

	t.Run("weekly cycle", func(t *testing.T) {
		start, end := d.GetBillingPeriod(base, "weekly")
		assert.Equal(t, base, start)
		assert.Equal(t, base.AddDate(0, 0, 7), end)
	})

	t.Run("yearly cycle", func(t *testing.T) {
		start, end := d.GetBillingPeriod(base, "yearly")
		assert.Equal(t, base, start)
		assert.Equal(t, base.AddDate(1, 0, 0), end)
	})

	t.Run("unknown cycle defaults to monthly", func(t *testing.T) {
		start, end := d.GetBillingPeriod(base, "biennial")
		assert.Equal(t, base, start)
		assert.Equal(t, base.AddDate(0, 1, 0), end)
	})
}
