package seeder

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

func (s *Seeder) seedRouter(ctx context.Context) error {
	if s.cfg.RouterAddress == "" {
		return nil // skip when no address configured
	}

	name := s.cfg.RouterName
	if name == "" {
		name = "Router Utama"
	}
	apiPort := s.cfg.RouterAPIPort
	if apiPort == 0 {
		apiPort = 8728
	}
	username := s.cfg.RouterUsername
	if username == "" {
		username = "admin"
	}

	encPass, err := encryptAESGCM(s.cfg.RouterPassword, s.cfg.EncryptionKey)
	if err != nil {
		return err
	}

	// No unique constraint on address, use WHERE NOT EXISTS to avoid duplicates
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO mikrotik_routers (name, address, api_port, username, password_encrypted, is_master, is_active)
		SELECT $1::varchar, $2::varchar, $3::int, $4::varchar, $5::text, true, true
		WHERE NOT EXISTS (SELECT 1 FROM mikrotik_routers WHERE address = $2::varchar AND deleted_at IS NULL)
	`, name, s.cfg.RouterAddress, apiPort, username, encPass)
	return err
}

// encryptAESGCM mirrors the logic in internal/service/router_service.go
func encryptAESGCM(plaintext, key string) (string, error) {
	k := []byte(key)
	padded := make([]byte, 32)
	copy(padded, k)

	block, err := aes.NewCipher(padded)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
