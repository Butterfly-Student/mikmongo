package seeder

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func (s *Seeder) seedCustomers(ctx context.Context) error {
	portalHash, err := bcrypt.GenerateFromPassword([]byte("portal123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt portal password: %w", err)
	}

	customers := []struct {
		code     string
		name     string
		phone    string
		email    string
		username string
	}{
		{"CST-00001", "Budi Santoso", "+6281200000001", "budi.santoso@example.com", "budi-santoso"},
		{"CST-00002", "Siti Rahayu", "+6281200000002", "siti.rahayu@example.com", "siti-rahayu"},
		{"CST-00003", "Ahmad Fauzi", "+6281200000003", "ahmad.fauzi@example.com", "ahmad-fauzi"},
		{"CST-00004", "Dewi Lestari", "+6281200000004", "dewi.lestari@example.com", "dewi-lestari"},
		{"CST-00005", "Rudi Hartono", "+6281200000005", "rudi.hartono@example.com", "rudi-hartono"},
	}

	for _, c := range customers {
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO customers (customer_code, full_name, phone, email, username, portal_password_hash, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, true)
			ON CONFLICT (customer_code) DO NOTHING
		`, c.code, c.name, c.phone, c.email, c.username, string(portalHash))
		if err != nil {
			return fmt.Errorf("insert customer %s: %w", c.code, err)
		}
	}

	// Advance the sequence counter so API-created customers start from CST-00006
	_, err = s.db.ExecContext(ctx, `
		UPDATE sequence_counters
		SET last_number = GREATEST(last_number, 5)
		WHERE name = 'customer_code'
	`)
	return err
}
