// Package parse provides per-topic classifiers that split raw RouterOS
// key=value strings into InfluxDB tags (string labels) and fields (float64
// measurements).
package parse

import (
	"strconv"
	"strings"
	"time"
)

// Parser classifies raw RouterOS sentence fields into TSDB tags and fields.
type Parser interface {
	// Parse splits raw key=value strings from RouterOS into:
	//   tags   – string labels used as InfluxDB tag keys
	//   fields – numeric measurements used as InfluxDB field keys
	//
	// Keys not matched by either set are silently dropped (non-numeric strings
	// that are not useful as dimensions).
	Parse(raw map[string]string) (tags map[string]string, fields map[string]float64)
}

// ForTopic returns the appropriate Parser for the given topic name.
// If no specific parser exists, a generic fallback is returned that attempts
// to parse every value as float64.
func ForTopic(topic string) Parser {
	switch topic {
	case "system-resource":
		return systemResourceParser{}
	case "interfaces":
		return interfacesParser{}
	case "interface-traffic":
		return interfaceTrafficParser{}
	case "queues":
		return queuesParser{}
	case "queue-stats":
		return queueStatsParser{}
	case "ip-addresses":
		return ipAddressesParser{}
	case "ppp-active":
		return pppActiveParser{}
	case "dhcp-leases":
		return dhcpLeasesParser{}
	default:
		return genericParser{}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// system-resource
// ─────────────────────────────────────────────────────────────────────────────

type systemResourceParser struct{}

func (systemResourceParser) Parse(raw map[string]string) (map[string]string, map[string]float64) {
	tags := map[string]string{}
	fields := map[string]float64{}

	numericKeys := map[string]bool{
		"cpu-load":       true,
		"free-memory":    true,
		"total-memory":   true,
		"free-hdd-space": true,
		"total-hdd-space": true,
		"bad-blocks":     true,
		"write-sect-total": true,
		"write-sect-since-reboot": true,
	}

	for k, v := range raw {
		if k == "version" || k == "platform" || k == "board-name" || k == "architecture-name" {
			tags[k] = v
			continue
		}
		if k == "uptime" {
			if secs := parseUptime(v); secs >= 0 {
				fields["uptime_seconds"] = secs
			}
			continue
		}
		if numericKeys[k] {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				fields[k] = f
			}
			continue
		}
		// drop the rest (cpu-frequency, etc.) if not numeric
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			fields[k] = f
		}
	}
	return tags, fields
}

// ─────────────────────────────────────────────────────────────────────────────
// interface-traffic  (/interface/monitor-traffic interval=1)
// ─────────────────────────────────────────────────────────────────────────────

// interfaceTrafficParser handles /interface/monitor-traffic output.
// RouterOS returns rx/tx-bits-per-second as plain integers.
type interfaceTrafficParser struct{}

func (interfaceTrafficParser) Parse(raw map[string]string) (map[string]string, map[string]float64) {
	tags := extractTags(raw, "name")
	fields := extractNumericFields(raw,
		"rx-packets-per-second",
		"rx-bits-per-second",
		"fp-rx-packets-per-second",
		"fp-rx-bits-per-second",
		"rx-drops-per-second",
		"rx-errors-per-second",
		"tx-packets-per-second",
		"tx-bits-per-second",
		"fp-tx-packets-per-second",
		"fp-tx-bits-per-second",
		"tx-drops-per-second",
		"tx-queue-drops-per-second",
		"tx-errors-per-second",
	)
	return tags, fields
}

// ─────────────────────────────────────────────────────────────────────────────
// interfaces  (/interface/print follow=)
// ─────────────────────────────────────────────────────────────────────────────

type interfacesParser struct{}

func (interfacesParser) Parse(raw map[string]string) (map[string]string, map[string]float64) {
	tags := extractTags(raw, "name", "type", "running", "disabled")
	fields := extractNumericFields(raw,
		"rx-byte", "tx-byte",
		"rx-packet", "tx-packet",
		"rx-error", "tx-error",
		"rx-drop", "tx-drop",
		"link-downs",
		"actual-mtu",
	)
	return tags, fields
}

// ─────────────────────────────────────────────────────────────────────────────
// queue-stats  (/queue/simple/print stats interval=1)
// ─────────────────────────────────────────────────────────────────────────────

