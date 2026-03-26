package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mikmongo/pkg/jwt"
)

func TestAgentPortalAuth_NoHeader(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAgentPortalAuthMiddleware(jwtSvc)

	r := buildTestRouter(mw.AuthenticateAgentPortal(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing authorization header")
}

func TestAgentPortalAuth_InvalidToken(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAgentPortalAuthMiddleware(jwtSvc)

	r := buildTestRouter(mw.AuthenticateAgentPortal(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer garbage.token.here")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestAgentPortalAuth_WrongRole(t *testing.T) {
	jwtSvc := testJWTService()
	// Generate a regular admin token (not agent_portal)
	token, err := jwtSvc.Generate("user-123", "admin@test.com", "admin")
	assert.NoError(t, err)

	mw := NewAgentPortalAuthMiddleware(jwtSvc)

	r := buildTestRouter(mw.AuthenticateAgentPortal(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "not an agent portal token")
}

func TestAgentPortalAuth_ValidToken(t *testing.T) {
	jwtSvc := testJWTService()
	token, err := jwtSvc.GenerateAgent("agent-456", "agentuser")
	assert.NoError(t, err)

	mw := NewAgentPortalAuthMiddleware(jwtSvc)

	var gotAgentID interface{}
	handler := func(c *gin.Context) {
		gotAgentID, _ = c.Get("agent_id")
		c.Status(http.StatusOK)
	}

	r := buildTestRouter(mw.AuthenticateAgentPortal(), handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "agent-456", gotAgentID)
}

func TestAgentPortalAuth_ExpiredToken(t *testing.T) {
	expiredSvc := jwt.NewService(testSecret, -1*time.Second, -1*time.Second)
	token, err := expiredSvc.GenerateAgent("agent-789", "expiredagent")
	assert.NoError(t, err)

	jwtSvc := testJWTService()
	mw := NewAgentPortalAuthMiddleware(jwtSvc)

	r := buildTestRouter(mw.AuthenticateAgentPortal(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestAgentPortalAuth_InvalidFormat(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAgentPortalAuthMiddleware(jwtSvc)

	r := buildTestRouter(mw.AuthenticateAgentPortal(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Token some-token")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization header format")
}
