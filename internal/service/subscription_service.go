package service

import (
	"context"
	"fmt"
	"time"

	"mikmongo/internal/domain/subscription"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/pkg/mikrotik"
	mkdomain "mikmongo/pkg/mikrotik/domain"

	"github.com/google/uuid"
)

// SubscriptionService handles subscription business logic
type SubscriptionService struct {
	subRepo     repository.SubscriptionRepository
	profileRepo repository.BandwidthProfileRepository
	settingRepo repository.SystemSettingRepository
	subDomain   *subscription.Domain
	routerSvc   *RouterService
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(
	subRepo repository.SubscriptionRepository,
	profileRepo repository.BandwidthProfileRepository,
	settingRepo repository.SystemSettingRepository,
	subDomain *subscription.Domain,
	routerSvc *RouterService,
) *SubscriptionService {
	return &SubscriptionService{
		subRepo:     subRepo,
		profileRepo: profileRepo,
		settingRepo: settingRepo,
		subDomain:   subDomain,
		routerSvc:   routerSvc,
	}
}

// Create creates a new subscription and creates PPP secret in MikroTik
func (s *SubscriptionService) Create(ctx context.Context, sub *model.Subscription) error {
	// Validate profile belongs to the same router
	planID, err := uuid.Parse(sub.PlanID)
	if err != nil {
		return fmt.Errorf("invalid plan ID: %w", err)
	}
	profile, err := s.profileRepo.GetByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}
	if profile.RouterID != sub.RouterID {
		return fmt.Errorf("profile does not belong to the specified router")
	}

	if sub.Password == "" {
		pass, err := s.subDomain.GeneratePassword(12)
		if err != nil {
			return fmt.Errorf("failed to generate password: %w", err)
		}
		sub.Password = pass
	}
	if err := s.subDomain.ValidateCredentials(sub.Username, sub.Password); err != nil {
		return err
	}
	sub.Status = "pending"

	// Connect to MikroTik first
	routerID, _ := uuid.Parse(sub.RouterID)
	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}

	// Create PPP secret in MikroTik
	if err := s.createInMikroTik(ctx, mt, sub, profile); err != nil {
		return fmt.Errorf("failed to create in mikrotik: %w", err)
	}

	// Save to database only if MikroTik succeeded
	if err := s.subRepo.Create(ctx, sub); err != nil {
		// Rollback: remove from MikroTik
		_ = s.removeFromMikroTik(ctx, mt, sub)
		return fmt.Errorf("failed to save subscription: %w", err)
	}

	return nil
}

// GetByID gets subscription by ID
func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	return s.subRepo.GetByID(ctx, id)
}

// GetByCustomerID gets subscriptions by customer ID
func (s *SubscriptionService) GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]model.Subscription, error) {
	return s.subRepo.GetByCustomerID(ctx, customerID)
}

// GetByUsername gets subscription by username
func (s *SubscriptionService) GetByUsername(ctx context.Context, username string) (*model.Subscription, error) {
	return s.subRepo.GetByUsername(ctx, username)
}

// List lists subscriptions with pagination
func (s *SubscriptionService) List(ctx context.Context, limit, offset int) ([]model.Subscription, int64, error) {
	subs, err := s.subRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.subRepo.Count(ctx)
	return subs, count, err
}

