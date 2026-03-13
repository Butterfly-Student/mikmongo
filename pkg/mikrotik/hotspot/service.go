package hotspot

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Service provides Hotspot operations
type Service struct {
	client *client.Client
	repo   *Repository
}

// NewService creates a new Hotspot service
func NewService(c *client.Client) *Service {
	return &Service{
		client: c,
		repo:   NewRepository(c),
	}
}

// ─── Users ────────────────────────────────────────────────────────────────────

func (s *Service) GetUsers(ctx context.Context, profile string) ([]*domain.HotspotUser, error) {
	return s.repo.GetHotspotUsers(ctx, profile)
}

func (s *Service) GetUsersByComment(ctx context.Context, comment string) ([]*domain.HotspotUser, error) {
	return s.repo.GetHotspotUsersByComment(ctx, comment)
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*domain.HotspotUser, error) {
	return s.repo.GetHotspotUserByID(ctx, id)
}

func (s *Service) GetUserByName(ctx context.Context, name string) (*domain.HotspotUser, error) {
	return s.repo.GetHotspotUserByName(ctx, name)
}

func (s *Service) AddUser(ctx context.Context, user *domain.HotspotUser) (string, error) {
	return s.repo.AddHotspotUser(ctx, user)
}

func (s *Service) UpdateUser(ctx context.Context, id string, user *domain.HotspotUser) error {
	return s.repo.UpdateHotspotUser(ctx, id, user)
}

func (s *Service) RemoveUser(ctx context.Context, id string) error {
	return s.repo.RemoveHotspotUser(ctx, id)
}

func (s *Service) RemoveUsersByComment(ctx context.Context, comment string) error {
	return s.repo.RemoveHotspotUsersByComment(ctx, comment)
}

func (s *Service) RemoveUsers(ctx context.Context, ids []string) error {
	return s.repo.RemoveHotspotUsers(ctx, ids)
}

func (s *Service) DisableUser(ctx context.Context, id string) error {
	return s.repo.DisableHotspotUser(ctx, id)
}

func (s *Service) EnableUser(ctx context.Context, id string) error {
	return s.repo.EnableHotspotUser(ctx, id)
}

func (s *Service) DisableUsers(ctx context.Context, ids []string) error {
	return s.repo.DisableHotspotUsers(ctx, ids)
}

func (s *Service) EnableUsers(ctx context.Context, ids []string) error {
	return s.repo.EnableHotspotUsers(ctx, ids)
}

func (s *Service) ResetUserCounters(ctx context.Context, id string) error {
	return s.repo.ResetHotspotUserCounters(ctx, id)
}

func (s *Service) ResetUserCountersMultiple(ctx context.Context, ids []string) error {
	return s.repo.ResetHotspotUserCountersMultiple(ctx, ids)
}

func (s *Service) GetUsersCount(ctx context.Context) (int, error) {
	return s.repo.GetHotspotUsersCount(ctx)
}

// ─── Profiles ─────────────────────────────────────────────────────────────────

func (s *Service) GetProfiles(ctx context.Context) ([]*domain.UserProfile, error) {
	return s.repo.GetUserProfiles(ctx)
}

func (s *Service) GetProfileByID(ctx context.Context, id string) (*domain.UserProfile, error) {
	return s.repo.GetUserProfileByID(ctx, id)
}

func (s *Service) GetProfileByName(ctx context.Context, name string) (*domain.UserProfile, error) {
	return s.repo.GetUserProfileByName(ctx, name)
}

func (s *Service) AddProfile(ctx context.Context, profile *domain.UserProfile) (string, error) {
	return s.repo.AddUserProfile(ctx, profile)
}

func (s *Service) UpdateProfile(ctx context.Context, id string, profile *domain.UserProfile) error {
	return s.repo.UpdateUserProfile(ctx, id, profile)
}

func (s *Service) RemoveProfile(ctx context.Context, id string) error {
	return s.repo.RemoveUserProfile(ctx, id)
}

func (s *Service) DisableProfile(ctx context.Context, id string) error {
	return s.repo.DisableUserProfile(ctx, id)
}

func (s *Service) EnableProfile(ctx context.Context, id string) error {
	return s.repo.EnableUserProfile(ctx, id)
}

func (s *Service) RemoveProfiles(ctx context.Context, ids []string) error {
	return s.repo.RemoveUserProfiles(ctx, ids)
}

func (s *Service) DisableProfiles(ctx context.Context, ids []string) error {
	return s.repo.DisableUserProfiles(ctx, ids)
}

func (s *Service) EnableProfiles(ctx context.Context, ids []string) error {
	return s.repo.EnableUserProfiles(ctx, ids)
}

// ─── Active ───────────────────────────────────────────────────────────────────

func (s *Service) GetActive(ctx context.Context) ([]*domain.HotspotActive, error) {
	return s.repo.GetHotspotActive(ctx)
}

func (s *Service) GetActiveCount(ctx context.Context) (int, error) {
	return s.repo.GetHotspotActiveCount(ctx)
}

func (s *Service) RemoveActive(ctx context.Context, id string) error {
	return s.repo.RemoveHotspotActive(ctx, id)
}

func (s *Service) RemoveActives(ctx context.Context, ids []string) error {
	return s.repo.RemoveHotspotActives(ctx, ids)
}

// ─── Hosts ────────────────────────────────────────────────────────────────────

func (s *Service) GetHosts(ctx context.Context) ([]*domain.HotspotHost, error) {
	return s.repo.GetHotspotHosts(ctx)
}

func (s *Service) RemoveHost(ctx context.Context, id string) error {
	return s.repo.RemoveHotspotHost(ctx, id)
}

// ─── Servers ──────────────────────────────────────────────────────────────────

func (s *Service) GetServers(ctx context.Context) ([]string, error) {
	return s.repo.GetHotspotServers(ctx)
}

// ─── IP Binding ───────────────────────────────────────────────────────────────

func (s *Service) GetIPBindings(ctx context.Context) ([]*domain.HotspotIPBinding, error) {
	return s.repo.GetIPBindings(ctx)
}

func (s *Service) AddIPBinding(ctx context.Context, b *domain.HotspotIPBinding) (string, error) {
	return s.repo.AddIPBinding(ctx, b)
}

func (s *Service) RemoveIPBinding(ctx context.Context, id string) error {
	return s.repo.RemoveIPBinding(ctx, id)
}

func (s *Service) EnableIPBinding(ctx context.Context, id string) error {
	return s.repo.EnableIPBinding(ctx, id)
}

func (s *Service) DisableIPBinding(ctx context.Context, id string) error {
	return s.repo.DisableIPBinding(ctx, id)
}

// ─── Streaming ────────────────────────────────────────────────────────────────

func (s *Service) ListenActive(ctx context.Context, resultChan chan<- []*domain.HotspotActive) (func() error, error) {
	return s.repo.ListenHotspotActive(ctx, resultChan)
}

func (s *Service) ListenInactive(ctx context.Context, resultChan chan<- []*domain.HotspotUser) (func() error, error) {
	return s.repo.ListenHotspotInactive(ctx, resultChan)
}