// queueStatsParser handles /queue/simple/print stats output.
// RouterOS rate fields are formatted strings like "15.1kbps/163.3kbps".
type queueStatsParser struct{}

func (queueStatsParser) Parse(raw map[string]string) (map[string]string, map[string]float64) {
	tags := extractTags(raw, "name", "target")
	fields := map[string]float64{}

	// rate = "upload/download" current throughput (formatted, e.g. "15.1kbps/163.3kbps")
	if v, ok := raw["rate"]; ok {
		ul, dl := splitSlashRate(v)
		if ul >= 0 {
			fields["rate-ul-bps"] = ul
		}
		if dl >= 0 {
			fields["rate-dl-bps"] = dl
		}
	}
	// packet-rate = "tx/rx"
	if v, ok := raw["packet-rate"]; ok {
		tx, rx := splitSlash(v)
		if tx >= 0 {
			fields["packet-rate-tx"] = tx
		}
		if rx >= 0 {
			fields["packet-rate-rx"] = rx
		}
	}
	// cumulative bytes tx/rx
	if v, ok := raw["bytes"]; ok {
		tx, rx := splitSlash(v)
		if tx >= 0 {
			fields["bytes-tx"] = tx
		}
		if rx >= 0 {
			fields["bytes-rx"] = rx
		}
	}
	return tags, fields
}

// ─────────────────────────────────────────────────────────────────────────────
// queues  (/queue/simple/print follow= — cumulative counters)
// ─────────────────────────────────────────────────────────────────────────────

type queuesParser struct{}

func (queuesParser) Parse(raw map[string]string) (map[string]string, map[string]float64) {
	tags := extractTags(raw, "name", "target", "disabled")
	fields := map[string]float64{}

	// "bytes" in queues is "tx/rx" combined string like "12345/67890"
	if v, ok := raw["bytes"]; ok {
		tx, rx := splitSlash(v)
		if tx >= 0 {
			fields["tx-byte"] = tx
		}
		if rx >= 0 {
			fields["rx-byte"] = rx
		}
	}
	if v, ok := raw["packets"]; ok {
		tx, rx := splitSlash(v)
		if tx >= 0 {
			fields["tx-packet"] = tx
		}
		if rx >= 0 {
			fields["rx-packet"] = rx
		}
	}
	if v, ok := raw["max-limit"]; ok {
		ul, dl := splitSlash(v)
		if ul >= 0 {
			fields["max-limit-ul"] = ul
		}
		if dl >= 0 {
			fields["max-limit-dl"] = dl
		}
	}
	return tags, fields
}

// ─────────────────────────────────────────────────────────────────────────────
// ip-addresses
// ─────────────────────────────────────────────────────────────────────────────

type ipAddressesParser struct{}

func (ipAddressesParser) Parse(raw map[string]string) (map[string]string, map[string]float64) {
	tags := extractTags(raw, ".id", "address", "interface", "network", "disabled", "dynamic", "invalid")
	// No meaningful numeric fields; emit a presence counter.
	fields := map[string]float64{"present": 1}
	return tags, fields
}

// ─────────────────────────────────────────────────────────────────────────────
// ppp-active
// ─────────────────────────────────────────────────────────────────────────────

type pppActiveParser struct{}

func (pppActiveParser) Parse(raw map[string]string) (map[string]string, map[string]float64) {
	tags := extractTags(raw, ".id", "name", "service", "address", "caller-id")
	fields := map[string]float64{}
	if v, ok := raw["uptime"]; ok {
		if secs := parseUptime(v); secs >= 0 {
			fields["uptime_seconds"] = secs
		}
	}
	return tags, fields
}

// ─────────────────────────────────────────────────────────────────────────────
// dhcp-leases
// ─────────────────────────────────────────────────────────────────────────────

type dhcpLeasesParser struct{}

func (dhcpLeasesParser) Parse(raw map[string]string) (map[string]string, map[string]float64) {
	tags := extractTags(raw, ".id", "address", "mac-address", "host-name", "status", "server")
	fields := map[string]float64{"present": 1}
	return tags, fields
}

// ─────────────────────────────────────────────────────────────────────────────
// Generic fallback
// ─────────────────────────────────────────────────────────────────────────────

