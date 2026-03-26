package seeder

import "context"

func (s *Seeder) seedSystemSettings(ctx context.Context) error {
	seeds := []struct {
		group, key, value, typ, label, description string
	}{
		// Base settings
		{"isolate", "pppoe_profile", "isolate", "string", "PPPoE Isolate Profile Name", ""},
		{"notification", "gowa_url", "", "string", "GoWA REST API URL", ""},
		{"notification", "gowa_sender", "", "string", "WhatsApp Sender Number", ""},
		{"billing", "due_days", "10", "integer", "Invoice Due Days", ""},
		{"billing", "reminder_intervals", "3,7,1", "string", "Reminder Intervals (days before due)", ""},
		{"billing", "late_fee_after_days", "0", "integer", "Late Fee After Days Overdue", ""},
		// Agent billing settings
		{"billing", "agent_default_billing_cycle", "monthly", "string", "Default Siklus Tagihan Agen", "Siklus tagihan default untuk agent baru: monthly atau weekly"},
		{"billing", "agent_default_billing_day_monthly", "1", "integer", "Default Hari Tagihan (Monthly)", "Tanggal invoice diproses untuk agent monthly (1-28)"},
		{"billing", "agent_default_billing_day_weekly", "1", "integer", "Default Hari Tagihan (Weekly)", "Hari dalam minggu untuk agent weekly (1=Senin, 7=Minggu)"},
	}

	for _, seed := range seeds {
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO system_settings (group_name, key_name, value, type, label, description)
			VALUES ($1::varchar, $2::varchar, NULLIF($3, '')::varchar, $4::varchar, $5::varchar, NULLIF($6, ''))
			ON CONFLICT (group_name, key_name) DO NOTHING
		`, seed.group, seed.key, seed.value, seed.typ, seed.label, seed.description)
		if err != nil {
			return err
		}
	}

	return nil
}
