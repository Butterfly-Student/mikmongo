package mikrotik

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	mikrotikpkg "mikmongo/pkg/mikrotik"
	"mikmongo/pkg/mikrotik/domain"
)

// PPPService provides MikroTik PPP operations
type PPPService struct {
	routerService RouterConnector
}

// NewPPPService creates a new PPP service
func NewPPPService(routerService RouterConnector) *PPPService {
	return &PPPService{
		routerService: routerService,
	}
}

// getClient creates a MikroTik client for the specified router
func (s *PPPService) getClient(ctx context.Context, routerID uuid.UUID) (*mikrotikpkg.Client, error) {
	return s.routerService.Connect(ctx, routerID)
}

// ─── Secrets ──────────────────────────────────────────────────────────────────

// GetSecrets retrieves all PPP secrets, optionally filtered by profile
func (s *PPPService) GetSecrets(ctx context.Context, routerID uuid.UUID, profile string) ([]*domain.PPPSecret, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.GetSecrets(ctx, profile)
}

// GetSecretByID retrieves a PPP secret by ID
func (s *PPPService) GetSecretByID(ctx context.Context, routerID uuid.UUID, id string) (*domain.PPPSecret, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.GetSecretByID(ctx, id)
}

// GetSecretByName retrieves a PPP secret by name
func (s *PPPService) GetSecretByName(ctx context.Context, routerID uuid.UUID, name string) (*domain.PPPSecret, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.GetSecretByName(ctx, name)
}

// AddSecret creates a new PPP secret
func (s *PPPService) AddSecret(ctx context.Context, routerID uuid.UUID, secret *domain.PPPSecret) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.AddSecret(ctx, secret)
}

// UpdateSecret updates an existing PPP secret
func (s *PPPService) UpdateSecret(ctx context.Context, routerID uuid.UUID, id string, secret *domain.PPPSecret) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.UpdateSecret(ctx, id, secret)
}

// RemoveSecret removes a PPP secret by ID
func (s *PPPService) RemoveSecret(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.RemoveSecret(ctx, id)
}

// DisableSecret disables a PPP secret
func (s *PPPService) DisableSecret(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.DisableSecret(ctx, id)
}

// EnableSecret enables a PPP secret
func (s *PPPService) EnableSecret(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.EnableSecret(ctx, id)
}

// RemoveSecrets removes multiple PPP secrets by IDs
func (s *PPPService) RemoveSecrets(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.RemoveSecrets(ctx, ids)
}

// DisableSecrets disables multiple PPP secrets
func (s *PPPService) DisableSecrets(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.DisableSecrets(ctx, ids)
}

// EnableSecrets enables multiple PPP secrets
func (s *PPPService) EnableSecrets(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.EnableSecrets(ctx, ids)
}

// ─── Profiles ─────────────────────────────────────────────────────────────────

// GetProfiles retrieves all PPP profiles
func (s *PPPService) GetProfiles(ctx context.Context, routerID uuid.UUID) ([]*domain.PPPProfile, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.GetProfiles(ctx)
}

// GetProfileByID retrieves a PPP profile by ID
func (s *PPPService) GetProfileByID(ctx context.Context, routerID uuid.UUID, id string) (*domain.PPPProfile, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.GetProfileByID(ctx, id)
}

// GetProfileByName retrieves a PPP profile by name
func (s *PPPService) GetProfileByName(ctx context.Context, routerID uuid.UUID, name string) (*domain.PPPProfile, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.GetProfileByName(ctx, name)
}

// AddProfile creates a new PPP profile
func (s *PPPService) AddProfile(ctx context.Context, routerID uuid.UUID, profile *domain.PPPProfile) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.AddProfile(ctx, profile)
}

// UpdateProfile updates an existing PPP profile
func (s *PPPService) UpdateProfile(ctx context.Context, routerID uuid.UUID, id string, profile *domain.PPPProfile) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.UpdateProfile(ctx, id, profile)
}

// RemoveProfile removes a PPP profile
func (s *PPPService) RemoveProfile(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.RemoveProfile(ctx, id)
}

// DisableProfile disables a PPP profile
func (s *PPPService) DisableProfile(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.DisableProfile(ctx, id)
}

// EnableProfile enables a PPP profile
func (s *PPPService) EnableProfile(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.EnableProfile(ctx, id)
}

// RemoveProfiles removes multiple PPP profiles
func (s *PPPService) RemoveProfiles(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.RemoveProfiles(ctx, ids)
}

// DisableProfiles disables multiple PPP profiles
func (s *PPPService) DisableProfiles(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.DisableProfiles(ctx, ids)
}

// EnableProfiles enables multiple PPP profiles
func (s *PPPService) EnableProfiles(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.EnableProfiles(ctx, ids)
}

// ─── Active Sessions ──────────────────────────────────────────────────────────

// GetActiveUsers retrieves all active PPP sessions
func (s *PPPService) GetActiveUsers(ctx context.Context, routerID uuid.UUID, service string) ([]*domain.PPPActive, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.GetActiveUsers(ctx, service)
}

// GetActiveByID retrieves an active PPP session by ID
func (s *PPPService) GetActiveByID(ctx context.Context, routerID uuid.UUID, id string) (*domain.PPPActive, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.GetActiveByID(ctx, id)
}

// DisconnectActive disconnects an active PPP session
func (s *PPPService) DisconnectActive(ctx context.Context, routerID uuid.UUID, id string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.DisconnectActive(ctx, id)
}

// DisconnectActives disconnects multiple active PPP sessions
func (s *PPPService) DisconnectActives(ctx context.Context, routerID uuid.UUID, ids []string) error {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to connect to router: %w", err)
	}
	defer client.Close()

	return client.PPP.DisconnectActives(ctx, ids)
}

// ─── Streaming (for WebSocket) ────────────────────────────────────────────────

// ListenActive streams active PPP sessions
func (s *PPPService) ListenActive(ctx context.Context, routerID uuid.UUID, resultChan chan<- []*domain.PPPActive) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	cleanup, err := client.PPP.ListenActive(ctx, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}

// ListenInactive streams inactive PPP secrets
func (s *PPPService) ListenInactive(ctx context.Context, routerID uuid.UUID, resultChan chan<- []*domain.PPPSecret) (func() error, error) {
	client, err := s.getClient(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	cleanup, err := client.PPP.ListenInactive(ctx, resultChan)
	if err != nil {
		client.Close()
		return nil, err
	}

	return func() error {
		defer client.Close()
		return cleanup()
	}, nil
}
