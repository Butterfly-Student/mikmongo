package mikrotik

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	mikrotikpkg "mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
)

// HotspotService provides MikroTik Hotspot operations
type HotspotService struct {
	routerService RouterConnector
}

// NewHotspotService creates a new Hotspot service
func NewHotspotService(routerService RouterConnector) *HotspotService {
	return &HotspotService{
		routerService: routerService,
	}
}

// getClient creates a MikroTik client for the specified router
func (s *HotspotService) getClient(ctx context.Context, routerID uuid.UUID) (*mikrotikpkg.Client, error) {
	return s.routerService.Connect(ctx, routerID)
}

// ─── Users ────────────────────────────────────────────────────────────────────

// GetUsers retrieves all hotspot users, optionally filtered by profile
func (s *HotspotService) GetUsers(ctx context.Context, routerID uuid.UUID, profile string) ([]*domain.HotspotUser, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetUsers(ctx, profile)
}

// GetUsersByComment retrieves hotspot users filtered by comment
func (s *HotspotService) GetUsersByComment(ctx context.Context, routerID uuid.UUID, comment string) ([]*domain.HotspotUser, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetUsersByComment(ctx, comment)
}

// GetUserByID retrieves a hotspot user by ID
func (s *HotspotService) GetUserByID(ctx context.Context, routerID uuid.UUID, id string) (*domain.HotspotUser, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetUserByID(ctx, id)
}

// GetUserByName retrieves a hotspot user by name
func (s *HotspotService) GetUserByName(ctx context.Context, routerID uuid.UUID, name string) (*domain.HotspotUser, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetUserByName(ctx, name)
}

// AddUser creates a new hotspot user
func (s *HotspotService) AddUser(ctx context.Context, routerID uuid.UUID, user *domain.HotspotUser) (string, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return "", fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.AddUser(ctx, user)
}

// UpdateUser updates an existing hotspot user
func (s *HotspotService) UpdateUser(ctx context.Context, routerID uuid.UUID, id string, user *domain.HotspotUser) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.UpdateUser(ctx, id, user)
}

// RemoveUser removes a hotspot user by ID
func (s *HotspotService) RemoveUser(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.RemoveUser(ctx, id)
}

// RemoveUsersByComment removes all hotspot users with the specified comment
func (s *HotspotService) RemoveUsersByComment(ctx context.Context, routerID uuid.UUID, comment string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.RemoveUsersByComment(ctx, comment)
}

// RemoveUsers removes multiple hotspot users by IDs
func (s *HotspotService) RemoveUsers(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.RemoveUsers(ctx, ids)
}

// DisableUser disables a hotspot user
func (s *HotspotService) DisableUser(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.DisableUser(ctx, id)
}

// EnableUser enables a hotspot user
func (s *HotspotService) EnableUser(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.EnableUser(ctx, id)
}

// DisableUsers disables multiple hotspot users
func (s *HotspotService) DisableUsers(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.DisableUsers(ctx, ids)
}

// EnableUsers enables multiple hotspot users
func (s *HotspotService) EnableUsers(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.EnableUsers(ctx, ids)
}

// ResetUserCounters resets counters for a hotspot user
func (s *HotspotService) ResetUserCounters(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.ResetUserCounters(ctx, id)
}

// ResetUserCountersMultiple resets counters for multiple hotspot users
func (s *HotspotService) ResetUserCountersMultiple(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.ResetUserCountersMultiple(ctx, ids)
}

// GetUsersCount returns the total count of hotspot users
func (s *HotspotService) GetUsersCount(ctx context.Context, routerID uuid.UUID) (int, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetUsersCount(ctx)
}

// ─── Profiles ─────────────────────────────────────────────────────────────────

// GetProfiles retrieves all hotspot user profiles
func (s *HotspotService) GetProfiles(ctx context.Context, routerID uuid.UUID) ([]*domain.UserProfile, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetProfiles(ctx)
}

// GetProfileByID retrieves a hotspot user profile by ID
func (s *HotspotService) GetProfileByID(ctx context.Context, routerID uuid.UUID, id string) (*domain.UserProfile, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetProfileByID(ctx, id)
}

// GetProfileByName retrieves a hotspot user profile by name
func (s *HotspotService) GetProfileByName(ctx context.Context, routerID uuid.UUID, name string) (*domain.UserProfile, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetProfileByName(ctx, name)
}

