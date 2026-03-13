package client

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-routeros/routeros/v3/proto"
)

// parseInt parses an integer from a RouterOS string value.
func parseInt(s string) int64 {
	if s == "" {
		return 0
	}
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

// parseFloat parses a float from a RouterOS string value.
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// parseBool parses a bool from a RouterOS "yes"/"true" string.
func parseBool(s string) bool {
	return s == "true" || s == "yes"
}

// formatInt formats an int64 to string.
func formatInt(n int64) string {
	return strconv.FormatInt(n, 10)
}

// ParseRate parses a rate string with unit (bps, kbps, Mbps, Gbps) to bits per second.
// Examples: "0bps" -> 0, "74.3kbps" -> 74300, "2.2Mbps" -> 2200000.
func ParseRate(s string) int64 {
	if s == "" || s == "0" {
		return 0
	}

	var value float64
	var unit string

	for i := len(s) - 1; i >= 0; i-- {
		if (s[i] >= '0' && s[i] <= '9') || s[i] == '.' {
			value = parseFloat(s[:i+1])
			unit = s[i+1:]
			break
		}
	}

	if unit == "" {
		return parseInt(s)
	}

	switch strings.ToLower(unit) {
	case "bps":
		return int64(value)
	case "kbps":
		return int64(value * 1_000)
	case "mbps":
		return int64(value * 1_000_000)
	case "gbps":
		return int64(value * 1_000_000_000)
	default:
		return parseInt(s)
	}
}

// SplitSlashValue splits a "in/out" value into (in, out).
// Example: "11519824317/96401664078" -> (11519824317, 96401664078)
func SplitSlashValue(value string) (int64, int64) {
	if value == "" {
		return 0, 0
	}
	parts := strings.Split(value, "/")
	if len(parts) == 2 {
		return parseInt(parts[0]), parseInt(parts[1])
	}
	return parseInt(value), 0
}

// SplitRateValue splits a rate with unit (bps, kbps, Mbps, Gbps).
// Example: "74.3kbps/2.2Mbps" -> (74300, 2200000)
func SplitRateValue(value string) (int64, int64) {
	if value == "" {
		return 0, 0
	}
	parts := strings.Split(value, "/")
	if len(parts) == 2 {
		return ParseRate(parts[0]), ParseRate(parts[1])
	}
	return ParseRate(value), 0
}

// ParseByteSize parses byte size strings like "6.4MiB", "6756KiB", "10MiB", "2GiB" to int64.
func ParseByteSize(s string) int64 {
	if s == "" {
		return 0
	}

	multipliers := map[string]int64{
		"KiB": 1024,
		"MiB": 1024 * 1024,
		"GiB": 1024 * 1024 * 1024,
		"TiB": 1024 * 1024 * 1024 * 1024,
		"KB":  1024,
		"MB":  1024 * 1024,
		"GB":  1024 * 1024 * 1024,
		"TB":  1024 * 1024 * 1024 * 1024,
		"B":   1,
		"MHz": 1,
		"GHz": 1000,
	}

	for suffix, multiplier := range multipliers {
		if len(s) > len(suffix) && s[len(s)-len(suffix):] == suffix {
			numStr := s[:len(s)-len(suffix)]
			if val, err := strconv.ParseFloat(numStr, 64); err == nil {
				return int64(val * float64(multiplier))
			}
			return 0
		}
	}

	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

// BatchDebounce is the silence window used by ListenBatches to detect the gap
// between RouterOS interval ticks. RouterOS sends a burst of !re sentences then
// goes silent until the next tick. 200ms safely fits inside any practical interval.
const BatchDebounce = 200 * time.Millisecond

// ListenBatches reads sentences from a RouterOS follow stream and emits complete
// batches. A batch is considered complete when no new sentence arrives within
// debounce duration — this detects the gap between RouterOS interval ticks.
func ListenBatches(
	ctx context.Context,
	sentences <-chan *proto.Sentence,
	debounce time.Duration,
) <-chan []*proto.Sentence {
	out := make(chan []*proto.Sentence, 4)

	go func() {
		defer close(out)

		for {
			var first *proto.Sentence
			select {
			case s, ok := <-sentences:
				if !ok {
					return
				}
				first = s
			case <-ctx.Done():
				return
			}

			batch := []*proto.Sentence{first}
			timer := time.NewTimer(debounce)

		collect:
			for {
				select {
				case s, ok := <-sentences:
					if !ok {
						timer.Stop()
						break collect
					}
					batch = append(batch, s)
					if !timer.Stop() {
						select {
						case <-timer.C:
						default:
						}
					}
					timer.Reset(debounce)
				case <-timer.C:
					break collect
				case <-ctx.Done():
					timer.Stop()
					return
				}
			}

			select {
			case out <- batch:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}
