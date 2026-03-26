package notification

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"mikmongo/pkg/gowa"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "strip spaces", input: "081 234 567 890", want: "6281234567890"},
		{name: "strip plus prefix", input: "+6281234567890", want: "6281234567890"},
		{name: "leading zero becomes 62", input: "081234567890", want: "6281234567890"},
		{name: "already international", input: "6281234567890", want: "6281234567890"},
		{name: "empty string", input: "", want: ""},
		{name: "plus and spaces combined", input: "+62 812 345 678", want: "62812345678"},
		{name: "no leading zero no plus", input: "81234567890", want: "81234567890"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizePhone(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func newTestGoWAClient(ts *httptest.Server, groupID string) *GoWAClient {
	client := gowa.New(&gowa.Config{
		BaseURL:  ts.URL,
		Username: "testuser",
		Password: "testpass",
		DeviceID: "device-1",
		Timeout:  10 * time.Second,
	})
	return NewGoWAClient(client, groupID)
}

func TestGoWAClient_SendMessage_Success(t *testing.T) {
	var receivedBody []byte
	var receivedAuth string
	var receivedPath string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path

		user, pass, ok := r.BasicAuth()
		if ok {
			receivedAuth = user + ":" + pass
		}

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		receivedBody = body

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := gowa.SendResponse{
			Code:    "200",
			Message: "success",
			Results: gowa.SendResult{MessageID: "msg-123", Status: "sent"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	client := newTestGoWAClient(ts, "")

	ctx := context.Background()
	err := client.SendMessage(ctx, "081234567890", "Hello test")
	require.NoError(t, err)

	assert.Equal(t, "/send/message", receivedPath)
	assert.Equal(t, "testuser:testpass", receivedAuth)

	var reqBody gowa.SendMessageRequest
	err = json.Unmarshal(receivedBody, &reqBody)
	require.NoError(t, err)
	assert.Equal(t, "6281234567890", reqBody.Phone, "phone should be normalized")
	assert.Equal(t, "Hello test", reqBody.Message)
}

func TestGoWAClient_SendMessage_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"code":"500","message":"internal error"}`))
	}))
	defer ts.Close()

	client := newTestGoWAClient(ts, "")

	ctx := context.Background()
	err := client.SendMessage(ctx, "081234567890", "Hello")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gowa send failed")
}

func TestGoWAClient_SendMessage_ContextCancel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := newTestGoWAClient(ts, "")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	err := client.SendMessage(ctx, "081234567890", "Hello")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gowa send failed")
}

func TestGoWAClient_SendGroupMessage_Success(t *testing.T) {
	var receivedBody []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = body

		w.Header().Set("Content-Type", "application/json")
		resp := gowa.SendResponse{
			Code:    "200",
			Message: "success",
			Results: gowa.SendResult{MessageID: "msg-456", Status: "sent"},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	groupJID := "120363347168689807@g.us"
	client := newTestGoWAClient(ts, groupJID)

	ctx := context.Background()
	err := client.SendGroupMessage(ctx, "Hello group")
	require.NoError(t, err)

	var reqBody gowa.SendMessageRequest
	err = json.Unmarshal(receivedBody, &reqBody)
	require.NoError(t, err)
	assert.Equal(t, groupJID, reqBody.Phone)
	assert.Equal(t, "Hello group", reqBody.Message)
}

func TestGoWAClient_SendGroupMessage_NoGroupID_Skips(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not be called when no group ID configured")
	}))
	defer ts.Close()

	client := newTestGoWAClient(ts, "") // no group ID

	ctx := context.Background()
	err := client.SendGroupMessage(ctx, "Hello group")
	require.NoError(t, err)
}
