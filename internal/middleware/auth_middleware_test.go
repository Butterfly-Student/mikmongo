package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mikmongo/pkg/jwt"
	"mikmongo/pkg/redis"
)

func init() { gin.SetMode(gin.TestMode) }

const testSecret = "test-secret-key-for-middleware-tests-32chars"

func testJWTService() *jwt.Service {
	return jwt.NewService(testSecret, 1*time.Hour, 24*time.Hour)
}

func expiredJWTService() *jwt.Service {
	return jwt.NewService(testSecret, -1*time.Second, -1*time.Second)
}

func testRedisClient(t *testing.T) *redis.Client {
	t.Helper()
	host := os.Getenv("TEST_REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	portStr := os.Getenv("TEST_REDIS_PORT")
	if portStr == "" {
		portStr = "6379"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Skip("invalid TEST_REDIS_PORT, skipping test")
	}
	dbStr := os.Getenv("TEST_REDIS_DB")
	if dbStr == "" {
		dbStr = "15"
	}
	db, err := strconv.Atoi(dbStr)
	if err != nil {
		t.Skip("invalid TEST_REDIS_DB, skipping test")
	}

	client := redis.NewClient(redis.Options{
		Host: host,
		Port: port,
		DB:   db,
	})
	if err := client.Ping(context.Background()); err != nil {
		t.Skipf("Redis not available at %s:%d, skipping test: %v", host, port, err)
	}
	return client
}

func buildTestRouter(middleware gin.HandlerFunc, handler gin.HandlerFunc) *gin.Engine {
	r := gin.New()
	r.GET("/test", middleware, handler)
	return r
}

func okHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}

// --- Authenticate tests (no Redis needed) ---

func TestAuthenticate_NoHeader(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAuthMiddleware(jwtSvc, nil)

	r := buildTestRouter(mw.Authenticate(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing authorization header")
}

func TestAuthenticate_InvalidFormat_NoBearer(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAuthMiddleware(jwtSvc, nil)

	r := buildTestRouter(mw.Authenticate(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Token some-token")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization header format")
}

func TestAuthenticate_InvalidFormat_NoParts(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAuthMiddleware(jwtSvc, nil)

	r := buildTestRouter(mw.Authenticate(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "justonepart")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid authorization header format")
}

func TestAuthenticate_InvalidToken(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAuthMiddleware(jwtSvc, nil)

	r := buildTestRouter(mw.Authenticate(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.string")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestAuthenticate_ExpiredToken(t *testing.T) {
	expiredSvc := expiredJWTService()
	token, err := expiredSvc.Generate("user-123", "test@test.com", "admin")
	assert.NoError(t, err)

	// Use a fresh service with the same secret to validate (it will see the token as expired)
	jwtSvc := testJWTService()
	mw := NewAuthMiddleware(jwtSvc, nil)

	r := buildTestRouter(mw.Authenticate(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

func TestAuthenticate_RefreshTokenRejected(t *testing.T) {
	jwtSvc := testJWTService()
	refreshToken, err := jwtSvc.GenerateRefreshToken("user-123", "test@test.com", "admin")
	assert.NoError(t, err)

	mw := NewAuthMiddleware(jwtSvc, nil)

	r := buildTestRouter(mw.Authenticate(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+refreshToken)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "refresh token not allowed")
}

func TestAuthenticate_WrongSecret(t *testing.T) {
	otherSvc := jwt.NewService("a-completely-different-secret-key-32chars!", 1*time.Hour, 24*time.Hour)
	token, err := otherSvc.Generate("user-123", "test@test.com", "admin")
	assert.NoError(t, err)

	jwtSvc := testJWTService()
	mw := NewAuthMiddleware(jwtSvc, nil)

	r := buildTestRouter(mw.Authenticate(), okHandler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "invalid token")
}

// --- Authenticate test requiring Redis ---

func TestAuthenticate_ValidToken_SetsContext(t *testing.T) {
	redisClient := testRedisClient(t)
	defer redisClient.Close()

	jwtSvc := testJWTService()
	token, err := jwtSvc.Generate("user-123", "test@test.com", "admin")
	assert.NoError(t, err)

	mw := NewAuthMiddleware(jwtSvc, redisClient)

	var gotUserID, gotEmail, gotRole interface{}
	handler := func(c *gin.Context) {
		gotUserID, _ = c.Get("user_id")
		gotEmail, _ = c.Get("email")
		gotRole, _ = c.Get("role")
		c.Status(http.StatusOK)
	}

	r := buildTestRouter(mw.Authenticate(), handler)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "user-123", gotUserID)
	assert.Equal(t, "test@test.com", gotEmail)
	assert.Equal(t, "admin", gotRole)
}

// --- Authorize tests (no Redis needed, role set manually) ---

func TestAuthorize_NoRole(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAuthMiddleware(jwtSvc, nil)

	r := gin.New()
	r.GET("/test", mw.Authorize("admin"), okHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "role not found")
}

func TestAuthorize_WrongRole(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAuthMiddleware(jwtSvc, nil)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.Set("role", "viewer")
		c.Next()
	}, mw.Authorize("admin"), okHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	assert.Contains(t, w.Body.String(), "insufficient permissions")
}

func TestAuthorize_MatchingRole(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAuthMiddleware(jwtSvc, nil)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	}, mw.Authorize("admin"), okHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthorize_MultipleRoles(t *testing.T) {
	jwtSvc := testJWTService()
	mw := NewAuthMiddleware(jwtSvc, nil)

	r := gin.New()
	r.GET("/test", func(c *gin.Context) {
		c.Set("role", "editor")
		c.Next()
	}, mw.Authorize("admin", "editor", "viewer"), okHandler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
