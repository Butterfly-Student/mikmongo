package parse_test

import (
	"testing"

	"mikmongo/internal/collector/parse"
)

// ─────────────────────────────────────────────────────────────────────────────
// ForTopic routing
// ─────────────────────────────────────────────────────────────────────────────

func TestForTopic_KnownTopics(t *testing.T) {
	topics := []string{
		"system-resource",
		"interfaces",
		"queues",
		"ip-addresses",
		"ppp-active",
		"dhcp-leases",
	}
	for _, topic := range topics {
		p := parse.ForTopic(topic)
		if p == nil {
			t.Errorf("ForTopic(%q) returned nil", topic)
		}
	}
}

func TestForTopic_Unknown_FallsBackToGeneric(t *testing.T) {
	p := parse.ForTopic("unknown-topic")
	if p == nil {
		t.Fatal("ForTopic(unknown) should return a generic parser, not nil")
	}
	raw := map[string]string{"some-value": "42", "label": "hello"}
	tags, fields := p.Parse(raw)
	if fields["some-value"] != 42 {
		t.Errorf("generic parser: expected fields[some-value]=42, got %v", fields["some-value"])
	}
	if tags["label"] != "hello" {
		t.Errorf("generic parser: expected tags[label]=hello, got %q", tags["label"])
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// system-resource parser
// ─────────────────────────────────────────────────────────────────────────────

func TestSystemResource_NumericFields(t *testing.T) {
	raw := map[string]string{
		"cpu-load":     "42",
		"free-memory":  "134217728",
		"total-memory": "268435456",
		"uptime":       "1d2h3m4s",
		"version":      "7.10.2",
	}
	p := parse.ForTopic("system-resource")
	tags, fields := p.Parse(raw)

	assertField(t, fields, "cpu-load", 42)
	assertField(t, fields, "free-memory", 134217728)
	assertField(t, fields, "total-memory", 268435456)
	if fields["uptime_seconds"] <= 0 {
		t.Errorf("expected uptime_seconds > 0, got %v", fields["uptime_seconds"])
	}
	if tags["version"] != "7.10.2" {
		t.Errorf("expected version as tag, got %q", tags["version"])
	}
}

func TestSystemResource_MissingFields_NoError(t *testing.T) {
	p := parse.ForTopic("system-resource")
	tags, fields := p.Parse(map[string]string{})
	_ = tags
	_ = fields // should not panic
}

// ─────────────────────────────────────────────────────────────────────────────
// interfaces parser
// ─────────────────────────────────────────────────────────────────────────────

func TestInterfaces_TagsAndFields(t *testing.T) {
	raw := map[string]string{
		"name":       "ether1",
		"type":       "ether",
		"running":    "true",
		"rx-byte":    "9999",
		"tx-byte":    "8888",
		"rx-packet":  "100",
		"tx-packet":  "200",
	}
	p := parse.ForTopic("interfaces")
	tags, fields := p.Parse(raw)

	if tags["name"] != "ether1" {
		t.Errorf("expected tag name=ether1, got %q", tags["name"])
	}
	assertField(t, fields, "rx-byte", 9999)
	assertField(t, fields, "tx-byte", 8888)
	assertField(t, fields, "rx-packet", 100)
	assertField(t, fields, "tx-packet", 200)
}

// ─────────────────────────────────────────────────────────────────────────────
// queues parser — slash-separated tx/rx
// ─────────────────────────────────────────────────────────────────────────────

func TestQueues_SlashBytes(t *testing.T) {
	raw := map[string]string{
		"name":      "user1",
		"target":    "192.168.1.10/32",
		"bytes":     "12345/67890",
		"packets":   "10/20",
		"max-limit": "5000000/10000000",
	}
	p := parse.ForTopic("queues")
	tags, fields := p.Parse(raw)

	if tags["name"] != "user1" {
		t.Errorf("expected tag name=user1")
	}
	assertField(t, fields, "tx-byte", 12345)
	assertField(t, fields, "rx-byte", 67890)
	assertField(t, fields, "tx-packet", 10)
	assertField(t, fields, "rx-packet", 20)
	assertField(t, fields, "max-limit-ul", 5000000)
	assertField(t, fields, "max-limit-dl", 10000000)
}

func TestQueues_MalformedBytes_Skipped(t *testing.T) {
	raw := map[string]string{"bytes": "notanumber/alsono"}
	p := parse.ForTopic("queues")
	_, fields := p.Parse(raw)
	if _, ok := fields["tx-byte"]; ok {
		t.Error("malformed bytes should not produce a tx-byte field")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// ip-addresses parser
// ─────────────────────────────────────────────────────────────────────────────

func TestIPAddresses_PresenceField(t *testing.T) {
	raw := map[string]string{
		".id":       "*1",
		"address":   "192.168.1.1/24",
		"interface": "ether1",
		"network":   "192.168.1.0",
		"disabled":  "false",
	}
	p := parse.ForTopic("ip-addresses")
	tags, fields := p.Parse(raw)

	if fields["present"] != 1 {
		t.Errorf("expected fields[present]=1, got %v", fields["present"])
	}
	if tags["interface"] != "ether1" {
		t.Errorf("expected tag interface=ether1")
	}
	if tags["address"] != "192.168.1.1/24" {
		t.Errorf("expected tag address=192.168.1.1/24")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// ppp-active parser — uptime string
// ─────────────────────────────────────────────────────────────────────────────

func TestPPPActive_UptimeSeconds(t *testing.T) {
	raw := map[string]string{
		"name":      "user1",
		"service":   "pppoe",
		"address":   "10.0.0.1",
		"caller-id": "AA:BB:CC:DD:EE:FF",
		"uptime":    "1h30m",
	}
	p := parse.ForTopic("ppp-active")
	tags, fields := p.Parse(raw)

	if tags["name"] != "user1" {
		t.Errorf("expected tag name=user1")
	}
	expected := float64(90 * 60) // 1h30m = 5400s
	if fields["uptime_seconds"] != expected {
		t.Errorf("expected uptime_seconds=%v, got %v", expected, fields["uptime_seconds"])
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// dhcp-leases parser
// ─────────────────────────────────────────────────────────────────────────────

func TestDHCPLeases_PresenceAndTags(t *testing.T) {
	raw := map[string]string{
		"address":     "192.168.1.100",
		"mac-address": "AA:BB:CC:DD:EE:FF",
		"host-name":   "mypc",
		"status":      "bound",
		"server":      "dhcp1",
	}
	p := parse.ForTopic("dhcp-leases")
	tags, fields := p.Parse(raw)

	if fields["present"] != 1 {
		t.Errorf("expected fields[present]=1")
	}
	if tags["status"] != "bound" {
		t.Errorf("expected tag status=bound")
	}
	if tags["mac-address"] != "AA:BB:CC:DD:EE:FF" {
		t.Errorf("expected tag mac-address")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// uptime parsing edge cases
// ─────────────────────────────────────────────────────────────────────────────

func TestPPPActive_WeeksDays_Uptime(t *testing.T) {
	cases := []struct {
		uptime  string
		wantSec float64
	}{
		{"1w", 7 * 24 * 3600},
		{"1d", 24 * 3600},
		{"2d3h", 2*24*3600 + 3*3600},
		{"10m30s", 10*60 + 30},
		{"0s", 0},
	}

	p := parse.ForTopic("ppp-active")
	for _, tc := range cases {
		raw := map[string]string{"uptime": tc.uptime}
		_, fields := p.Parse(raw)
		if fields["uptime_seconds"] != tc.wantSec {
			t.Errorf("uptime %q: expected %v seconds, got %v", tc.uptime, tc.wantSec, fields["uptime_seconds"])
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// helper
// ─────────────────────────────────────────────────────────────────────────────

func assertField(t *testing.T, fields map[string]float64, key string, want float64) {
	t.Helper()
	got, ok := fields[key]
	if !ok {
		t.Errorf("expected field %q to exist", key)
		return
	}
	if got != want {
		t.Errorf("field %q: got %v, want %v", key, got, want)
	}
}
