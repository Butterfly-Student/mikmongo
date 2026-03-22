package dto

import (
	"time"

	"mikmongo/internal/model"
)

// === REQUEST ===

// CreateUserRequest is the request body for creating a user.
type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ToModel converts the create request to a model.User.
func (r *CreateUserRequest) ToModel() *model.User {
	return &model.User{
		FullName: r.FullName,
		Email:    r.Email,
		Phone:    r.Phone,
		Role:     r.Role,
	}
}

// UpdateUserRequest is the request body for updating a user.
// All fields are pointers — only non-nil fields are applied.
type UpdateUserRequest struct {
	FullName *string `json:"full_name"`
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
}

// ApplyTo applies non-nil fields to the existing model.
func (r *UpdateUserRequest) ApplyTo(m *model.User) {
	if r.FullName != nil {
		m.FullName = *r.FullName
	}
	if r.Email != nil {
		m.Email = *r.Email
	}
	if r.Phone != nil {
		m.Phone = *r.Phone
	}
	if r.Role != nil {
		m.Role = *r.Role
	}
	if r.IsActive != nil {
		m.IsActive = *r.IsActive
	}
}

// === RESPONSE ===

// UserResponse is the safe response struct.
// PasswordHash, BearerKey, LastIP, and DeletedAt are excluded.
type UserResponse struct {
	ID        string     `json:"id"`
	FullName  string     `json:"full_name"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Role      string     `json:"role"`
	IsActive  bool       `json:"is_active"`
	LastLogin *time.Time `json:"last_login,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// === CONVERTERS ===

// UserToResponse converts a model to a response DTO.
func UserToResponse(m *model.User) UserResponse {
	return UserResponse{
		ID:        m.ID,
		FullName:  m.FullName,
		Email:     m.Email,
		Phone:     m.Phone,
		Role:      m.Role,
		IsActive:  m.IsActive,
		LastLogin: m.LastLogin,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// UsersToResponse converts a slice of models to response DTOs.
func UsersToResponse(ms []model.User) []UserResponse {
	result := make([]UserResponse, len(ms))
	for i := range ms {
		result[i] = UserToResponse(&ms[i])
	}
	return result
}
