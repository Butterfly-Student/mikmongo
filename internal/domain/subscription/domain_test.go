package subscription

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/model"
)

func TestValidateStatusTransition(t *testing.T) {
	d := NewDomain()

	validCases := []struct {
		from, to string
	}{
		{"pending", "active"},
		{"pending", "terminated"},
		{"active", "suspended"},
		{"active", "isolated"},
		{"active", "expired"},
		{"active", "terminated"},
		{"suspended", "active"},
		{"suspended", "isolated"},
		{"suspended", "terminated"},
		{"isolated", "active"},
		{"isolated", "suspended"},
		{"isolated", "terminated"},
		{"expired", "active"},
		{"expired", "terminated"},
	}

	for _, tc := range validCases {
		t.Run(tc.from+"→"+tc.to+" valid", func(t *testing.T) {
			err := d.ValidateStatusTransition(tc.from, tc.to)
			assert.NoError(t, err)
		})
	}

	invalidCases := []struct {
		from, to string
	}{
		{"terminated", "active"},
		{"terminated", "pending"},
		{"pending", "suspended"},
		{"pending", "isolated"},
		{"active", "pending"},
		{"suspended", "pending"},
		{"isolated", "pending"},
		{"expired", "suspended"},
	}

	for _, tc := range invalidCases {
		t.Run(tc.from+"→"+tc.to+" invalid", func(t *testing.T) {
			err := d.ValidateStatusTransition(tc.from, tc.to)
			assert.Error(t, err)
		})
	}

	t.Run("unknown status → error", func(t *testing.T) {
		err := d.ValidateStatusTransition("unknown", "active")
		assert.Error(t, err)
	})
}

func TestCanActivate(t *testing.T) {
	d := NewDomain()

	t.Run("pending → can activate", func(t *testing.T) {
		sub := &model.Subscription{Status: "pending"}
		assert.NoError(t, d.CanActivate(sub))
	})

	t.Run("suspended → can activate", func(t *testing.T) {
		sub := &model.Subscription{Status: "suspended"}
		assert.NoError(t, d.CanActivate(sub))
	})

	t.Run("active → cannot activate", func(t *testing.T) {
		sub := &model.Subscription{Status: "active"}
		assert.Error(t, d.CanActivate(sub))
	})

	t.Run("isolated → cannot activate via CanActivate", func(t *testing.T) {
		sub := &model.Subscription{Status: "isolated"}
		assert.Error(t, d.CanActivate(sub))
	})

	t.Run("terminated → cannot activate", func(t *testing.T) {
		sub := &model.Subscription{Status: "terminated"}
		assert.Error(t, d.CanActivate(sub))
	})
}

func TestCanSuspend(t *testing.T) {
	d := NewDomain()

	t.Run("active → can suspend", func(t *testing.T) {
		sub := &model.Subscription{Status: "active"}
		assert.NoError(t, d.CanSuspend(sub))
	})

	t.Run("isolated → can suspend", func(t *testing.T) {
		sub := &model.Subscription{Status: "isolated"}
		assert.NoError(t, d.CanSuspend(sub))
	})

	t.Run("pending → cannot suspend", func(t *testing.T) {
		sub := &model.Subscription{Status: "pending"}
		assert.Error(t, d.CanSuspend(sub))
	})

	t.Run("terminated → cannot suspend", func(t *testing.T) {
		sub := &model.Subscription{Status: "terminated"}
		assert.Error(t, d.CanSuspend(sub))
	})
}

func TestCanIsolate(t *testing.T) {
	d := NewDomain()

	t.Run("active → can isolate", func(t *testing.T) {
		sub := &model.Subscription{Status: "active"}
		assert.NoError(t, d.CanIsolate(sub))
	})

	t.Run("suspended → can isolate", func(t *testing.T) {
		sub := &model.Subscription{Status: "suspended"}
		assert.NoError(t, d.CanIsolate(sub))
	})

	t.Run("pending → cannot isolate", func(t *testing.T) {
		sub := &model.Subscription{Status: "pending"}
		assert.Error(t, d.CanIsolate(sub))
	})

	t.Run("isolated → cannot isolate again", func(t *testing.T) {
		sub := &model.Subscription{Status: "isolated"}
		assert.Error(t, d.CanIsolate(sub))
	})
}

