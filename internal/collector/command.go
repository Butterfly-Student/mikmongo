// Package collector provides a background engine that streams and polls
// MikroTik monitoring data into Redis (cache + pub/sub).
//
// It is designed to be self-contained: no existing handlers, services,
// or routes are modified. Integration is done later by wiring the Manager
// into cmd/server/main.go.
package collector

// DataCategory controls how a command's results are stored and forwarded.
// Both categories stream events via RouterOS =follow= — no polling occurs.
type DataCategory int

const (
	// SlowChanging data is streamed via RouterOS =follow= and cached with a
	// long TTL. Cache is refreshed only when RouterOS fires a change event.
	// Examples: IP addresses, PPP active sessions, DHCP leases.
	SlowChanging DataCategory = iota

	// RealTime data is streamed via RouterOS =follow= and cached with a
	// short TTL. Each update is also published to a Redis Pub/Sub channel.
	RealTime
)

// Command describes a single RouterOS monitoring command.
type Command struct {
	// Args is the full RouterOS command with parameters.
	Args []string
	// Category determines TTL and whether pub/sub is used.
	Category DataCategory
	// Topic is the logical name used in Redis keys and channels.
	Topic string
}

// DefaultCommands returns the standard set of monitoring commands.
func DefaultCommands() []Command {
	return []Command{
		// ── RealTime (streamed with =follow=) ────────────────────────────

		{
			Args:     []string{"/system/resource/print", "=follow="},
			Category: RealTime,
			Topic:    "system-resource",
		},
		{
			Args: []string{
				"/interface/print",
				"=.proplist=name,type,rx-byte,tx-byte,rx-packet,tx-packet,running",
				"=follow=",
			},
			Category: RealTime,
			Topic:    "interfaces",
		},
		{
			Args: []string{
				"/queue/simple/print",
				"=.proplist=name,max-limit,bytes,packets,target",
				"=follow=",
			},
			Category: RealTime,
			Topic:    "queues",
		},

		// ── SlowChanging (event-driven via =follow=, long TTL, no pub/sub) ─
		// RouterOS pushes a notification only when data actually changes,
		// so this adds zero periodic load to the router.

		{
			Args: []string{
				"/ip/address/print",
				"=.proplist=.id,address,interface,network,disabled",
				"=follow=",
			},
			Category: SlowChanging,
			Topic:    "ip-addresses",
		},
		{
			Args: []string{
				"/ppp/active/print",
				"=.proplist=.id,name,service,address,uptime,caller-id",
				"=follow=",
			},
			Category: SlowChanging,
			Topic:    "ppp-active",
		},
		{
			Args: []string{
				"/ip/dhcp-server/lease/print",
				"=.proplist=.id,address,mac-address,host-name,status,server",
				"=follow=",
			},
			Category: SlowChanging,
			Topic:    "dhcp-leases",
		},
	}
}

// FilterByCategory returns commands matching the given category.
func FilterByCategory(commands []Command, cat DataCategory) []Command {
	var out []Command
	for _, cmd := range commands {
		if cmd.Category == cat {
			out = append(out, cmd)
		}
	}
	return out
}
