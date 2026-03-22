package seeder

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

func (s *Seeder) seedRouters(ctx context.Context) error {
	routers := []struct {
		name     string
		address  string
		apiPort  int
		username string
		password string
		isMaster bool
	}{
		{"Router Utama", "192.168.233.1", 8728, "admin", "r00t", true},
		{"Router 2", "192.168.27.1", 8728, "admin", "r00t", false},
	}

	for _, r := range routers {
		encPass, err := encryptAESGCM(r.password, s.cfg.EncryptionKey)
		if err != nil {
			return fmt.Errorf("encrypt router %s: %w", r.address, err)
		}

		_, err = s.db.ExecContext(ctx, `
			INSERT INTO mikrotik_routers (name, address, api_port, username, password_encrypted, is_master, is_active)
			SELECT $1::varchar, $2::varchar, $3::int, $4::varchar, $5::text, $6::boolean, true
			WHERE NOT EXISTS (SELECT 1 FROM mikrotik_routers WHERE address = $2::varchar AND deleted_at IS NULL)
		`, r.name, r.address, r.apiPort, r.username, encPass, r.isMaster)
		if err != nil {
			return fmt.Errorf("insert router %s: %w", r.address, err)
		}
	}
	return nil
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