// ListByRouterID lists subscriptions by router ID
func (s *SubscriptionService) ListByRouterID(ctx context.Context, routerID uuid.UUID, limit, offset int) ([]model.Subscription, int64, error) {
	subs, err := s.subRepo.ListByRouterID(ctx, routerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.subRepo.CountByRouterID(ctx, routerID)
	return subs, count, err
}

// Update updates a subscription and updates PPP secret in MikroTik
func (s *SubscriptionService) Update(ctx context.Context, sub *model.Subscription) error {
	// Get existing subscription
	existingSub, err := s.subRepo.GetByID(ctx, uuid.MustParse(sub.ID))
	if err != nil {
		return err
	}

	// Validate profile belongs to the same router if PlanID changed
	if sub.PlanID != "" && sub.PlanID != existingSub.PlanID {
		planID, err := uuid.Parse(sub.PlanID)
		if err != nil {
			return fmt.Errorf("invalid plan ID: %w", err)
		}
		profile, err := s.profileRepo.GetByID(ctx, planID)
		if err != nil {
			return fmt.Errorf("plan not found: %w", err)
		}
		if profile.RouterID != sub.RouterID {
			return fmt.Errorf("profile does not belong to the specified router")
		}
	}

	// Connect to MikroTik
	routerID, _ := uuid.Parse(sub.RouterID)
	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}

	// Get profile for MikroTik sync
	planID, _ := uuid.Parse(sub.PlanID)
	profile, err := s.profileRepo.GetByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	// Update PPP secret in MikroTik
	if err := s.updateInMikroTik(ctx, mt, sub, profile); err != nil {
		return fmt.Errorf("failed to update in mikrotik: %w", err)
	}

	// Update database only if MikroTik succeeded
	if err := s.subRepo.Update(ctx, sub); err != nil {
		// Rollback: restore old data in MikroTik
		_ = s.updateInMikroTik(ctx, mt, existingSub, profile)
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// Delete deletes a subscription and removes PPP secret from MikroTik
func (s *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	// Get subscription before delete
	sub, err := s.subRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Connect to MikroTik
	routerID, _ := uuid.Parse(sub.RouterID)
	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}

	// Remove from MikroTik first
	if err := s.removeFromMikroTik(ctx, mt, sub); err != nil {
		return fmt.Errorf("failed to remove from mikrotik: %w", err)
	}

	// Delete from database only if MikroTik succeeded
	return s.subRepo.Delete(ctx, id)
}

// Activate activates a subscription: enable PPP secret + set status=active
func (s *SubscriptionService) Activate(ctx context.Context, id uuid.UUID) error {
	sub, err := s.subRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.subDomain.ValidateStatusTransition(sub.Status, "active"); err != nil {
		return err
	}

	// Connect to MikroTik and enable PPP secret
	routerID, _ := uuid.Parse(sub.RouterID)
	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}

	if err := s.enableInMikroTik(ctx, mt, sub); err != nil {
		return fmt.Errorf("failed to enable in mikrotik: %w", err)
	}

	now := time.Now()
	sub.Status = "active"
	sub.ActivatedAt = &now
	return s.subRepo.Update(ctx, sub)
}

// createInMikroTik creates PPP secret in MikroTik
func (s *SubscriptionService) createInMikroTik(ctx context.Context, mt *mikrotik.Client, sub *model.Subscription, profile *model.BandwidthProfile) error {
	secret := &mkdomain.PPPSecret{
		Name:     sub.Username,
		Password: sub.Password,
		Profile:  profile.Name,
		Comment:  fmt.Sprintf("sub:%s", sub.ID),
	}
	if sub.StaticIP != nil {
		secret.RemoteAddress = *sub.StaticIP
	}

	return mt.PPP.AddSecret(ctx, secret)
}

// updateInMikroTik updates PPP secret in MikroTik
func (s *SubscriptionService) updateInMikroTik(ctx context.Context, mt *mikrotik.Client, sub *model.Subscription, profile *model.BandwidthProfile) error {
	secret := &mkdomain.PPPSecret{
		Name:     sub.Username,
		Password: sub.Password,
		Profile:  profile.Name,
		Comment:  fmt.Sprintf("sub:%s", sub.ID),
	}
	if sub.StaticIP != nil {
		secret.RemoteAddress = *sub.StaticIP
	}

	existing, err := mt.PPP.GetSecretByName(ctx, sub.Username)
	if err == nil && existing != nil {
		return mt.PPP.UpdateSecret(ctx, existing.ID, secret)
	}
	return mt.PPP.AddSecret(ctx, secret)
}

// getIsolateProfile returns the isolate profile name from system_settings or the profile's override
func (s *SubscriptionService) getIsolateProfile(ctx context.Context, profileOverride *string) string {
	if profileOverride != nil && *profileOverride != "" {
		return *profileOverride
	}
	setting, err := s.settingRepo.GetByGroupAndKey(ctx, "isolate", "pppoe_profile")
	if err != nil || setting.Value == nil {
		return "isolate"
	}
	return *setting.Value
}

