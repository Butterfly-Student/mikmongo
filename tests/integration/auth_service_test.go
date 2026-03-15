//go:build integration

package integration

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"mikmongo/internal/model"
	"mikmongo/internal/repository/postgres"
	"mikmongo/internal/service"
	"mikmongo/pkg/jwt"
)

func newIntegrationAuthService(suite *TestSuite) (*service.AuthService, *jwt.Service) {
	repos := postgres.NewRepository(suite.DB)
	jwtSvc := jwt.NewService("test-secret-key-must-be-32chars!", 15*time.Minute, 7*24*time.Hour)
	authSvc := service.NewAuthService(repos.UserRepo, jwtSvc, suite.RedisClient)
	return authSvc, jwtSvc
}

func hashPasswordBcrypt(password string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(h)
}

func TestLogin_Integration(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	authSvc, _ := newIntegrationAuthService(suite)

	user := &model.User{
		ID:           uuid.New().String(),
		FullName:     "Admin Test",
		Email:        "admin-intg@test.com",
		PasswordHash: hashPasswordBcrypt("password123"),
		Role:         "admin",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	require.NoError(t, repos.UserRepo.Create(suite.Ctx, user))

	resp, err := authSvc.Login(suite.Ctx, "admin-intg@test.com", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
	assert.Equal(t, user.ID, resp.User.ID)
	assert.Empty(t, resp.User.PasswordHash, "password hash must be cleared in response")
}

func TestLogout_Integration(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	authSvc, jwtSvc := newIntegrationAuthService(suite)

	user := &model.User{
		ID:           uuid.New().String(),
		FullName:     "Logout Test",
		Email:        "logout@test.com",
		PasswordHash: hashPasswordBcrypt("password123"),
		Role:         "admin",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	require.NoError(t, repos.UserRepo.Create(suite.Ctx, user))

	// Login to get tokens
	resp, err := authSvc.Login(suite.Ctx, "logout@test.com", "password123")
	require.NoError(t, err)

	// Parse token to get JTI
	claims, err := jwtSvc.Validate(resp.AccessToken)
	require.NoError(t, err)
	jti := claims.ID

	// Logout with real JTI and remaining TTL
	remainingTTL := time.Until(claims.ExpiresAt.Time)
	err = authSvc.Logout(suite.Ctx, jti, remainingTTL)
	require.NoError(t, err)

	// ASSERT: JTI is now blacklisted in Redis
	blacklisted, err := suite.RedisClient.IsBlacklisted(suite.Ctx, jti)
	require.NoError(t, err)
	assert.True(t, blacklisted, "JTI must be blacklisted in Redis after logout")
}

func TestRefreshToken_Integration(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	authSvc, jwtSvc := newIntegrationAuthService(suite)

	user := &model.User{
		ID:           uuid.New().String(),
		FullName:     "Refresh Test",
		Email:        "refresh@test.com",
		PasswordHash: hashPasswordBcrypt("password123"),
		Role:         "admin",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	require.NoError(t, repos.UserRepo.Create(suite.Ctx, user))

	// Login to get token pair
	resp, err := authSvc.Login(suite.Ctx, "refresh@test.com", "password123")
	require.NoError(t, err)

	// Parse refresh token to get old JTI
	oldClaims, err := jwtSvc.Validate(resp.RefreshToken)
	require.NoError(t, err)
	oldJTI := oldClaims.ID

	// Refresh
	refreshResp, err := authSvc.RefreshToken(suite.Ctx, resp.RefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, refreshResp.AccessToken)
	assert.NotEmpty(t, refreshResp.RefreshToken)
	assert.NotEqual(t, resp.AccessToken, refreshResp.AccessToken)

	// ASSERT: Old refresh token JTI is blacklisted
	blacklisted, err := suite.RedisClient.IsBlacklisted(suite.Ctx, oldJTI)
	require.NoError(t, err)
	assert.True(t, blacklisted, "old refresh token JTI must be blacklisted after rotation")

	// ASSERT: Using old refresh token again must fail
	_, err = authSvc.RefreshToken(suite.Ctx, resp.RefreshToken)
	assert.Error(t, err, "using a revoked refresh token must return error")
}

func TestChangePassword_Integration(t *testing.T) {
	suite := SetupSuite(t)
	defer suite.TearDownSuite(t)
	defer suite.Cleanup(t)

	repos := postgres.NewRepository(suite.DB)
	authSvc, _ := newIntegrationAuthService(suite)

	userID := uuid.New()
	user := &model.User{
		ID:           userID.String(),
		FullName:     "Change Pass",
		Email:        "changepass@test.com",
		PasswordHash: hashPasswordBcrypt("old-password"),
		Role:         "admin",
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	require.NoError(t, repos.UserRepo.Create(suite.Ctx, user))

	// Change password
	err := authSvc.ChangePassword(suite.Ctx, userID, "old-password", "new-password-123")
	require.NoError(t, err)

	// ASSERT: pwd_changed key stored in Redis
	pwdChangedAt, err := suite.RedisClient.GetPasswordChangedAt(suite.Ctx, userID.String())
	require.NoError(t, err)
	assert.False(t, pwdChangedAt.IsZero(), "password changed timestamp must be stored in Redis")
	assert.WithinDuration(t, time.Now(), pwdChangedAt, 5*time.Second)

	// Should be able to log in with new password
	resp, err := authSvc.Login(suite.Ctx, "changepass@test.com", "new-password-123")
	require.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)

	// Old password should fail
	_, err = authSvc.Login(suite.Ctx, "changepass@test.com", "old-password")
	assert.Error(t, err)
}
