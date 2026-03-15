package mikhmon

import (
	"math/rand"
	"strings"
	"time"

	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

// generatorRepository implements GeneratorRepository interface
type generatorRepository struct {
	rng *rand.Rand
}

// NewGeneratorRepository creates a new generator repository
func NewGeneratorRepository() GeneratorRepository {
	return &generatorRepository{
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// charSetMap defines character sets for different modes
var charSetMap = map[string]string{
	mikhmonDomain.CharSetLower:   "abcdefghijklmnopqrstuvwxyz",
	mikhmonDomain.CharSetUpper:   "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	mikhmonDomain.CharSetUpplow:  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	mikhmonDomain.CharSetLower1:  "abcdefghijklmnopqrstuvwxyz0123456789",
	mikhmonDomain.CharSetUpper1:  "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	mikhmonDomain.CharSetUpplow1: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	mikhmonDomain.CharSetMix:     "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	mikhmonDomain.CharSetMix1:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%",
	mikhmonDomain.CharSetMix2:    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*",
	mikhmonDomain.CharSetNumeric: "0123456789",
}

// GenerateUsername generates a random username
func (r *generatorRepository) GenerateUsername(config *mikhmonDomain.GeneratorConfig) string {
	chars := r.getCharSet(config.CharSet)
	username := r.generateString(chars, config.Length)

	if config.Prefix != "" {
		username = config.Prefix + username
	}

	return username
}

// GeneratePassword generates a random password
func (r *generatorRepository) GeneratePassword(config *mikhmonDomain.GeneratorConfig) string {
	// For password, use a mixed charset for better security
	chars := r.getCharSet(mikhmonDomain.CharSetUpplow1)
	return r.generateString(chars, config.Length)
}

// GeneratePair generates a username/password pair based on mode
func (r *generatorRepository) GeneratePair(mode string, config *mikhmonDomain.GeneratorConfig) *mikhmonDomain.GeneratorResult {
	username := r.GenerateUsername(config)

	var password string
	if mode == mikhmonDomain.VoucherModeVoucher {
		// Voucher mode: username = password
		password = username
	} else {
		// User/Password mode: different password
		password = r.GeneratePassword(config)
	}

	return &mikhmonDomain.GeneratorResult{
		Username: username,
		Password: password,
	}
}

// getCharSet returns the character set for a given mode
func (r *generatorRepository) getCharSet(mode string) string {
	if chars, ok := charSetMap[mode]; ok {
		return chars
	}
	// Default to alphanumeric
	return charSetMap[mikhmonDomain.CharSetUpplow1]
}

// generateString generates a random string from the given character set
func (r *generatorRepository) generateString(chars string, length int) string {
	if length <= 0 {
		length = 6 // Default length
	}

	var sb strings.Builder
	charLen := len(chars)

	for i := 0; i < length; i++ {
		sb.WriteByte(chars[r.rng.Intn(charLen)])
	}

	return sb.String()
}

// ShuffleString shuffles a string ( Fisher-Yates algorithm )
func (r *generatorRepository) ShuffleString(s string) string {
	bytes := []byte(s)
	n := len(bytes)

	for i := n - 1; i > 0; i-- {
		j := r.rng.Intn(i + 1)
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}

	return string(bytes)
}