// Isolate changes subscription profile to the isolate profile on MikroTik
func (s *SubscriptionService) Isolate(ctx context.Context, id uuid.UUID, reason string) error {
	sub, err := s.subRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if sub.Status == "isolated" {
		return nil
	}

	var isolateProfileOverride *string
	planID, err := uuid.Parse(sub.PlanID)
	if err == nil {
		if profile, err := s.profileRepo.GetByID(ctx, planID); err == nil {
			isolateProfileOverride = profile.IsolateProfileName
		}
	}

	routerID, _ := uuid.Parse(sub.RouterID)
	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}

	isolateProfile := s.getIsolateProfile(ctx, isolateProfileOverride)
	if err := s.applyProfile(ctx, mt, sub, isolateProfile); err != nil {
		return fmt.Errorf("failed to apply isolate profile: %w", err)
	}

	r := reason
	sub.SuspendReason = &r
	sub.Status = "isolated"
	return s.subRepo.Update(ctx, sub)
}

// applyProfile sets a new profile name on the PPP secret
func (s *SubscriptionService) applyProfile(ctx context.Context, mt *mikrotik.Client, sub *model.Subscription, profileName string) error {
	existing, err := mt.PPP.GetSecretByName(ctx, sub.Username)
	if err != nil {
		return err
	}
	return mt.PPP.UpdateSecret(ctx, existing.ID, &mkdomain.PPPSecret{Profile: profileName})
}

// Restore reverts subscription from isolated back to its original bandwidth profile
func (s *SubscriptionService) Restore(ctx context.Context, id uuid.UUID) error {
	sub, err := s.subRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	planID, err := uuid.Parse(sub.PlanID)
	if err != nil {
		return fmt.Errorf("invalid plan ID: %w", err)
	}
	profile, err := s.profileRepo.GetByID(ctx, planID)
	if err != nil {
		return err
	}

	routerID, _ := uuid.Parse(sub.RouterID)
	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}

	if err := s.applyProfile(ctx, mt, sub, profile.Name); err != nil {
		return fmt.Errorf("failed to restore profile: %w", err)
	}

	sub.Status = "active"
	sub.SuspendReason = nil
	return s.subRepo.Update(ctx, sub)
}

// Suspend disables the PPP secret on MikroTik and marks the subscription as suspended
func (s *SubscriptionService) Suspend(ctx context.Context, id uuid.UUID, reason string) error {
	sub, err := s.subRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	routerID, _ := uuid.Parse(sub.RouterID)
	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}

	if err := s.disableInMikroTik(ctx, mt, sub); err != nil {
		return fmt.Errorf("failed to disable in mikrotik: %w", err)
	}

	r := reason
	sub.SuspendReason = &r
	sub.Status = "suspended"
	return s.subRepo.Update(ctx, sub)
}

// disableInMikroTik disables the PPP secret
func (s *SubscriptionService) disableInMikroTik(ctx context.Context, mt *mikrotik.Client, sub *model.Subscription) error {
	existing, err := mt.PPP.GetSecretByName(ctx, sub.Username)
	if err != nil {
		return err
	}
	return mt.PPP.DisableSecret(ctx, existing.ID)
}

// enableInMikroTik enables the PPP secret
func (s *SubscriptionService) enableInMikroTik(ctx context.Context, mt *mikrotik.Client, sub *model.Subscription) error {
	existing, err := mt.PPP.GetSecretByName(ctx, sub.Username)
	if err != nil {
		return err
	}
	return mt.PPP.EnableSecret(ctx, existing.ID)
}

// Terminate removes the subscription from MikroTik and marks it as terminated
func (s *SubscriptionService) Terminate(ctx context.Context, id uuid.UUID) error {
	sub, err := s.subRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	routerID, _ := uuid.Parse(sub.RouterID)
	mt, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}

	if err := s.removeFromMikroTik(ctx, mt, sub); err != nil {
		return fmt.Errorf("failed to remove from mikrotik: %w", err)
	}

	now := time.Now()
	sub.Status = "terminated"
	sub.TerminatedAt = &now
	return s.subRepo.Update(ctx, sub)
}

// removeFromMikroTik removes the PPP secret from MikroTik
func (s *SubscriptionService) removeFromMikroTik(ctx context.Context, mt *mikrotik.Client, sub *model.Subscription) error {
	existing, err := mt.PPP.GetSecretByName(ctx, sub.Username)
	if err != nil {
		return nil // already gone
	}
	return mt.PPP.RemoveSecret(ctx, existing.ID)
}
