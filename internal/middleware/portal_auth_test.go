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

func TestPortalAuth_NoHeader(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewPortalAuthMiddleware(jwtSvc)

	r := buildTestRouter(mw.AuthenticatePortal(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing authorization header")
}

func TestPortalAuth_InvalidToken(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewPortalAuthMiddleware(jwtSvc)

	r := buildTestRouter(mw.AuthenticatePortal(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer garbage.token.here")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestPortalAuth_WrongRole(t *testing.T) {
	jwtSvc := testJWTService()
	// Generate a regular admin token (not portal)
	token, err := jwtSvc.Generate("user-123", "admin@test.com", "admin")
	assert.NoError(t, err)

	mw := NewPortalAuthMiddleware(jwtSvc)

	r := buildTestRouter(mw.AuthenticatePortal(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "not a portal token")
}

func TestPortalAuth_ValidToken(t *testing.T) {
	jwtSvc := testJWTService()
	token, err := jwtSvc.GeneratePortal("customer-789", "CUST-001")
	assert.NoError(t, err)

	mw := NewPortalAuthMiddleware(jwtSvc)

	var gotCustomerID interface{}
	handler := func(c *gin.Context) {
		gotCustomerID, _ = c.Get("customer_id")
		c.Status(http.StatusOK)
	}

	r := buildTestRouter(mw.AuthenticatePortal(), handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "customer-789", gotCustomerID)
}

func TestPortalAuth_ExpiredToken(t *testing.T) {
	expiredSvc := jwt.NewService(testSecret, -1*time.Second, -1*time.Second)
	token, err := expiredSvc.GeneratePortal("customer-expired", "CUST-EXP")
	assert.NoError(t, err)

	jwtSvc := testJWTService()
	mw := NewPortalAuthMiddleware(jwtSvc)

	r := buildTestRouter(mw.AuthenticatePortal(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestPortalAuth_InvalidFormat(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewPortalAuthMiddleware(jwtSvc)

	r := buildTestRouter(mw.AuthenticatePortal(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization header format")
}
