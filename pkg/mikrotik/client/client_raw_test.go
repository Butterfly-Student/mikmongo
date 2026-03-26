package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunRaw_NotConnected(t *testing.T) {
	c := &Client{} // no connection

	results, err := c.RunRaw(context.Background(), []string{"/ip/address/print"})
	require.Error(t, err)
	assert.Nil(t, results)
	assert.Contains(t, err.Error(), "not connected")
}

func TestListenRaw_NotConnected(t *testing.T) {
	c := &Client{} // no connection

	ch := make(chan map[string]string, 10)
	cancel, err := c.ListenRaw(context.Background(), []string{"/interface/print"}, ch)
	require.Error(t, err)
	assert.Nil(t, cancel)
}
