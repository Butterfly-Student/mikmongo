package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"mikmongo/internal/model"
	"mikmongo/internal/repository"
	"mikmongo/pkg/jwt"
)

// RedisClientInterface defines the redis operations used by AuthService
type RedisClientInterface interface {
	BlacklistToken(ctx context.Context, jti string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, jti string) (bool, error)
	SetPasswordChangedAt(ctx context.Context, userID string, t time.Time, ttl time.Duration) error
}

// AuthService handles authentication business logic
type AuthService struct {
	userRepo    repository.UserRepository
	jwtService  *jwt.Service
	redisClient RedisClientInterface
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repository.UserRepository, jwtService *jwt.Service, redisClient RedisClientInterface) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		jwtService:  jwtService,
		redisClient: redisClient,
	}
}

// LoginRequest holds login credentials
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse holds login response data
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *model.User `json:"user"`
}

// Login authenticates a user and returns a JWT token pair
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update last login
	userID, _ := uuid.Parse(user.ID)
	_ = s.userRepo.UpdateLastLogin(ctx, userID, "", time.Now())

	// Clear sensitive data
	user.PasswordHash = ""
	user.BearerKey = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// Logout blacklists a token's JTI in Redis
func (s *AuthService) Logout(ctx context.Context, tokenJTI string, remainingTTL time.Duration) error {
	return s.redisClient.BlacklistToken(ctx, tokenJTI, remainingTTL)
}

// RefreshToken validates a refresh token and returns a new token pair
func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenString string) (*LoginResponse, error) {
	claims, err := s.jwtService.Validate(refreshTokenString)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return nil, errors.New("token is not a refresh token")
	}

	// Check if refresh token is blacklisted
	if claims.ID != "" {
		blacklisted, err := s.redisClient.IsBlacklisted(ctx, claims.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check token status: %w", err)
		}
		if blacklisted {
			return nil, errors.New("refresh token has been revoked")
		}
	}

	// Check user is still active
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, errors.New("invalid user id in token")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Blacklist old refresh token (one-time use)
	if claims.ID != "" && claims.ExpiresAt != nil {
		remainingTTL := time.Until(claims.ExpiresAt.Time)
		if remainingTTL > 0 {
			_ = s.redisClient.BlacklistToken(ctx, claims.ID, remainingTTL)
		}
	}

	// Generate new token pair
	accessToken, refreshToken, err := s.jwtService.GenerateTokenPair(user.ID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Clear sensitive data
	user.PasswordHash = ""
	user.BearerKey = ""

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// ChangePasswordRequest holds change password data
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// ChangePassword changes a user's password and invalidates old tokens
func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("invalid current password")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.PasswordHash = string(hashed)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	// Invalidate all existing tokens by recording password change time
	_ = s.redisClient.SetPasswordChangedAt(ctx, userID.String(), time.Now(), 7*24*time.Hour)

	return nil
}

// CreateUser creates a new admin user
func (s *AuthService) CreateUser(ctx context.Context, user *model.User, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = string(hashed)
	return s.userRepo.Create(ctx, user)
}

// GetUser gets user by ID
func (s *AuthService) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = ""
	user.BearerKey = ""
	return user, nil
}

// ListUsers lists all users
func (s *AuthService) ListUsers(ctx context.Context, limit, offset int) ([]model.User, int64, error) {
	users, err := s.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	for i := range users {
		users[i].PasswordHash = ""
		users[i].BearerKey = ""
	}
	count, err := s.userRepo.Count(ctx)
	return users, count, err
}

// UpdateUser updates a user
func (s *AuthService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user
func (s *AuthService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.userRepo.Delete(ctx, id)
}
