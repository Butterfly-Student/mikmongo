package seeder

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

func (s *Seeder) seedUser(ctx context.Context) error {
	email := s.cfg.AdminEmail
	if email == "" {
		email = "admin@mikmongo.local"
	}
	password := s.cfg.AdminPassword
	if password == "" {
		password = "admin123"
	}
	name := s.cfg.AdminName
	if name == "" {
		name = "Super Admin"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	phone := s.cfg.AdminPhone
	if phone == "" {
		phone = "+6280000000000"
	}

	_, err = s.db.ExecContext(ctx, `
		INSERT INTO users (full_name, email, phone, password_hash, role, is_active)
		VALUES ($1, $2, $3, $4, 'superadmin', true)
		ON CONFLICT (email) DO NOTHING
	`, name, email, phone, string(hash))
	return err
}
