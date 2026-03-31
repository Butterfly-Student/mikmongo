package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── ListenManyArgsContext unit tests (no real router needed) ─────────────────

// TestListenManyArgsContext_NotConnected verifies that calling ListenManyArgsContext
// without an active connection returns an immediate error and no channel.
func TestListenManyArgsContext_NotConnected(t *testing.T) {
	c := &Client{} // no connection

	commands := [][]string{
		{"/interface/print", "=follow="},
		{"/system/resource/print", "=follow="},
	}

	ch, err := c.ListenManyArgsContext(context.Background(), commands, 32)
	require.Error(t, err)
	assert.Nil(t, ch)
	assert.Contains(t, err.Error(), "not connected")
}

// TestListenManyArgsContext_EmptyCommands verifies that an empty command list
// returns an error immediately (before even checking the connection).
func TestListenManyArgsContext_EmptyCommands(t *testing.T) {
	c := &Client{} // no connection

	ch, err := c.ListenManyArgsContext(context.Background(), nil, 32)
	require.Error(t, err)
	assert.Nil(t, ch)
}

// TestListenManyArgsContext_EmptySlice uses an explicitly empty (not nil) slice.
func TestListenManyArgsContext_EmptySlice(t *testing.T) {
	c := &Client{}

	ch, err := c.ListenManyArgsContext(context.Background(), [][]string{}, 32)
	require.Error(t, err)
	assert.Nil(t, ch)
}

// ─── StreamEvent struct tests ──────────────────────────────────────────────────

// TestStreamEvent_ZeroValue ensures the zero value of StreamEvent is usable
// and does not panic on access.
func TestStreamEvent_ZeroValue(t *testing.T) {
	var ev StreamEvent
	assert.Equal(t, 0, ev.Index)
	assert.Nil(t, ev.Args)
	assert.Nil(t, ev.Map)
	assert.Nil(t, ev.Err)
}

// TestStreamEvent_ErrorField verifies constructing an error event.
func TestStreamEvent_ErrorField(t *testing.T) {
	args := []string{"/interface/print", "=follow="}
	ev := StreamEvent{
		Index: 2,
		Args:  args,
		Err:   context.DeadlineExceeded,
	}
	assert.Equal(t, 2, ev.Index)
	assert.Equal(t, args, ev.Args)
	assert.ErrorIs(t, ev.Err, context.DeadlineExceeded)
	assert.Nil(t, ev.Map)
}

// TestStreamEvent_MapField verifies constructing a data event.
func TestStreamEvent_MapField(t *testing.T) {
	m := map[string]string{
		"name":     "ether1",
		"running":  "true",
		"rx-byte":  "123456",
		"tx-byte":  "654321",
	}
	ev := StreamEvent{Index: 0, Args: []string{"/interface/print"}, Map: m}
	assert.Equal(t, "ether1", ev.Map["name"])
	assert.Equal(t, "123456", ev.Map["rx-byte"])
	assert.Nil(t, ev.Err)
}

// ─── ListenArgs / ListenArgsContext unit tests ────────────────────────────────

func TestListenArgs_NotConnected(t *testing.T) {
	c := &Client{}

	lr, err := c.ListenArgs([]string{"/interface/print", "=follow="})
	require.Error(t, err)
	assert.Nil(t, lr)
	assert.Contains(t, err.Error(), "not connected")
}

func TestListenArgsContext_NotConnected(t *testing.T) {
	c := &Client{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	lr, err := c.ListenArgsContext(ctx, []string{"/interface/print", "=follow="})
	require.Error(t, err)
	assert.Nil(t, lr)
}

func TestListenArgsQueue_NotConnected(t *testing.T) {
	c := &Client{}

	lr, err := c.ListenArgsQueue([]string{"/interface/print", "=follow="}, 64)
	require.Error(t, err)
	assert.Nil(t, lr)
}
