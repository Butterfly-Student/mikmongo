package dto

import (
	"time"

	"mikmongo/internal/model"
)

// === REQUEST ===

// CreateCustomerRequest is the request body for creating a customer with subscription.
type CreateCustomerRequest struct {
	FullName  string   `json:"full_name" binding:"required"`
	Phone     string   `json:"phone" binding:"required"`
	Email     *string  `json:"email"`
	Address   *string  `json:"address"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	// Subscription fields
	PlanID   string  `json:"plan_id"`
	RouterID string  `json:"router_id"`
	Username string  `json:"username"`
	Password string  `json:"password"`
	StaticIP *string `json:"static_ip"`
}

// ToCustomerModel converts the request to a model.Customer.
func (r *CreateCustomerRequest) ToCustomerModel() *model.Customer {
	return &model.Customer{
		FullName:  r.FullName,
		Phone:     r.Phone,
		Email:     r.Email,
		Address:   r.Address,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
	}
}

// UpdateCustomerRequest is the request body for updating a customer.
// All fields are pointers — only non-nil fields are applied.
type UpdateCustomerRequest struct {
	FullName     *string  `json:"full_name"`
	Username     *string  `json:"username"`
	Phone        *string  `json:"phone"`
	Email        *string  `json:"email"`
	Address      *string  `json:"address"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	IDCardNumber *string  `json:"id_card_number"`
	Notes        *string  `json:"notes"`
	Tags         *string  `json:"tags"`
	IsActive     *bool    `json:"is_active"`
}

// ApplyTo applies non-nil fields to the existing model.
func (r *UpdateCustomerRequest) ApplyTo(m *model.Customer) {
	if r.FullName != nil {
		m.FullName = *r.FullName
	}
	if r.Username != nil {
		m.Username = r.Username
	}
	if r.Phone != nil {
		m.Phone = *r.Phone
	}
	if r.Email != nil {
		m.Email = r.Email
	}
	if r.Address != nil {
		m.Address = r.Address
	}
	if r.Latitude != nil {
		m.Latitude = r.Latitude
	}
	if r.Longitude != nil {
		m.Longitude = r.Longitude
	}
	if r.IDCardNumber != nil {
		m.IDCardNumber = r.IDCardNumber
	}
	if r.Notes != nil {
		m.Notes = r.Notes
	}
	if r.Tags != nil {
		m.Tags = r.Tags
	}
	if r.IsActive != nil {
		m.IsActive = *r.IsActive
	}
}

// === RESPONSE ===

// CustomerResponse is the safe response struct.
// PortalPasswordHash, PortalLastLogin, DeletedAt are excluded.
type CustomerResponse struct {
	ID           string     `json:"id"`
	CustomerCode string     `json:"customer_code"`
	FullName     string     `json:"full_name"`
	Username     *string    `json:"username,omitempty"`
	Email        *string    `json:"email,omitempty"`
	Phone        string     `json:"phone"`
	IDCardNumber *string    `json:"id_card_number,omitempty"`
	Address      *string    `json:"address,omitempty"`
	Latitude     *float64   `json:"latitude,omitempty"`
	Longitude    *float64   `json:"longitude,omitempty"`
	IsActive     bool       `json:"is_active"`
	Notes        *string    `json:"notes,omitempty"`
	Tags         *string    `json:"tags,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// === CONVERTERS ===

// CustomerToResponse converts a model to a response DTO.
func CustomerToResponse(m *model.Customer) CustomerResponse {
	return CustomerResponse{
		ID:           m.ID,
		CustomerCode: m.CustomerCode,
		FullName:     m.FullName,
		Username:     m.Username,
		Email:        m.Email,
		Phone:        m.Phone,
		IDCardNumber: m.IDCardNumber,
		Address:      m.Address,
		Latitude:     m.Latitude,
		Longitude:    m.Longitude,
		IsActive:     m.IsActive,
		Notes:        m.Notes,
		Tags:         m.Tags,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

// CustomersToResponse converts a slice of models to response DTOs.
func CustomersToResponse(ms []model.Customer) []CustomerResponse {
	result := make([]CustomerResponse, len(ms))
	for i := range ms {
		result[i] = CustomerToResponse(&ms[i])
	}
	return result
}
