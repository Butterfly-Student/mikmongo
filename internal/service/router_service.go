package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	mikrotik "github.com/Butterfly-Student/go-ros"
	"github.com/Butterfly-Student/go-ros/client"
	"mikmongo/pkg/redis"
	"mikmongo/utils"
)

// RouterService handles router business logic
type RouterService struct {
	routerRepo  repository.RouterDeviceRepository
	encKey      []byte // 32-byte AES key (from config or env)
	redisClient *redis.Client
	manager     *client.Manager
}

// NewRouterService creates a new router service
func NewRouterService(
	routerRepo repository.RouterDeviceRepository,
	encKey string,
	redisClient *redis.Client,
	logger *zap.Logger,
) *RouterService {
	key := []byte(encKey)
	// Pad or truncate to 32 bytes
	padded := make([]byte, 32)
	copy(padded, key)

	// Create manager for connection pooling
	manager := client.NewManager(logger)

	return &RouterService{
		routerRepo:  routerRepo,
		encKey:      padded,
		redisClient: redisClient,
		manager:     manager,
	}
}

// Create creates a new router device (encrypts password)
func (s *RouterService) Create(ctx context.Context, router *model.MikrotikRouter, plainPassword string) error {
	enc, err := utils.Encrypt(s.encKey, plainPassword)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}
	router.PasswordEncrypted = enc
	return s.routerRepo.Create(ctx, router)
}

// GetByID gets router by ID
func (s *RouterService) GetByID(ctx context.Context, id uuid.UUID) (*model.MikrotikRouter, error) {
	return s.routerRepo.GetByID(ctx, id)
}

// GetDevice gets router device by ID (alias for GetByID)
func (s *RouterService) GetDevice(ctx context.Context, id uuid.UUID) (*model.MikrotikRouter, error) {
	return s.routerRepo.GetByID(ctx, id)
}

// List lists all routers
func (s *RouterService) List(ctx context.Context, limit, offset int) ([]model.MikrotikRouter, error) {
	return s.routerRepo.List(ctx, limit, offset)
}

// Update updates router (re-encrypts password if provided)
func (s *RouterService) Update(ctx context.Context, router *model.MikrotikRouter, plainPassword string) error {
	if plainPassword != "" {
		enc, err := utils.Encrypt(s.encKey, plainPassword)
		if err != nil {
			return fmt.Errorf("failed to encrypt password: %w", err)
		}
		router.PasswordEncrypted = enc

		// If password changed, unregister from manager to force reconnect with new credentials
		s.manager.Unregister(router.ID)
	}
	return s.routerRepo.Update(ctx, router)
}

// Delete deletes a router
func (s *RouterService) Delete(ctx context.Context, id uuid.UUID) error {
	// Unregister from manager before deleting
	router, err := s.routerRepo.GetByID(ctx, id)
	if err == nil {
		s.manager.Unregister(router.ID)
	}
	return s.routerRepo.Delete(ctx, id)
}

// getRouterConfig returns the client config for a router
func (s *RouterService) getRouterConfig(ctx context.Context, routerID uuid.UUID) (client.Config, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return client.Config{}, fmt.Errorf("router not found: %w", err)
	}

	password, err := utils.Decrypt(s.encKey, router.PasswordEncrypted)
	if err != nil {
		return client.Config{}, fmt.Errorf("failed to decrypt password: %w", err)
	}

	return client.Config{
		Host:     router.Address,
		Port:     router.APIPort,
		Username: router.Username,
		Password: password,
		UseTLS:   router.UseSSL,
		Timeout:  10 * time.Second,
	}, nil
}

