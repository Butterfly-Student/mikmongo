//go:build integration

package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"mikmongo/internal/middleware"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
	"mikmongo/pkg/jwt"
)

// newE2EAuthRouter builds a minimal Gin router with real auth middleware for E2E testing.
// Uses real Redis and real JWT service — no mocks.
func newE2EAuthRouter(suite *TestSuite, jwtSvc *jwt.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mw := middleware.NewAuthMiddleware(jwtSvc, suite.RedisClient)
	r.GET("/protected", mw.Authenticate(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})
	return r
}

func createTestUserModel(t *testing.T, suite *TestSuite, email string) *model.User {
	t.Helper()
	h, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	return &model.User{
		ID:           uuid.New().String(),
		FullName:     "E2E Test User",
		Email:        email,
		PasswordHash: string(h),
		Role:         "admin",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// TestE2E_LogoutBlacklistsToken verifies that an access token is rejected after logout:
//
//	Login → GET /protected → 200
//	→ Logout (blacklist JTI in Redis)
//	→ GET /protected with same token → 401
func TestE2E_LogoutBlacklistsToken(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	jwtSvc := jwt.NewService("test-secret-key-must-be-32chars!", 15*time.Minute, 7*24*time.Hour)
	authSvc := service.NewAuthService(repos.UserRepo, jwtSvc, suite.RedisClient)
	router := newE2EAuthRouter(suite, jwtSvc)

	user := createTestUserModel(t, suite, "e2e-logout@test.com")
	require.NoError(t, repos.UserRepo.Create(suite.Ctx, user))

	// Login
	resp, err := authSvc.Login(suite.Ctx, "e2e-logout@test.com", "password123")
	require.NoError(t, err)

	// Protected endpoint → 200
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+resp.AccessToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "token should be valid before logout")

	// Logout
	claims, err := jwtSvc.Validate(resp.AccessToken)
	require.NoError(t, err)
	require.NoError(t, authSvc.Logout(suite.Ctx, claims.ID, time.Until(claims.ExpiresAt.Time)))

	// Same token → 401
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req2.Header.Set("Authorization", "Bearer "+resp.AccessToken)
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code, "revoked token must be rejected")
}

// TestE2E_RefreshTokenRotation_OldTokenRejected verifies that after a refresh,
// the old refresh token cannot be used again (service-level check).
func TestE2E_RefreshTokenRotation_OldTokenRejected(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	jwtSvc := jwt.NewService("test-secret-key-must-be-32chars!", 15*time.Minute, 7*24*time.Hour)
	authSvc := service.NewAuthService(repos.UserRepo, jwtSvc, suite.RedisClient)

	user := createTestUserModel(t, suite, "e2e-refresh@test.com")
	require.NoError(t, repos.UserRepo.Create(suite.Ctx, user))

	// Login
	resp, err := authSvc.Login(suite.Ctx, "e2e-refresh@test.com", "password123")
	require.NoError(t, err)

	oldRefreshToken := resp.RefreshToken

	// Refresh — old refresh JTI gets blacklisted, new tokens issued
	newResp, err := authSvc.RefreshToken(suite.Ctx, oldRefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, newResp.AccessToken)
	assert.NotEmpty(t, newResp.RefreshToken)

	// Using old refresh token again must fail
	_, err = authSvc.RefreshToken(suite.Ctx, oldRefreshToken)
	assert.Error(t, err, "using a revoked refresh token must return error")
}

// TestE2E_PasswordChange_OldTokenInvalidated verifies that tokens issued before
// a password change are rejected by the auth middleware.
//
//	Login → GET /protected → 200
//	→ ChangePassword (stores pwd_changed:<userID> in Redis)
//	→ GET /protected with old token → 401
func TestE2E_PasswordChange_OldTokenInvalidated(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	jwtSvc := jwt.NewService("test-secret-key-must-be-32chars!", 15*time.Minute, 7*24*time.Hour)
	authSvc := service.NewAuthService(repos.UserRepo, jwtSvc, suite.RedisClient)
	router := newE2EAuthRouter(suite, jwtSvc)

	user := createTestUserModel(t, suite, "e2e-pwdchange@test.com")
	require.NoError(t, repos.UserRepo.Create(suite.Ctx, user))

	// Login → get old access token
	resp, err := authSvc.Login(suite.Ctx, "e2e-pwdchange@test.com", "password123")
	require.NoError(t, err)

	// Protected endpoint → 200 (token is valid before password change)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+resp.AccessToken)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "token should be valid before password change")

	// Give at least 1 second so IssuedAt < pwdChangedAt
	time.Sleep(1 * time.Second)

	// Change password
	userID, err := uuid.Parse(user.ID)
	require.NoError(t, err)
	err = authSvc.ChangePassword(suite.Ctx, userID, "password123", "newpassword-456")
	require.NoError(t, err)

	// Old token → 401 (issued before password change)
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req2.Header.Set("Authorization", "Bearer "+resp.AccessToken)
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusUnauthorized, w2.Code, "token issued before password change must be rejected")
}
