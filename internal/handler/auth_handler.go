package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/service"
	"mikmongo/pkg/jwt"
	"mikmongo/pkg/response"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	service    *service.AuthService
	jwtService *jwt.Service
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(svc *service.AuthService, jwtService *jwt.Service) *AuthHandler {
	return &AuthHandler{
		service:    svc,
		jwtService: jwtService,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	response.OK(c, result)
}

// Logout handles user logout by blacklisting the current token
func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		response.BadRequest(c, "missing token")
		return
	}

	tokenString := parts[1]
	claims, err := h.jwtService.Validate(tokenString)
	if err != nil {
		response.Unauthorized(c, "invalid token")
		return
	}

	if claims.ID == "" {
		response.BadRequest(c, "token has no JTI")
		return
	}

	// Calculate remaining TTL
	var remainingTTL time.Duration
	if claims.ExpiresAt != nil {
		remainingTTL = time.Until(claims.ExpiresAt.Time)
	}
	if remainingTTL <= 0 {
		// Token already expired, no need to blacklist
		response.OK(c, gin.H{"message": "logged out"})
		return
	}

	if err := h.service.Logout(c.Request.Context(), claims.ID, remainingTTL); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "logged out"})
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.service.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.OK(c, result)
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}
	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.service.ChangePassword(c.Request.Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "password changed"})
}

// GetMe returns the authenticated user's profile
func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}
	user, err := h.service.GetUser(c.Request.Context(), userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.OK(c, user)
}
