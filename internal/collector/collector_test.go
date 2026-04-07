package collector

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers / mocks
// ─────────────────────────────────────────────────────────────────────────────

// recordingSink captures all Write calls for assertion in tests.
type recordingSink struct {
	mu     sync.Mutex
	points []DataPoint
}

func (r *recordingSink) Write(_ context.Context, point DataPoint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.points = append(r.points, point)
	return nil
}

func (r *recordingSink) Close() error { return nil }

func (r *recordingSink) all() []DataPoint {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]DataPoint, len(r.points))
	copy(out, r.points)
	return out
}

// ─────────────────────────────────────────────────────────────────────────────
// command.go tests
// ─────────────────────────────────────────────────────────────────────────────

func TestDefaultCommands_Count(t *testing.T) {
	cmds := DefaultCommands()
	if len(cmds) == 0 {
		t.Fatal("DefaultCommands() returned no commands")
	}
}

func TestDefaultCommands_Categories(t *testing.T) {
	cmds := DefaultCommands()

	realtimeCmds := FilterByCategory(cmds, RealTime)
	slowCmds := FilterByCategory(cmds, SlowChanging)

	if len(realtimeCmds) == 0 {
		t.Error("expected at least one RealTime command")
	}
	if len(slowCmds) == 0 {
		t.Error("expected at least one SlowChanging command")
	}
}

func TestDefaultCommands_ExpectedTopics(t *testing.T) {
	expected := map[string]bool{
		"system-resource": false,
		"interfaces":      false,
		"queues":          false,
		"ip-addresses":    false,
		"ppp-active":      false,
		"dhcp-leases":     false,
	}

	for _, cmd := range DefaultCommands() {
		if _, ok := expected[cmd.Topic]; ok {
			expected[cmd.Topic] = true
		}
	}

	for topic, found := range expected {
		if !found {
			t.Errorf("expected topic %q not found in DefaultCommands()", topic)
		}
	}
}

func TestDefaultCommands_AllArgsNonEmpty(t *testing.T) {
	for _, cmd := range DefaultCommands() {
		if len(cmd.Args) == 0 {
			t.Errorf("command topic=%q has empty Args", cmd.Topic)
		}
	}
}

func TestFilterByCategory(t *testing.T) {
	cmds := []Command{
		{Topic: "a", Category: RealTime},
		{Topic: "b", Category: SlowChanging},
		{Topic: "c", Category: RealTime},
	}

	rt := FilterByCategory(cmds, RealTime)
	if len(rt) != 2 {
		t.Errorf("expected 2 RealTime commands, got %d", len(rt))
	}

	sc := FilterByCategory(cmds, SlowChanging)
	if len(sc) != 1 {
		t.Errorf("expected 1 SlowChanging command, got %d", len(sc))
	}
}

