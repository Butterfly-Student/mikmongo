package mikhmon

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	goroshotspot "github.com/Butterfly-Student/go-ros/repository/hotspot"
	gorosmikhmon "github.com/Butterfly-Student/go-ros/repository/mikhmon"

	"mikmongo/internal/service"
)

type MikhmonProfileService struct {
	routerSvc *service.RouterService
}

func NewMikhmonProfileService(routerSvc *service.RouterService) *MikhmonProfileService {
	return &MikhmonProfileService{
		routerSvc: routerSvc,
	}
}

func (s *MikhmonProfileService) CreateProfile(ctx context.Context, routerID uuid.UUID, req *mikhmonDomain.ProfileRequest) error {
	repo, err := s.getProfileRepo(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to get profile repository: %w", err)
	}
	return repo.CreateProfile(ctx, req)
}

func (s *MikhmonProfileService) UpdateProfile(ctx context.Context, routerID uuid.UUID, id string, req *mikhmonDomain.ProfileRequest) error {
	repo, err := s.getProfileRepo(ctx, routerID)
	if err != nil {
		return fmt.Errorf("failed to get profile repository: %w", err)
	}
	return repo.UpdateProfile(ctx, id, req)
}

func (s *MikhmonProfileService) GenerateOnLoginScript(data *mikhmonDomain.OnLoginScriptData) string {
	repo := gorosmikhmon.NewProfileRepository(nil)
	return repo.GenerateOnLoginScript(data)
}

func (s *MikhmonProfileService) getProfileRepo(ctx context.Context, routerID uuid.UUID) (gorosmikhmon.ProfileRepository, error) {
	c, err := s.routerSvc.GetMikrotikClient(ctx, routerID)
	if err != nil {
		return nil, err
	}
	conn := c.Conn()
	hotspotRepo := goroshotspot.NewRepository(conn)
	return gorosmikhmon.NewProfileRepository(hotspotRepo), nil
}
