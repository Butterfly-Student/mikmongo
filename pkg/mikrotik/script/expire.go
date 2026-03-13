package script

import (
	"context"

	"mikmongo/pkg/mikrotik/client"
)

const ExpireMonitorName = "Mikhmon-Expire-Monitor"

// EnsureExpireMonitor ensures scheduler "Mikhmon-Expire-Monitor" exists and is enabled.
// Returns status: "created", "enabled", or "existing".
func EnsureExpireMonitor(ctx context.Context, c *client.Client, script string) (string, error) {
	reply, err := c.RunContext(ctx, "/system/scheduler/print", "?name="+ExpireMonitorName)
	if err != nil {
		return "", err
	}

	if len(reply.Re) == 0 {
		_, err := c.RunContext(ctx,
			"/system/scheduler/add",
			"=name="+ExpireMonitorName,
			"=start-time=00:00:00",
			"=interval=00:01:00",
			"=on-event="+script,
			"=disabled=no",
			"=comment=Mikhmon Expire Monitor",
		)
		if err != nil {
			return "", err
		}
		return "created", nil
	}

	entry := reply.Re[0].Map
	if entry["disabled"] == "true" || entry["disabled"] == "yes" {
		_, err := c.RunContext(ctx,
			"/system/scheduler/set",
			"=.id="+entry[".id"],
			"=interval=00:01:00",
			"=on-event="+script,
			"=disabled=no",
		)
		if err != nil {
			return "", err
		}
		return "enabled", nil
	}

	return "existing", nil
}

// Manager wraps OnLoginGenerator and expire monitor for facade use.
type Manager struct {
	client    *client.Client
	generator *OnLoginGenerator
}

// NewManager creates a new script Manager.
func NewManager(c *client.Client) *Manager {
	return &Manager{
		client:    c,
		generator: NewOnLoginGenerator(),
	}
}

// Generator returns the OnLoginGenerator.
func (m *Manager) Generator() *OnLoginGenerator {
	return m.generator
}

// EnsureExpireMonitor delegates to the package-level function.
func (m *Manager) EnsureExpireMonitor(ctx context.Context, scriptStr string) (string, error) {
	return EnsureExpireMonitor(ctx, m.client, scriptStr)
}
