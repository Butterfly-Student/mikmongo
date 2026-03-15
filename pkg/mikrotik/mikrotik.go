package mikrotik

import (
	"github.com/Butterfly-Student/go-ros/client"
	ppprepo "github.com/Butterfly-Student/go-ros/repository/ppp"
)

// PPPClient exposes PPP secret and profile operations for ISP billing.
// It combines SecretRepository and ProfileRepository into a single access point.
type PPPClient struct {
	ppprepo.SecretRepository
	ppprepo.ProfileRepository
}

// Client is the MikroTik RouterOS facade used by internal services.
type Client struct {
	PPP *PPPClient
}

// NewClient creates a connected Client from a Config.
func NewClient(cfg client.Config) (*Client, error) {
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	return NewClientFromConnection(c), nil
}

// NewClientFromConnection creates a Client facade from a managed connection.
func NewClientFromConnection(c *client.Client) *Client {
	return &Client{
		PPP: &PPPClient{
			SecretRepository:  ppprepo.NewSecretRepository(c),
			ProfileRepository: ppprepo.NewProfileRepository(c),
		},
	}
}
