package seeder

import "context"

func (s *Seeder) seedSystemSettings(ctx context.Context) error {
	seeds := []struct {
		group, key, value, typ, label string
	}{
		{"isolate", "pppoe_profile", "isolate", "string", "PPPoE Isolate Profile Name"},
		{"notification", "gowa_url", "", "string", "GoWA REST API URL"},
		{"notification", "gowa_sender", "", "string", "WhatsApp Sender Number"},
		{"billing", "due_days", "10", "integer", "Invoice Due Days"},
		{"billing", "reminder_intervals", "3,7,1", "string", "Reminder Intervals (days before due)"},
		{"billing", "late_fee_after_days", "0", "integer", "Late Fee After Days Overdue"},
	}

	for _, seed := range seeds {
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO system_settings (group_name, key_name, value, type, label)
			VALUES ($1::varchar, $2::varchar, NULLIF($3, '')::varchar, $4::varchar, $5::varchar)
			ON CONFLICT (group_name, key_name) DO NOTHING
		`, seed.group, seed.key, seed.value, seed.typ, seed.label)
		if err != nil {
			return err
		}
	}

	return nil
}