func TestFilterByCategory_Empty(t *testing.T) {
	result := FilterByCategory(nil, RealTime)
	if result != nil && len(result) != 0 {
		t.Errorf("expected nil/empty slice, got %v", result)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// keys.go tests
// ─────────────────────────────────────────────────────────────────────────────

func TestCacheKey(t *testing.T) {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	key := CacheKey(id, "system-resource")
	expected := "monitor:12345678-1234-1234-1234-123456789012:system-resource"
	if key != expected {
		t.Errorf("CacheKey = %q; want %q", key, expected)
	}
}

func TestPubSubChannel(t *testing.T) {
	id := uuid.MustParse("12345678-1234-1234-1234-123456789012")
	ch := PubSubChannel(id, "interfaces")
	expected := "monitor:12345678-1234-1234-1234-123456789012:interfaces"
	if ch != expected {
		t.Errorf("PubSubChannel = %q; want %q", ch, expected)
	}
}

func TestTTLFor(t *testing.T) {
	if TTLFor(RealTime) != TTLRealtime {
		t.Errorf("TTLFor(RealTime) = %v; want %v", TTLFor(RealTime), TTLRealtime)
	}
	if TTLFor(SlowChanging) != TTLSlowChanging {
		t.Errorf("TTLFor(SlowChanging) = %v; want %v", TTLFor(SlowChanging), TTLSlowChanging)
	}
}

func TestTTLValues(t *testing.T) {
	if TTLRealtime >= TTLSlowChanging {
		t.Error("TTLRealtime should be shorter than TTLSlowChanging")
	}
	if TTLRealtime <= 0 {
		t.Error("TTLRealtime must be positive")
	}
	if TTLSlowChanging <= 0 {
		t.Error("TTLSlowChanging must be positive")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// sink.go tests (using recordingSink mock)
// ─────────────────────────────────────────────────────────────────────────────

func TestRecordingSink_Write(t *testing.T) {
	sink := &recordingSink{}
	id := uuid.New()

	point := DataPoint{
		RouterID:   id,
		RouterHost: "192.168.1.1",
		Topic:      "system-resource",
		Category:   RealTime,
		Timestamp:  time.Now(),
		RawFields:  map[string]string{"cpu-load": "42", "free-memory": "256000"},
		Tags:       map[string]string{},
		Fields:     map[string]float64{"cpu-load": 42, "free-memory": 256000},
	}

	if err := sink.Write(context.Background(), point); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	all := sink.all()
	if len(all) != 1 {
		t.Fatalf("expected 1 recorded point, got %d", len(all))
	}
	if all[0].Topic != "system-resource" {
		t.Errorf("recorded topic = %q; want %q", all[0].Topic, "system-resource")
	}
}

func TestMultiSink_FansOut(t *testing.T) {
	sink1 := &recordingSink{}
	sink2 := &recordingSink{}

	// construct a MultiSink without zap (use a no-op logger via MultiSink directly)
	ms := &MultiSink{sinks: []DataSink{sink1, sink2}}

	point := DataPoint{
		RouterID:   uuid.New(),
		RouterHost: "192.168.88.1",
		Topic:      "interfaces",
		Category:   RealTime,
		Timestamp:  time.Now(),
		RawFields:  map[string]string{"name": "ether1", "running": "true"},
		Tags:       map[string]string{"name": "ether1", "running": "true"},
		Fields:     map[string]float64{},
	}

	if err := ms.Write(context.Background(), point); err != nil {
		t.Fatalf("MultiSink.Write() error = %v", err)
	}

	if len(sink1.all()) != 1 {
		t.Error("sink1 did not receive the data point")
	}
	if len(sink2.all()) != 1 {
		t.Error("sink2 did not receive the data point")
	}
}

func TestMultiSink_Close(t *testing.T) {
	ms := &MultiSink{sinks: []DataSink{&recordingSink{}, &recordingSink{}}}
	if err := ms.Close(); err != nil {
		t.Errorf("MultiSink.Close() unexpected error: %v", err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// manager.go lifecycle tests (no real router needed)
// ─────────────────────────────────────────────────────────────────────────────

// Test that Manager.ListRunning returns empty when nothing is started.
func TestManager_InitiallyEmpty(t *testing.T) {
	mgr := &Manager{
		collectors: make(map[uuid.UUID]*Collector),
	}
	ids := mgr.ListRunning()
	if len(ids) != 0 {
		t.Errorf("expected no running collectors, got %d", len(ids))
	}
}

// Test Manager.IsRunning on unknown ID.
func TestManager_IsRunning_Unknown(t *testing.T) {
	mgr := &Manager{
		collectors: make(map[uuid.UUID]*Collector),
	}
	if mgr.IsRunning(uuid.New()) {
		t.Error("IsRunning should return false for an unknown router")
	}
}

// Test Manager.StopAll on an empty manager does not panic.
func TestManager_StopAll_Empty(t *testing.T) {
	mgr := &Manager{
		collectors: make(map[uuid.UUID]*Collector),
		logger:     noopLogger(),
	}
	// Should not panic.
	mgr.StopAll()
}

// ─────────────────────────────────────────────────────────────────────────────
// collector.go — Collector Start/Stop idempotency (without real router)
// ─────────────────────────────────────────────────────────────────────────────

// TestCollector_StopBeforeStart verifies that Stop() is a no-op when the
// collector was never started.
func TestCollector_StopBeforeStart(t *testing.T) {
	c := &Collector{
		routerID:   uuid.New(),
		routerHost: "",
		sink:       &recordingSink{},
		logger:     noopLogger(),
		commands:   DefaultCommands(),
	}
	// Should not panic or block.
	c.Stop()
	if c.IsRunning() {
		t.Error("collector should not be running after Stop() without Start()")
	}
}

// TestCollector_IsRunning_InitialFalse verifies initial state.
func TestCollector_IsRunning_InitialFalse(t *testing.T) {
	c := &Collector{}
	if c.IsRunning() {
		t.Error("new collector should not be running")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// helpers
// ─────────────────────────────────────────────────────────────────────────────

// noopLogger returns a no-op zap logger for tests that don't need log output.
func noopLogger() *zap.Logger {
	return zap.NewNop()
}
