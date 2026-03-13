package ppp

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Service provides PPP operations
type Service struct {
	client *client.Client
	repo   *Repository
}

// NewService creates a new PPP service
func NewService(c *client.Client) *Service {
	return &Service{
		client: c,
		repo:   NewRepository(c),
	}
}

// ─── Secrets ──────────────────────────────────────────────────────────────────

func (s *Service) GetSecrets(ctx context.Context, profile string) ([]*domain.PPPSecret, error) {
	return s.repo.GetPPPSecrets(ctx, profile)
}

func (s *Service) GetSecretByID(ctx context.Context, id string) (*domain.PPPSecret, error) {
	return s.repo.GetPPPSecretByID(ctx, id)
}

func (s *Service) GetSecretByName(ctx context.Context, name string) (*domain.PPPSecret, error) {
	return s.repo.GetPPPSecretByName(ctx, name)
}

func (s *Service) AddSecret(ctx context.Context, secret *domain.PPPSecret) error {
	return s.repo.AddPPPSecret(ctx, secret)
}

func (s *Service) UpdateSecret(ctx context.Context, id string, secret *domain.PPPSecret) error {
	return s.repo.UpdatePPPSecret(ctx, id, secret)
}

func (s *Service) RemoveSecret(ctx context.Context, id string) error {
	return s.repo.RemovePPPSecret(ctx, id)
}

func (s *Service) DisableSecret(ctx context.Context, id string) error {
	return s.repo.DisablePPPSecret(ctx, id)
}

func (s *Service) EnableSecret(ctx context.Context, id string) error {
	return s.repo.EnablePPPSecret(ctx, id)
}

func (s *Service) RemoveSecrets(ctx context.Context, ids []string) error {
	return s.repo.RemovePPPSecrets(ctx, ids)
}

func (s *Service) DisableSecrets(ctx context.Context, ids []string) error {
	return s.repo.DisablePPPSecrets(ctx, ids)
}

func (s *Service) EnableSecrets(ctx context.Context, ids []string) error {
	return s.repo.EnablePPPSecrets(ctx, ids)
}

// ─── Profiles ─────────────────────────────────────────────────────────────────

func (s *Service) GetProfiles(ctx context.Context) ([]*domain.PPPProfile, error) {
	return s.repo.GetPPPProfiles(ctx)
}

func (s *Service) GetProfileByID(ctx context.Context, id string) (*domain.PPPProfile, error) {
	return s.repo.GetPPPProfileByID(ctx, id)
}

func (s *Service) GetProfileByName(ctx context.Context, name string) (*domain.PPPProfile, error) {
	return s.repo.GetPPPProfileByName(ctx, name)
}

func (s *Service) AddProfile(ctx context.Context, profile *domain.PPPProfile) error {
	return s.repo.AddPPPProfile(ctx, profile)
}

func (s *Service) UpdateProfile(ctx context.Context, id string, profile *domain.PPPProfile) error {
	return s.repo.UpdatePPPProfile(ctx, id, profile)
}

func (s *Service) RemoveProfile(ctx context.Context, id string) error {
	return s.repo.RemovePPPProfile(ctx, id)
}

func (s *Service) DisableProfile(ctx context.Context, id string) error {
	return s.repo.DisablePPPProfile(ctx, id)
}

func (s *Service) EnableProfile(ctx context.Context, id string) error {
	return s.repo.EnablePPPProfile(ctx, id)
}

func (s *Service) RemoveProfiles(ctx context.Context, ids []string) error {
	return s.repo.RemovePPPProfiles(ctx, ids)
}

func (s *Service) DisableProfiles(ctx context.Context, ids []string) error {
	return s.repo.DisablePPPProfiles(ctx, ids)
}

func (s *Service) EnableProfiles(ctx context.Context, ids []string) error {
	return s.repo.EnablePPPProfiles(ctx, ids)
}

// ─── Active ───────────────────────────────────────────────────────────────────

func (s *Service) GetActiveUsers(ctx context.Context, service string) ([]*domain.PPPActive, error) {
	return s.repo.GetPPPActive(ctx, service)
}

func (s *Service) GetActiveByID(ctx context.Context, id string) (*domain.PPPActive, error) {
	return s.repo.GetPPPActiveByID(ctx, id)
}

func (s *Service) DisconnectActive(ctx context.Context, id string) error {
	return s.repo.RemovePPPActive(ctx, id)
}

func (s *Service) DisconnectActives(ctx context.Context, ids []string) error {
	return s.repo.RemovePPPActives(ctx, ids)
}

// ─── Streaming ────────────────────────────────────────────────────────────────

func (s *Service) ListenActive(ctx context.Context, resultChan chan<- []*domain.PPPActive) (func() error, error) {
	return s.repo.ListenPPPActive(ctx, resultChan)
}

func (s *Service) ListenInactive(ctx context.Context, resultChan chan<- []*domain.PPPSecret) (func() error, error) {
	return s.repo.ListenPPPInactive(ctx, resultChan)
}
