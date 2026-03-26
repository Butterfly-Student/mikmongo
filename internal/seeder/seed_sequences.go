package seeder

import "context"

func (s *Seeder) seedSequenceCounters(ctx context.Context) error {
	counters := []struct {
		name, prefix string
		padding      int
	}{
		{"invoice_number", "INV", 6},
		{"payment_number", "PAY", 6},
		{"customer_code", "CST", 5},
		{"cash_entry_number", "KAS", 6},
	}

	for _, c := range counters {
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO sequence_counters (name, prefix, padding, last_number)
			VALUES ($1::varchar, $2::varchar, $3::int, 0)
			ON CONFLICT (name) DO NOTHING
		`, c.name, c.prefix, c.padding)
		if err != nil {
			return err
		}
	}

	return nil
}
