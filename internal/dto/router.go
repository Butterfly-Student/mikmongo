package dto

import (
	"time"

	"mikmongo/internal/model"
)

// === REQUEST ===

// CreateRouterRequest is the request body for creating a router.
type CreateRouterRequest struct {
	Name     string  `json:"name" binding:"required"`
	Address  string  `json:"address" binding:"required"`
	Area     *string `json:"area"`
	APIPort  *int    `json:"api_port"`
	RESTPort *int    `json:"rest_port"`
	Username string  `json:"username" binding:"required"`
	Password string  `json:"password" binding:"required"`
	UseSSL   *bool   `json:"use_ssl"`
	IsMaster *bool   `json:"is_master"`
	Notes    *string `json:"notes"`
}

// ToModel converts the create request to a model.MikrotikRouter.
func (r *CreateRouterRequest) ToModel() *model.MikrotikRouter {
	m := &model.MikrotikRouter{
		Name:     r.Name,
		Address:  r.Address,
		Area:     r.Area,
		Username: r.Username,
		Notes:    r.Notes,
	}
	if r.APIPort != nil {
		m.APIPort = *r.APIPort
	}
	if r.RESTPort != nil {
		m.RESTPort = *r.RESTPort
	}
	if r.UseSSL != nil {
		m.UseSSL = *r.UseSSL
	}
	if r.IsMaster != nil {
		m.IsMaster = *r.IsMaster
	}
	return m
}

// UpdateRouterRequest is the request body for updating a router.
// All fields are pointers — only non-nil fields are applied.
type UpdateRouterRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	Area     *string `json:"area"`
	APIPort  *int    `json:"api_port"`
	RESTPort *int    `json:"rest_port"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	UseSSL   *bool   `json:"use_ssl"`
	IsMaster *bool   `json:"is_master"`
	IsActive *bool   `json:"is_active"`
	Notes    *string `json:"notes"`
}

// ApplyTo applies non-nil fields to the existing model.
func (r *UpdateRouterRequest) ApplyTo(m *model.MikrotikRouter) {
	if r.Name != nil {
		m.Name = *r.Name
	}
	if r.Address != nil {
		m.Address = *r.Address
	}
	if r.Area != nil {
		m.Area = r.Area
	}
	if r.APIPort != nil {
		m.APIPort = *r.APIPort
	}
	if r.RESTPort != nil {
		m.RESTPort = *r.RESTPort
	}
	if r.Username != nil {
		m.Username = *r.Username
	}
	if r.UseSSL != nil {
		m.UseSSL = *r.UseSSL
	}
	if r.IsMaster != nil {
		m.IsMaster = *r.IsMaster
	}
	if r.IsActive != nil {
		m.IsActive = *r.IsActive
	}
	if r.Notes != nil {
		m.Notes = r.Notes
	}
}

// === RESPONSE ===

// RouterResponse is the safe response struct.
// PasswordEncrypted and DeletedAt are excluded.
type RouterResponse struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Address    string     `json:"address"`
	Area       *string    `json:"area,omitempty"`
	APIPort    int        `json:"api_port"`
	RESTPort   int        `json:"rest_port"`
	Username   string     `json:"username"`
	UseSSL     bool       `json:"use_ssl"`
	IsMaster   bool       `json:"is_master"`
	IsActive   bool       `json:"is_active"`
	Status     string     `json:"status"`
	LastSeenAt *time.Time `json:"last_seen_at,omitempty"`
	Notes      *string    `json:"notes,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// === CONVERTERS ===

// RouterToResponse converts a model to a response DTO.
func RouterToResponse(m *model.MikrotikRouter) RouterResponse {
	return RouterResponse{
		ID:         m.ID,
		Name:       m.Name,
		Address:    m.Address,
		Area:       m.Area,
		APIPort:    m.APIPort,
		RESTPort:   m.RESTPort,
		Username:   m.Username,
		UseSSL:     m.UseSSL,
		IsMaster:   m.IsMaster,
		IsActive:   m.IsActive,
		Status:     m.Status,
		LastSeenAt: m.LastSeenAt,
		Notes:      m.Notes,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

// RoutersToResponse converts a slice of models to response DTOs.
func RoutersToResponse(ms []model.MikrotikRouter) []RouterResponse {
	result := make([]RouterResponse, len(ms))
	for i := range ms {
		result[i] = RouterToResponse(&ms[i])
	}
	return result
}
