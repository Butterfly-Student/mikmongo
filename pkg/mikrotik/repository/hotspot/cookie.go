package hotspot

import "github.com/Butterfly-Student/go-ros/client"

// cookieRepository implements CookieRepository interface
type cookieRepository struct {
	client *client.Client
}

// NewCookieRepository creates a new cookie repository
func NewCookieRepository(c *client.Client) CookieRepository {
	return &cookieRepository{client: c}
}
