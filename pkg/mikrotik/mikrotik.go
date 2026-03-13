// Package mikrotik provides RouterOS API client
package mikrotik

import (
	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/firewall"
	"mikmongo/pkg/mikrotik/hotspot"
	"mikmongo/pkg/mikrotik/ipaddress"
	"mikmongo/pkg/mikrotik/ippool"
	"mikmongo/pkg/mikrotik/monitor"
	"mikmongo/pkg/mikrotik/ppp"
	"mikmongo/pkg/mikrotik/queue"
	"mikmongo/pkg/mikrotik/report"
	"mikmongo/pkg/mikrotik/script"
	"mikmongo/pkg/mikrotik/voucher"
)

// Client is the main Mikrotik client facade
type Client struct {
	conn      *client.Client
	PPP       *ppp.Service
	Hotspot   *hotspot.Service
	Queue     *queue.Service
	Firewall  *firewall.Service
	Monitor   *monitor.Service
	Report    *report.Service
	Voucher   *voucher.Generator
	Script    *script.Manager
	IPPool    *ippool.Service
	IPAddress *ipaddress.Service
}

// NewClient creates a new Mikrotik client facade using the provided Config.
func NewClient(cfg client.Config) (*Client, error) {
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:      c,
		PPP:       ppp.NewService(c),
		Hotspot:   hotspot.NewService(c),
		Queue:     queue.NewService(c),
		Firewall:  firewall.NewService(c),
		Monitor:   monitor.NewService(c),
		Report:    report.NewService(c),
		Voucher:   voucher.NewGenerator(),
		Script:    script.NewManager(c),
		IPPool:    ippool.NewService(c),
		IPAddress: ipaddress.NewService(c),
	}, nil
}

// Conn returns the underlying client connection.
func (c *Client) Conn() *client.Client {
	return c.conn
}

// Close closes the underlying connection.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// NewClientFromConnection creates a facade from an existing connection.
// This is useful when using a connection manager that handles connection lifecycle.
func NewClientFromConnection(conn *client.Client) *Client {
	return &Client{
		conn:      conn,
		PPP:       ppp.NewService(conn),
		Hotspot:   hotspot.NewService(conn),
		Queue:     queue.NewService(conn),
		Firewall:  firewall.NewService(conn),
		Monitor:   monitor.NewService(conn),
		Report:    report.NewService(conn),
		Voucher:   voucher.NewGenerator(),
		Script:    script.NewManager(conn),
		IPPool:    ippool.NewService(conn),
		IPAddress: ipaddress.NewService(conn),
	}
}
