package mikhmon

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	"github.com/Butterfly-Student/go-ros/repository/hotspot"
)

// voucherRepository implements VoucherRepository interface
type voucherRepository struct {
	client        *client.Client
	hotspotRepo   hotspot.Repository
	generatorRepo GeneratorRepository
}

// NewVoucherRepository creates a new voucher repository
func NewVoucherRepository(c *client.Client, hr hotspot.Repository, gr GeneratorRepository) VoucherRepository {
	return &voucherRepository{
		client:        c,
		hotspotRepo:   hr,
		generatorRepo: gr,
	}
}

// GenerateBatch generates a batch of vouchers
func (r *voucherRepository) GenerateBatch(ctx context.Context, req *mikhmonDomain.VoucherGenerateRequest) (*mikhmonDomain.VoucherBatch, error) {
	// Generate unique batch code
	code := generateBatchCode()

	batch := &mikhmonDomain.VoucherBatch{
		Code:      code,
		Quantity:  req.Quantity,
		Profile:   req.Profile,
		Server:    req.Server,
		TimeLimit: req.TimeLimit,
		DataLimit: req.DataLimit,
		Vouchers:  make([]mikhmonDomain.Voucher, 0, req.Quantity),
	}

	// Generate voucher comment format: vc-[CODE]-[DATE] or up-[CODE]-[DATE]
	date := time.Now().Format("01.02.06") // MM.DD.YY
	var comment string
	if req.Comment != "" {
		comment = fmt.Sprintf("%s-%s-%s-%s", req.Mode, code, date, req.Comment)
	} else {
		comment = fmt.Sprintf("%s-%s-%s", req.Mode, code, date)
	}

	// Generate config
	genConfig := &mikhmonDomain.GeneratorConfig{
		Length:  req.NameLength,
		Prefix:  req.Prefix,
		CharSet: req.CharSet,
	}

	// Generate vouchers
	for i := 0; i < req.Quantity; i++ {
		result := r.generatorRepo.GeneratePair(req.Mode, genConfig)

		voucher := mikhmonDomain.Voucher{
			Name:     result.Username,
			Password: result.Password,
			Profile:  req.Profile,
			Server:   req.Server,
			Comment:  comment,
			Code:     code,
			Mode:     req.Mode,
			Date:     date,
		}

		// Create user in MikroTik
		hotspotUser := &domain.HotspotUser{
			Name:            result.Username,
			Password:        result.Password,
			Profile:         req.Profile,
			Server:          req.Server,
			Comment:         comment,
			LimitUptime:     req.TimeLimit,
			LimitBytesTotal: parseDataLimit(req.DataLimit),
		}

		_, err := r.hotspotRepo.User().AddUser(ctx, hotspotUser)
		if err != nil {
			return nil, fmt.Errorf("failed to add voucher user %s: %w", result.Username, err)
		}

		batch.Vouchers = append(batch.Vouchers, voucher)
	}

	return batch, nil
}

// GetAllVouchers retrieves all hotspot users as vouchers
func (r *voucherRepository) GetAllVouchers(ctx context.Context) ([]*mikhmonDomain.Voucher, error) {
	users, err := r.hotspotRepo.User().GetUsers(ctx, "")
	if err != nil {
		return nil, err
	}

	vouchers := make([]*mikhmonDomain.Voucher, 0, len(users))
	for _, user := range users {
		voucher := &mikhmonDomain.Voucher{
			ID:       user.ID,
			Name:     user.Name,
			Password: user.Password,
			Profile:  user.Profile,
			Server:   user.Server,
			Comment:  user.Comment,
		}
		vouchers = append(vouchers, voucher)
	}

	return vouchers, nil
}

// GetVouchersByComment retrieves vouchers by comment
func (r *voucherRepository) GetVouchersByComment(ctx context.Context, comment string) ([]*mikhmonDomain.Voucher, error) {
	users, err := r.hotspotRepo.User().GetUsersByComment(ctx, comment)
	if err != nil {
		return nil, err
	}

	vouchers := make([]*mikhmonDomain.Voucher, 0, len(users))
	for _, user := range users {
		voucher := &mikhmonDomain.Voucher{
			ID:       user.ID,
			Name:     user.Name,
			Password: user.Password,
			Profile:  user.Profile,
			Server:   user.Server,
			Comment:  user.Comment,
		}
		vouchers = append(vouchers, voucher)
	}

	return vouchers, nil
}

// GetVouchersByCode retrieves vouchers by batch code
func (r *voucherRepository) GetVouchersByCode(ctx context.Context, code string) ([]*mikhmonDomain.Voucher, error) {
	// Search for vouchers with comment containing the code
	// Format: vc-[CODE]-[DATE] or up-[CODE]-[DATE]
	commentPattern := fmt.Sprintf("-%s-", code)

	// Get all users and filter by code
	users, err := r.hotspotRepo.User().GetUsers(ctx, "")
	if err != nil {
		return nil, err
	}

	vouchers := make([]*mikhmonDomain.Voucher, 0)
	for _, user := range users {
		if contains(user.Comment, commentPattern) {
			voucher := &mikhmonDomain.Voucher{
				ID:       user.ID,
				Name:     user.Name,
				Password: user.Password,
				Profile:  user.Profile,
				Server:   user.Server,
				Comment:  user.Comment,
			}
			vouchers = append(vouchers, voucher)
		}
	}

	return vouchers, nil
}

// RemoveVoucherBatch removes all vouchers in a batch
func (r *voucherRepository) RemoveVoucherBatch(ctx context.Context, comment string) error {
	vouchers, err := r.GetVouchersByComment(ctx, comment)
	if err != nil {
		return err
	}

	for _, voucher := range vouchers {
		err := r.hotspotRepo.User().RemoveUser(ctx, voucher.ID)
		if err != nil {
			return fmt.Errorf("failed to remove voucher %s: %w", voucher.Name, err)
		}
	}

	return nil
}

// Helper functions

func generateBatchCode() string {
	// Generate 3-digit random code
	return fmt.Sprintf("%03d", rand.Intn(900)+100)
}

func parseDataLimit(limit string) int64 {
	if limit == "" {
		return 0
	}

	// Parse format like "1G", "500M", "1024K"
	var value int64
	var unit string

	fmt.Sscanf(limit, "%d%s", &value, &unit)

	switch unit {
	case "G", "g":
		return value * 1024 * 1024 * 1024
	case "M", "m":
		return value * 1024 * 1024
	case "K", "k":
		return value * 1024
	default:
		return value
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				findSubstring(s, substr) >= 0))
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