func TestCanRestore(t *testing.T) {
	d := NewDomain()

	t.Run("isolated → can restore", func(t *testing.T) {
		sub := &model.Subscription{Status: "isolated"}
		assert.NoError(t, d.CanRestore(sub))
	})

	t.Run("active → cannot restore", func(t *testing.T) {
		sub := &model.Subscription{Status: "active"}
		assert.Error(t, d.CanRestore(sub))
	})

	t.Run("suspended → cannot restore", func(t *testing.T) {
		sub := &model.Subscription{Status: "suspended"}
		assert.Error(t, d.CanRestore(sub))
	})
}

func TestCanTerminate(t *testing.T) {
	d := NewDomain()

	t.Run("active → can terminate", func(t *testing.T) {
		sub := &model.Subscription{Status: "active"}
		assert.NoError(t, d.CanTerminate(sub))
	})

	t.Run("pending → can terminate", func(t *testing.T) {
		sub := &model.Subscription{Status: "pending"}
		assert.NoError(t, d.CanTerminate(sub))
	})

	t.Run("suspended → can terminate", func(t *testing.T) {
		sub := &model.Subscription{Status: "suspended"}
		assert.NoError(t, d.CanTerminate(sub))
	})

	t.Run("already terminated → error", func(t *testing.T) {
		sub := &model.Subscription{Status: "terminated"}
		assert.Error(t, d.CanTerminate(sub))
	})
}

func TestIsExpired(t *testing.T) {
	d := NewDomain()
	now := time.Now()

	t.Run("expiry date in the past → expired", func(t *testing.T) {
		past := now.AddDate(0, 0, -1)
		sub := &model.Subscription{ExpiryDate: &past}
		assert.True(t, d.IsExpired(sub, now))
	})

	t.Run("expiry date in the future → not expired", func(t *testing.T) {
		future := now.AddDate(0, 0, 1)
		sub := &model.Subscription{ExpiryDate: &future}
		assert.False(t, d.IsExpired(sub, now))
	})

	t.Run("expiry date nil → not expired", func(t *testing.T) {
		sub := &model.Subscription{ExpiryDate: nil}
		assert.False(t, d.IsExpired(sub, now))
	})
}

func TestValidateCredentials(t *testing.T) {
	d := NewDomain()

	t.Run("valid credentials", func(t *testing.T) {
		assert.NoError(t, d.ValidateCredentials("validuser", "password123"))
	})

	t.Run("username too short (< 3 chars)", func(t *testing.T) {
		err := d.ValidateCredentials("ab", "password123")
		assert.Error(t, err)
	})

	t.Run("username too long (> 100 chars)", func(t *testing.T) {
		longUsername := strings.Repeat("a", 101)
		err := d.ValidateCredentials(longUsername, "password123")
		assert.Error(t, err)
	})

	t.Run("password too short (< 6 chars)", func(t *testing.T) {
		err := d.ValidateCredentials("validuser", "short")
		assert.Error(t, err)
	})

	t.Run("username exactly 3 chars → valid", func(t *testing.T) {
		assert.NoError(t, d.ValidateCredentials("abc", "password123"))
	})

	t.Run("username exactly 100 chars → valid", func(t *testing.T) {
		username100 := strings.Repeat("a", 100)
		assert.NoError(t, d.ValidateCredentials(username100, "password123"))
	})

	t.Run("password exactly 6 chars → valid", func(t *testing.T) {
		assert.NoError(t, d.ValidateCredentials("validuser", "sixchr"))
	})
}

func TestGeneratePassword(t *testing.T) {
	d := NewDomain()

	t.Run("correct length", func(t *testing.T) {
		pwd, err := d.GeneratePassword(12)
		require.NoError(t, err)
		assert.Len(t, pwd, 12)
	})

	t.Run("only alphanumeric characters", func(t *testing.T) {
		pwd, err := d.GeneratePassword(50)
		require.NoError(t, err)
		for _, c := range pwd {
			assert.True(t, (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'),
				"character %c is not alphanumeric", c)
		}
	})

	t.Run("different each call (entropy)", func(t *testing.T) {
		pwd1, err1 := d.GeneratePassword(16)
		pwd2, err2 := d.GeneratePassword(16)
		require.NoError(t, err1)
		require.NoError(t, err2)
		// Extremely unlikely to be equal with 16-char random passwords
		assert.NotEqual(t, pwd1, pwd2)
	})

	t.Run("zero length → error", func(t *testing.T) {
		_, err := d.GeneratePassword(0)
		assert.Error(t, err)
	})

	t.Run("negative length → error", func(t *testing.T) {
		_, err := d.GeneratePassword(-1)
		assert.Error(t, err)
	})
}
