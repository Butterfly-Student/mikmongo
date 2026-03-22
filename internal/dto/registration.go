package dto

import (
	"time"

	"mikmongo/internal/model"
)

// === REQUEST ===

// CreateRegistrationRequest is the request body for creating a registration.
type CreateRegistrationRequest struct {
	FullName           string   `json:"full_name" binding:"required"`
	Email              *string  `json:"email"`
	Phone              string   `json:"phone" binding:"required"`
	Address            *string  `json:"address"`
	Latitude           *float64 `json:"latitude"`
	Longitude          *float64 `json:"longitude"`
	Notes              *string  `json:"notes"`
	BandwidthProfileID *string  `json:"bandwidth_profile_id"`
}

// ToModel converts the create request to a model.CustomerRegistration.
func (r *CreateRegistrationRequest) ToModel() *model.CustomerRegistration {
	return &model.CustomerRegistration{
		FullName:           r.FullName,
		Email:              r.Email,
		Phone:              r.Phone,
		Address:            r.Address,
		Latitude:           r.Latitude,
		Longitude:          r.Longitude,
		Notes:              r.Notes,
		BandwidthProfileID: r.BandwidthProfileID,
	}
}

// ApproveRegistrationRequest is the request body for approving a registration.
type ApproveRegistrationRequest struct {
	RouterID  string  `json:"router_id" binding:"required"`
	ProfileID *string `json:"profile_id"`
}

// RejectRegistrationRequest is the request body for rejecting a registration.
type RejectRegistrationRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// === RESPONSE ===

// RegistrationResponse is the safe response struct.
// DeletedAt is excluded.
type RegistrationResponse struct {
	ID                 string     `json:"id"`
	FullName           string     `json:"full_name"`
	Email              *string    `json:"email,omitempty"`
	Phone              string     `json:"phone"`
	Address            *string    `json:"address,omitempty"`
	Latitude           *float64   `json:"latitude,omitempty"`
	Longitude          *float64   `json:"longitude,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
	BandwidthProfileID *string    `json:"bandwidth_profile_id,omitempty"`
	Status             string     `json:"status"`
	RejectionReason    *string    `json:"rejection_reason,omitempty"`
	ApprovedBy         *string    `json:"approved_by,omitempty"`
	ApprovedAt         *time.Time `json:"approved_at,omitempty"`
	CustomerID         *string    `json:"customer_id,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// === CONVERTERS ===

// RegistrationToResponse converts a model to a response DTO.
func RegistrationToResponse(m *model.CustomerRegistration) RegistrationResponse {
	return RegistrationResponse{
		ID:                 m.ID,
		FullName:           m.FullName,
		Email:              m.Email,
		Phone:              m.Phone,
		Address:            m.Address,
		Latitude:           m.Latitude,
		Longitude:          m.Longitude,
		Notes:              m.Notes,
		BandwidthProfileID: m.BandwidthProfileID,
		Status:             m.Status,
		RejectionReason:    m.RejectionReason,
		ApprovedBy:         m.ApprovedBy,
		ApprovedAt:         m.ApprovedAt,
		CustomerID:         m.CustomerID,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

// RegistrationsToResponse converts a slice of models to response DTOs.
func RegistrationsToResponse(ms []model.CustomerRegistration) []RegistrationResponse {
	result := make([]RegistrationResponse, len(ms))
	for i := range ms {
		result[i] = RegistrationToResponse(&ms[i])
	}
	return result
}
