package service

import (
	"context"
	"encoding/json"
	"fmt"

	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/pkg/mikrotik"
	mkdomain "mikmongo/pkg/mikrotik/domain"

	"github.com/google/uuid"
)

// BandwidthProfileService handles bandwidth profile business logic
type BandwidthProfileService struct {
	profileRepo repository.BandwidthProfileRepository
	routerSvc   *RouterService
}

// NewBandwidthProfileService creates a new bandwidth profile service
func NewBandwidthProfileService(profileRepo repository.BandwidthProfileRepository, routerSvc *RouterService) *BandwidthProfileService {
	return &BandwidthProfileService{
		profileRepo: profileRepo,
		routerSvc:   routerSvc,
	}
}

// Create creates a new bandwidth profile and syncs PPP profile to MikroTik
func (s *BandwidthProfileService) Create(ctx context.Context, profile *model.BandwidthProfile) error {
	// Connect to MikroTik first
	routerID, err := uuid.Parse(profile.RouterID)
	if err != nil {
		return err
	}

	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return err
	}

	// Sync PPP profile to MikroTik
	if err := s.syncPPPProfile(ctx, mt, profile); err != nil {
		return fmt.Errorf("failed to sync to mikrotik: %w", err)
	}

	// Save to database only if MikroTik succeeded
	if err := s.profileRepo.Create(ctx, profile); err != nil {
		// Rollback: remove from MikroTik
		existing, _ := mt.PPP.GetProfileByName(ctx, profile.Name)
		if existing != nil {
			_ = mt.PPP.RemoveProfile(ctx, existing.ID)
		}
		return err
	}

	return nil
}

// syncPPPProfile syncs PPP profile to MikroTik
func (s *BandwidthProfileService) syncPPPProfile(ctx context.Context, mt *mikrotik.Client, profile *model.BandwidthProfile) error {
	var config model.PPPProfileConfig
	if len(profile.MikrotikConfig) > 0 {
		if err := json.Unmarshal(profile.MikrotikConfig, &config); err != nil {
			return fmt.Errorf("failed to unmarshal PPP config: %w", err)
		}
	}

	pppProfile := &mkdomain.PPPProfile{
		Name:           profile.Name,
		LocalAddress:   config.LocalAddress,
		RemoteAddress:  config.RemoteAddress,
		DNSServer:      config.DNSServer,
		SessionTimeout: config.SessionTimeout,
		IdleTimeout:    config.IdleTimeout,
		RateLimit:      config.RateLimit,
		UseCompression: config.UseCompression,
		UseEncryption:  config.UseEncryption,
		OnlyOne:        config.OnlyOne,
		ChangeTCPMSS:   config.ChangeTCPMSS,
		Bridge:         config.Bridge,
		AddressList:    config.AddressList,
	}

	// Cek apakah profile sudah ada
	existing, _ := mt.PPP.GetProfileByName(ctx, profile.Name)
	if existing != nil {
		return mt.PPP.UpdateProfile(ctx, existing.ID, pppProfile)
	}
	return mt.PPP.AddProfile(ctx, pppProfile)
}

// GetByID gets profile by ID
func (s *BandwidthProfileService) GetByID(ctx context.Context, id uuid.UUID) (*model.BandwidthProfile, error) {
	return s.profileRepo.GetByID(ctx, id)
}

// GetByCode gets profile by code
func (s *BandwidthProfileService) GetByCode(ctx context.Context, code string) (*model.BandwidthProfile, error) {
	return s.profileRepo.GetByCode(ctx, code)
}

// List lists profiles with pagination
func (s *BandwidthProfileService) List(ctx context.Context, limit, offset int) ([]model.BandwidthProfile, int64, error) {
	profiles, err := s.profileRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.profileRepo.Count(ctx)
	return profiles, count, err
}

// ListActive lists active profiles
func (s *BandwidthProfileService) ListActive(ctx context.Context) ([]model.BandwidthProfile, error) {
	return s.profileRepo.ListActive(ctx)
}

// ListByRouterID lists profiles by router ID
func (s *BandwidthProfileService) ListByRouterID(ctx context.Context, routerID uuid.UUID, limit, offset int) ([]model.BandwidthProfile, int64, error) {
	profiles, err := s.profileRepo.ListByRouterID(ctx, routerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.profileRepo.CountByRouterID(ctx, routerID)
	return profiles, count, err
}

// ListActiveByRouterID lists active profiles by router ID
func (s *BandwidthProfileService) ListActiveByRouterID(ctx context.Context, routerID uuid.UUID) ([]model.BandwidthProfile, error) {
	return s.profileRepo.ListActiveByRouterID(ctx, routerID)
}

// Update updates a profile and syncs PPP profile to MikroTik
func (s *BandwidthProfileService) Update(ctx context.Context, profile *model.BandwidthProfile) error {
	// Get existing profile
	existingProfile, err := s.profileRepo.GetByID(ctx, uuid.MustParse(profile.ID))
	if err != nil {
		return err
	}

	// Connect to MikroTik
	routerID, err := uuid.Parse(profile.RouterID)
	if err != nil {
		return err
	}

	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return err
	}

	// Sync PPP profile to MikroTik
	if err := s.syncPPPProfile(ctx, mt, profile); err != nil {
		return fmt.Errorf("failed to sync to mikrotik: %w", err)
	}

	// Update database only if MikroTik succeeded
	if err := s.profileRepo.Update(ctx, profile); err != nil {
		// Rollback: restore old profile in MikroTik
		_ = s.syncPPPProfile(ctx, mt, existingProfile)
		return err
	}

	return nil
}

// Delete soft-deletes a profile and removes PPP profile from MikroTik
func (s *BandwidthProfileService) Delete(ctx context.Context, id uuid.UUID) error {
	// Get profile before delete
	profile, err := s.profileRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Remove from MikroTik first
	routerID, err := uuid.Parse(profile.RouterID)
	if err != nil {
		return err
	}

	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}

	// Remove PPP profile from MikroTik
	existing, _ := mt.PPP.GetProfileByName(ctx, profile.Name)
	if existing != nil {
		if err := mt.PPP.RemoveProfile(ctx, existing.ID); err != nil {
			return fmt.Errorf("failed to remove from mikrotik: %w", err)
		}
	}

	// Soft delete from database only if MikroTik succeeded
	return s.profileRepo.Delete(ctx, id)
}