// AddProfile creates a new hotspot user profile
func (s *HotspotService) AddProfile(ctx context.Context, routerID uuid.UUID, profile *domain.UserProfile) (string, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return "", fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.AddProfile(ctx, profile)
}

// UpdateProfile updates an existing hotspot user profile
func (s *HotspotService) UpdateProfile(ctx context.Context, routerID uuid.UUID, id string, profile *domain.UserProfile) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.UpdateProfile(ctx, id, profile)
}

// RemoveProfile removes a hotspot user profile
func (s *HotspotService) RemoveProfile(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.RemoveProfile(ctx, id)
}

// DisableProfile disables a hotspot user profile
func (s *HotspotService) DisableProfile(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.DisableProfile(ctx, id)
}

// EnableProfile enables a hotspot user profile
func (s *HotspotService) EnableProfile(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.EnableProfile(ctx, id)
}

// RemoveProfiles removes multiple hotspot user profiles
func (s *HotspotService) RemoveProfiles(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.RemoveProfiles(ctx, ids)
}

// DisableProfiles disables multiple hotspot user profiles
func (s *HotspotService) DisableProfiles(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.DisableProfiles(ctx, ids)
}

// EnableProfiles enables multiple hotspot user profiles
func (s *HotspotService) EnableProfiles(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.EnableProfiles(ctx, ids)
}

// ─── Active Sessions ──────────────────────────────────────────────────────────

// GetActive retrieves all active hotspot sessions
func (s *HotspotService) GetActive(ctx context.Context, routerID uuid.UUID) ([]*domain.HotspotActive, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetActive(ctx)
}

// GetActiveCount returns the count of active hotspot sessions
func (s *HotspotService) GetActiveCount(ctx context.Context, routerID uuid.UUID) (int, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return 0, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetActiveCount(ctx)
}

// RemoveActive removes an active hotspot session
func (s *HotspotService) RemoveActive(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.RemoveActive(ctx, id)
}

// RemoveActives removes multiple active hotspot sessions
func (s *HotspotService) RemoveActives(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.RemoveActives(ctx, ids)
}

// ─── Hosts ────────────────────────────────────────────────────────────────────

// GetHosts retrieves all hotspot hosts
func (s *HotspotService) GetHosts(ctx context.Context, routerID uuid.UUID) ([]*domain.HotspotHost, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetHosts(ctx)
}

// RemoveHost removes a hotspot host
func (s *HotspotService) RemoveHost(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.RemoveHost(ctx, id)
}

// ─── Servers ──────────────────────────────────────────────────────────────────

// GetServers retrieves all hotspot server names
func (s *HotspotService) GetServers(ctx context.Context, routerID uuid.UUID) ([]string, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetServers(ctx)
}

// ─── IP Bindings ──────────────────────────────────────────────────────────────

// GetIPBindings retrieves all hotspot IP bindings
func (s *HotspotService) GetIPBindings(ctx context.Context, routerID uuid.UUID) ([]*domain.HotspotIPBinding, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.GetIPBindings(ctx)
}

// AddIPBinding creates a new hotspot IP binding
func (s *HotspotService) AddIPBinding(ctx context.Context, routerID uuid.UUID, binding *domain.HotspotIPBinding) (string, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return "", fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.AddIPBinding(ctx, binding)
}

// RemoveIPBinding removes a hotspot IP binding
func (s *HotspotService) RemoveIPBinding(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.RemoveIPBinding(ctx, id)
}

// EnableIPBinding enables a hotspot IP binding
func (s *HotspotService) EnableIPBinding(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.EnableIPBinding(ctx, id)
}

// DisableIPBinding disables a hotspot IP binding
func (s *HotspotService) DisableIPBinding(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.Hotspot.DisableIPBinding(ctx, id)
}

// ─── Streaming (for WebSocket) ────────────────────────────────────────────────

// ListenActive streams active hotspot sessions
func (s *HotspotService) ListenActive(ctx context.Context, routerID uuid.UUID, resultChan chan<- []*domain.HotspotActive) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	// Note: Don't close client here for streaming operations
	// The cleanup function returned will handle it

	cleanup, err := client.Hotspot.ListenActive(ctx, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	// Wrap cleanup to also close the client
	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}

// ListenInactive streams inactive hotspot users
func (s *HotspotService) ListenInactive(ctx context.Context, routerID uuid.UUID, resultChan chan<- []*domain.HotspotUser) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	cleanup, err := client.Hotspot.ListenInactive(ctx, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}
