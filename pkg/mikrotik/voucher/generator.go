package voucher

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"mikmongo/pkg/mikrotik/domain"
)

// Generator generates voucher codes
type Generator struct{}

// NewGenerator creates a new voucher generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateBatch generates a batch of vouchers.
// Comment is intentionally empty; the caller sets the correct format.
func (g *Generator) GenerateBatch(req *domain.VoucherGenerateRequest) []*domain.Voucher {
	vouchers := make([]*domain.Voucher, 0, req.Quantity)

	for i := 0; i < req.Quantity; i++ {
		var username, password string

		switch req.Mode {
		case "vc":
			username = g.generateVoucherCode(req.Prefix, req.NameLength, req.CharacterSet)
			password = username
		case "up":
			username = g.generateUsername(req.Prefix, req.NameLength, req.CharacterSet)
			password = g.generatePassword(req.NameLength)
		}

		vouchers = append(vouchers, &domain.Voucher{
			Username:  username,
			Password:  password,
			Profile:   req.Profile,
			Server:    req.Server,
			TimeLimit: req.TimeLimit,
			DataLimit: req.DataLimit,
		})
	}

	return vouchers
}

func (g *Generator) generateVoucherCode(prefix string, length int, charset string) string {
	chars := g.getCharset(charset)
	var sb strings.Builder
	sb.Grow(length + len(prefix))
	sb.WriteString(prefix)

	switch charset {
	case "num":
		sb.WriteString(g.randomString(chars, length))
	case "lower1", "upper1", "upplow1":
		letterLen := length - 2
		if letterLen < 1 {
			letterLen = 2
		}
		sb.WriteString(g.randomString(g.getLetters(charset), letterLen))
		sb.WriteString(g.randomString("0123456789", length-letterLen))
	default:
		sb.WriteString(g.randomString(chars, length))
	}

	return sb.String()
}

func (g *Generator) generateUsername(prefix string, length int, charset string) string {
	chars := g.getCharset(charset)
	var sb strings.Builder
	sb.Grow(length + len(prefix))
	sb.WriteString(prefix)
	sb.WriteString(g.randomString(chars, length))
	return sb.String()
}

func (g *Generator) generatePassword(length int) string {
	return g.randomString("0123456789", length)
}

func (g *Generator) randomString(chars string, length int) string {
	result := make([]byte, length)
	charLen := big.NewInt(int64(len(chars)))
	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, charLen)
		result[i] = chars[n.Int64()]
	}
	return string(result)
}

func (g *Generator) getCharset(charset string) string {
	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"

	switch charset {
	case "lower":
		return lower
	case "upper":
		return upper
	case "upplow":
		return lower + upper
	case "lower1":
		return lower + digits
	case "upper1":
		return upper + digits
	case "upplow1":
		return lower + upper + digits
	case "mix":
		return lower + digits
	case "mix1":
		return upper + digits
	case "mix2":
		return lower + upper + digits
	case "num":
		return digits
	default:
		return lower + digits
	}
}

func (g *Generator) getLetters(charset string) string {
	switch charset {
	case "lower1":
		return "abcdefghijklmnopqrstuvwxyz"
	case "upper1":
		return "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "upplow1":
		return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	default:
		return "abcdefghijklmnopqrstuvwxyz"
	}
}

// ParseDataLimit parses data limit string to bytes.
// Supports formats: 100M, 1G, 500K
func ParseDataLimit(limit string) int64 {
	if limit == "" {
		return 0
	}
	limit = strings.ToUpper(limit)
	var multiplier int64 = 1

	if strings.HasSuffix(limit, "G") {
		multiplier = 1073741824
		limit = strings.TrimSuffix(limit, "G")
	} else if strings.HasSuffix(limit, "M") {
		multiplier = 1048576
		limit = strings.TrimSuffix(limit, "M")
	} else if strings.HasSuffix(limit, "K") {
		multiplier = 1024
		limit = strings.TrimSuffix(limit, "K")
	}

	var value int64
	fmt.Sscanf(limit, "%d", &value)
	return value * multiplier
}
