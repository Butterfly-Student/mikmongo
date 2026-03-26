package mikrotik

import (
	"context"

	"github.com/Butterfly-Student/go-ros/client"
	"github.com/Butterfly-Student/go-ros/domain"
	hotspotrepo "github.com/Butterfly-Student/go-ros/repository/hotspot"
	ppprepo "github.com/Butterfly-Student/go-ros/repository/ppp"
)

// PPPClient exposes PPP operations for ISP billing.
type PPPClient struct {
	ppprepo.SecretRepository
	ppprepo.ProfileRepository
	active ppprepo.ActiveRepository
}

// GetActiveUsers returns active PPP sessions filtered by service type.
func (c *PPPClient) GetActiveUsers(ctx context.Context, service string) ([]*domain.PPPActive, error) {
	return c.active.GetActive(ctx, service)
}

// HotspotClient exposes Hotspot operations.
type HotspotClient struct {
	hotspotrepo.UserRepository
	hotspotrepo.ProfileRepository
	active  hotspotrepo.ActiveRepository
	host    hotspotrepo.HostRepository
	server  hotspotrepo.ServerRepository
}

// GetActive returns active hotspot sessions.
func (c *HotspotClient) GetActive(ctx context.Context) ([]*domain.HotspotActive, error) {
	return c.active.GetActive(ctx)
}

// GetActiveCount returns the count of active hotspot sessions.
func (c *HotspotClient) GetActiveCount(ctx context.Context) (int, error) {
	return c.active.GetActiveCount(ctx)
}

// GetHosts returns hotspot hosts.
func (c *HotspotClient) GetHosts(ctx context.Context) ([]*domain.HotspotHost, error) {
	return c.host.GetHosts(ctx)
}

// GetServers returns hotspot server names.
func (c *HotspotClient) GetServers(ctx context.Context) ([]string, error) {
	return c.server.GetServers(ctx)
}

// Client is the MikroTik RouterOS facade used by internal services.
type Client struct {
	PPP     *PPPClient
	Hotspot *HotspotClient
	conn    *client.Client
}

// Conn returns the underlying raw RouterOS client connection.
func (c *Client) Conn() *client.Client {
	return c.conn
}

// Close closes the underlying connection.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// NewClient creates a connected Client from a Config.
func NewClient(cfg client.Config) (*Client, error) {
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	return NewClientFromConnection(c), nil
}

// RunRaw executes any RouterOS command and returns raw map results.
func (c *Client) RunRaw(ctx context.Context, args []string) ([]map[string]string, error) {
	return c.conn.RunRaw(ctx, args)
}

// ListenRaw starts a streaming RouterOS command and sends raw data to resultChan.
func (c *Client) ListenRaw(ctx context.Context, args []string, resultChan chan<- map[string]string) (func() error, error) {
	return c.conn.ListenRaw(ctx, args, resultChan)
}

// NewClientFromConnection creates a Client facade from a managed connection.
func NewClientFromConnection(c *client.Client) *Client {
	return &Client{
		conn: c,
		PPP: &PPPClient{
			SecretRepository:  ppprepo.NewSecretRepository(c),
			ProfileRepository: ppprepo.NewProfileRepository(c),
			active:            ppprepo.NewActiveRepository(c),
		},
		Hotspot: &HotspotClient{
			UserRepository:    hotspotrepo.NewUserRepository(c),
			ProfileRepository: hotspotrepo.NewProfileRepository(c),
			active:            hotspotrepo.NewActiveRepository(c),
			host:              hotspotrepo.NewHostRepository(c),
			server:            hotspotrepo.NewServerRepository(c),
		},
	}
}
