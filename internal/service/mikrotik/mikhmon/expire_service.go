package mikhmon

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	gorosmikhmon "github.com/Butterfly-Student/go-ros/repository/mikhmon"
	gorossystem "github.com/Butterfly-Student/go-ros/repository/system"

	"mikmongo/internal/service"
)

type MikhmonExpireService struct {
	routerSvc *service.RouterService
}

func NewMikhmonExpireService(routerSvc *service.RouterService) *MikhmonExpireService {
	return &MikhmonExpireService{
		routerSvc: routerSvc,
	}
}

func (s *MikhmonExpireService) SetupExpireMonitor(ctx context.Context, routerID uuid.UUID) error {
	repo, err := s.getExpireRepo(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to get expire repository: %w", err)
	}
	return repo.SetupExpireMonitor(ctx)
}

func (s *MikhmonExpireService) DisableExpireMonitor(ctx context.Context, routerID uuid.UUID) error {
	repo, err := s.getExpireRepo(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to get expire repository: %w", err)
	}
	return repo.DisableExpireMonitor(ctx)
}

func (s *MikhmonExpireService) IsExpireMonitorEnabled(ctx context.Context, routerID uuid.UUID) (bool, error) {
	repo, err := s.getExpireRepo(ctx, routerID)
	if err != nil {
		return false, fmt.Errorf("failed to get expire repository: %w", err)
	}
	return repo.IsExpireMonitorEnabled(ctx)
}

func (s *MikhmonExpireService) GenerateExpireMonitorScript() string {
	repo := gorosmikhmon.NewExpireRepository(nil, nil)
	return repo.GenerateExpireMonitorScript()
}

func (s *MikhmonExpireService) getExpireRepo(ctx context.Context, routerID uuid.UUID) (gorosmikhmon.ExpireRepository, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	conn := c.Conn()
	systemRepo := gorossystem.NewRepository(conn)
	return gorosmikhmon.NewExpireRepository(conn, systemRepo), nil
}
