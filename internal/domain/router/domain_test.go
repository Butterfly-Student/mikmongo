package router

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mikmongo/internal/model"
)

func TestValidateConnection(t *testing.T) {
	d := NewDomain()

	t.Run("valid host and port", func(t *testing.T) {
		assert.NoError(t, d.ValidateConnection("192.168.1.1", 8728))
	})

	t.Run("host empty → error", func(t *testing.T) {
		assert.Error(t, d.ValidateConnection("", 8728))
	})

	t.Run("port 0 → error", func(t *testing.T) {
		assert.Error(t, d.ValidateConnection("192.168.1.1", 0))
	})

	t.Run("port 65536 → error", func(t *testing.T) {
		assert.Error(t, d.ValidateConnection("192.168.1.1", 65536))
	})

	t.Run("port negative → error", func(t *testing.T) {
		assert.Error(t, d.ValidateConnection("192.168.1.1", -1))
	})

	t.Run("port 1 → valid", func(t *testing.T) {
		assert.NoError(t, d.ValidateConnection("router.local", 1))
	})

	t.Run("port 65535 → valid", func(t *testing.T) {
		assert.NoError(t, d.ValidateConnection("router.local", 65535))
	})
}

func TestIsOnline(t *testing.T) {
	d := NewDomain()

	t.Run("status online → true", func(t *testing.T) {
		r := &model.MikrotikRouter{Status: "online"}
		assert.True(t, d.IsOnline(r))
	})

	t.Run("status offline → false", func(t *testing.T) {
		r := &model.MikrotikRouter{Status: "offline"}
		assert.False(t, d.IsOnline(r))
	})

	t.Run("status unknown → false", func(t *testing.T) {
		r := &model.MikrotikRouter{Status: "unknown"}
		assert.False(t, d.IsOnline(r))
	})
}

func TestCanConnect(t *testing.T) {
	d := NewDomain()

	t.Run("active and status online → true", func(t *testing.T) {
		r := &model.MikrotikRouter{IsActive: true, Status: "online"}
		assert.True(t, d.CanConnect(r))
	})

	t.Run("active and status unknown → true", func(t *testing.T) {
		r := &model.MikrotikRouter{IsActive: true, Status: "unknown"}
		assert.True(t, d.CanConnect(r))
	})

	t.Run("active and status offline → false", func(t *testing.T) {
		r := &model.MikrotikRouter{IsActive: true, Status: "offline"}
		assert.False(t, d.CanConnect(r))
	})

	t.Run("not active and status online → false", func(t *testing.T) {
		r := &model.MikrotikRouter{IsActive: false, Status: "online"}
		assert.False(t, d.CanConnect(r))
	})

	t.Run("not active and status offline → false", func(t *testing.T) {
		r := &model.MikrotikRouter{IsActive: false, Status: "offline"}
		assert.False(t, d.CanConnect(r))
	})
}

func TestShouldSync(t *testing.T) {
	d := NewDomain()

	t.Run("last seen nil → should sync (never synced)", func(t *testing.T) {
		assert.True(t, d.ShouldSync(nil, 5))
	})

	t.Run("last seen long ago → should sync", func(t *testing.T) {
		old := time.Now().Add(-30 * time.Minute)
		assert.True(t, d.ShouldSync(&old, 5))
	})

	t.Run("last seen recently → should not sync", func(t *testing.T) {
		recent := time.Now().Add(-1 * time.Minute)
		assert.False(t, d.ShouldSync(&recent, 5))
	})

	t.Run("last seen exactly at interval → should sync", func(t *testing.T) {
		// Just over the interval threshold
		atInterval := time.Now().Add(-6 * time.Minute)
		assert.True(t, d.ShouldSync(&atInterval, 5))
	})
}

func TestIsStale(t *testing.T) {
	d := NewDomain()

	t.Run("last seen nil → stale", func(t *testing.T) {
		r := &model.MikrotikRouter{LastSeenAt: nil}
		assert.True(t, d.IsStale(r, 10))
	})

	t.Run("last seen > staleMinutes ago → stale", func(t *testing.T) {
		old := time.Now().Add(-20 * time.Minute)
		r := &model.MikrotikRouter{LastSeenAt: &old}
		assert.True(t, d.IsStale(r, 10))
	})

	t.Run("last seen < staleMinutes ago → not stale", func(t *testing.T) {
		recent := time.Now().Add(-2 * time.Minute)
		r := &model.MikrotikRouter{LastSeenAt: &recent}
		assert.False(t, d.IsStale(r, 10))
	})
}