// genericParser tries to parse every value as float64; strings become tags.
type genericParser struct{}

func (genericParser) Parse(raw map[string]string) (map[string]string, map[string]float64) {
	tags := map[string]string{}
	fields := map[string]float64{}
	for k, v := range raw {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			fields[k] = f
		} else {
			tags[k] = v
		}
	}
	return tags, fields
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// extractTags picks the listed keys from raw into a tags map.
func extractTags(raw map[string]string, keys ...string) map[string]string {
	tags := make(map[string]string, len(keys))
	for _, k := range keys {
		if v, ok := raw[k]; ok {
			tags[k] = v
		}
	}
	return tags
}

// extractNumericFields picks the listed keys from raw and parses them as float64.
func extractNumericFields(raw map[string]string, keys ...string) map[string]float64 {
	fields := make(map[string]float64, len(keys))
	for _, k := range keys {
		if v, ok := raw[k]; ok {
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				fields[k] = f
			}
		}
	}
	return fields
}

// splitSlash parses a "tx/rx" string like "12345/67890" into two float64 values.
// Returns (-1, -1) if parsing fails.
func splitSlash(s string) (float64, float64) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return -1, -1
	}
	a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err1 != nil || err2 != nil {
		return -1, -1
	}
	return a, b
}

// splitSlashRate parses a RouterOS formatted rate string like "15.1kbps/163.3kbps"
// or "1.2Mbps/8.0Mbps" into two float64 values in bits per second.
// Returns (-1, -1) if parsing fails.
func splitSlashRate(s string) (float64, float64) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 {
		return -1, -1
	}
	a := parseRateBps(strings.TrimSpace(parts[0]))
	b := parseRateBps(strings.TrimSpace(parts[1]))
	return a, b
}

// parseRateBps converts a RouterOS rate string (e.g. "15.1kbps", "1.2Mbps", "500bps")
// to float64 bits per second. Returns -1 on failure.
func parseRateBps(s string) float64 {
	s = strings.ToLower(s)
	multipliers := []struct {
		suffix string
		mult   float64
	}{
		{"gbps", 1e9},
		{"mbps", 1e6},
		{"kbps", 1e3},
		{"bps", 1},
	}
	for _, m := range multipliers {
		if strings.HasSuffix(s, m.suffix) {
			numStr := strings.TrimSuffix(s, m.suffix)
			if f, err := strconv.ParseFloat(strings.TrimSpace(numStr), 64); err == nil {
				return f * m.mult
			}
		}
	}
	// Fallback: try plain number
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return -1
}

// parseUptime converts a RouterOS uptime string (e.g. "2d3h15m40s", "1w2d")
// into total seconds. Returns -1 on parse failure.
func parseUptime(s string) float64 {
	// RouterOS format: Xw Xd Xh Xm Xs (may omit leading zero segments)
	d, err := time.ParseDuration(toGoDuration(s))
	if err != nil {
		return -1
	}
	return d.Seconds()
}

// toGoDuration converts RouterOS uptime (weeks, days) into a Go parseable duration.
// Go's time.ParseDuration does not support 'w' or 'd'.
func toGoDuration(s string) string {
	// Replace weeks and days into hours.
	result := s
	result = replaceUnit(result, "w", 7*24)
	result = replaceUnit(result, "d", 24)
	return result
}

// replaceUnit replaces a unit suffix (e.g. "3w") with the equivalent hours (e.g. "504h").
func replaceUnit(s, unit string, hoursPerUnit int) string {
	var out strings.Builder
	remaining := s
	for {
		idx := strings.Index(remaining, unit)
		if idx < 0 {
			out.WriteString(remaining)
			break
		}
		// Find the start of the number before this unit.
		numStart := idx
		for numStart > 0 && remaining[numStart-1] >= '0' && remaining[numStart-1] <= '9' {
			numStart--
		}
		numStr := remaining[numStart:idx]
		out.WriteString(remaining[:numStart])
		if n, err := strconv.Atoi(numStr); err == nil {
			out.WriteString(strconv.Itoa(n*hoursPerUnit))
			out.WriteString("h")
		}
		remaining = remaining[idx+len(unit):]
	}
	return out.String()
}