// GetMikrotikClient returns a connected MikroTik client for a router using Manager
func (s *RouterService) GetMikrotikClient(ctx context.Context, routerID uuid.UUID) (*mikrotik.Client, error) {
	cfg, err := s.getRouterConfig(ctx, routerID)
	if err != nil {
		return nil, err
	}

	// Use manager to get or create connection
	clientConn, err := s.manager.GetOrConnect(ctx, routerID.String(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to router: %w", err)
	}

	// Create facade client using the connection from manager
	return mikrotik.NewClientFromConnection(clientConn), nil
}

// GetRouterHost returns the IP/hostname of the router.
// Used by the collector to populate the "host" tag in TSDB data points.
// Does not decrypt credentials — cheaper than getRouterConfig.
func (s *RouterService) GetRouterHost(ctx context.Context, routerID uuid.UUID) (string, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return "", fmt.Errorf("router not found: %w", err)
	}
	return router.Address, nil
}

// Connect connects to a MikroTik router using its stored credentials (backward compatibility)
func (s *RouterService) Connect(ctx context.Context, routerID uuid.UUID) (*mikrotik.Client, error) {
	return s.GetMikrotikClient(ctx, routerID)
}

// GetMikrotikAdapter returns a MikrotikClientAdapter for a router.
// This implements MikrotikProvider, enabling SubscriptionService to be unit-tested
// by substituting a mock provider.
func (s *RouterService) GetMikrotikAdapter(ctx context.Context, routerID uuid.UUID) (MikrotikClientAdapter, error) {
	client, err := s.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	return &mikrotikClientWrapper{client: client}, nil
}

// TestConnection tests connection to a router
func (s *RouterService) TestConnection(ctx context.Context, routerID uuid.UUID) error {
	cfg, err := s.getRouterConfig(ctx, routerID)
	if err != nil {
		return err
	}
	return s.manager.TestConnection(ctx, cfg)
}

// SyncDevice syncs router status (health check + update LastSeenAt)
func (s *RouterService) SyncDevice(ctx context.Context, id uuid.UUID) error {
	router, err := s.routerRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Just test connection, don't keep the client
	cfg, err := s.getRouterConfig(ctx, id)
	if err != nil {
		// Mark offline
		router.Status = "offline"
		_ = s.routerRepo.Update(ctx, router)
		return err
	}

	if err := s.manager.TestConnection(ctx, cfg); err != nil {
		// Mark offline
		router.Status = "offline"
		_ = s.routerRepo.Update(ctx, router)
		return err
	}

	now := time.Now()
	router.Status = "online"
	router.LastSeenAt = &now
	return s.routerRepo.UpdateLastSync(ctx, id)
}

// SyncAllDevices syncs all active routers
func (s *RouterService) SyncAllDevices(ctx context.Context) error {
	devices, err := s.routerRepo.GetActive(ctx)
	if err != nil {
		return err
	}
	for _, device := range devices {
		deviceUUID, err := uuid.Parse(device.ID)
		if err != nil {
			continue
		}
		_ = s.SyncDevice(ctx, deviceUUID)
	}
	return nil
}

// SelectRouter sets the active router for a user
func (s *RouterService) SelectRouter(ctx context.Context, userID string, routerID uuid.UUID) (*model.MikrotikRouter, error) {
	router, err := s.routerRepo.GetByID(ctx, routerID)
	if err != nil {
		return nil, fmt.Errorf("router not found: %w", err)
	}

	if !router.IsActive {
		return nil, errors.New("router is not active")
	}

	if err := s.redisClient.SetSelectedRouter(ctx, userID, routerID.String(), 7*24*time.Hour); err != nil {
		return nil, fmt.Errorf("failed to save selected router: %w", err)
	}

	return router, nil
}

// GetSelectedRouter returns the currently selected router for a user
func (s *RouterService) GetSelectedRouter(ctx context.Context, userID string) (*model.MikrotikRouter, error) {
	routerIDStr, err := s.redisClient.GetSelectedRouter(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get selected router: %w", err)
	}

	if routerIDStr == "" {
		return nil, nil
	}

	routerID, err := uuid.Parse(routerIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid router id: %w", err)
	}

	return s.routerRepo.GetByID(ctx, routerID)
}

// CloseAllConnections closes all managed connections (call on shutdown)
func (s *RouterService) CloseAllConnections() {
	s.manager.CloseAll()
}
