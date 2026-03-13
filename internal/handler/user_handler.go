package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mikmongo/internal/model"
	"mikmongo/internal/service"
	"mikmongo/pkg/response"
)

// UserHandler handles user management HTTP requests
type UserHandler struct {
	service *service.AuthService
}

// NewUserHandler creates a new user handler
func NewUserHandler(svc *service.AuthService) *UserHandler {
	return &UserHandler{service: svc}
}

// Create handles user creation
func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		model.User
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.service.CreateUser(c.Request.Context(), &req.User, req.Password); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	req.User.PasswordHash = ""
	response.Created(c, req.User)
}

// Get handles getting a user by ID
func (h *UserHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	user, err := h.service.GetUser(c.Request.Context(), id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.OK(c, user)
}

// List handles listing users
func (h *UserHandler) List(c *gin.Context) {
	limit, offset := getPagination(c)
	users, count, err := h.service.ListUsers(c.Request.Context(), limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.WithMeta(c, http.StatusOK, users, &response.Meta{Total: count, Limit: limit, Offset: offset})
}

// Delete handles deleting a user
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}
	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.OK(c, gin.H{"message": "deleted"})
}
