package mikrotik

import (
	"context"

	"github.com/google/uuid"
	"mikmongo/pkg/mikrotik"
)

// RouterConnector defines the interface for connecting to routers
type RouterConnector interface {
	Connect(ctx context.Context, routerID uuid.UUID) (*mikrotik.Client, error)
}
