package model

import (
	"time"

	"gorm.io/gorm"
)

// CustomerRegistration represents new customer registration requests
type CustomerRegistration struct {
	ID                 string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	FullName           string         `gorm:"type:varchar(100);not null" json:"full_name"`
	Email              *string        `gorm:"type:varchar(100)" json:"email"`
	Phone              string         `gorm:"type:varchar(20);not null" json:"phone"`
	Address            *string        `gorm:"type:text" json:"address"`
	Latitude           *float64       `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude          *float64       `gorm:"type:decimal(11,8)" json:"longitude"`
	Notes              *string        `gorm:"type:text" json:"notes"`
	BandwidthProfileID *string        `gorm:"type:uuid;index:,where:bandwidth_profile_id IS NOT NULL AND deleted_at IS NULL" json:"bandwidth_profile_id"`
	Status             string         `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending', 'approved', 'rejected');index:,where:deleted_at IS NULL" json:"status"`
	RejectionReason    *string        `gorm:"type:text" json:"rejection_reason"`
	ApprovedBy         *string        `gorm:"type:uuid" json:"approved_by"`
	ApprovedAt         *time.Time     `gorm:"type:timestamptz" json:"approved_at"`
	CustomerID         *string        `gorm:"type:uuid;index:,where:deleted_at IS NULL" json:"customer_id"`
	CreatedAt          time.Time      `gorm:"type:timestamptz;not null;default:NOW()" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"type:timestamptz;not null;default:NOW()" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	BandwidthProfile *BandwidthProfile `gorm:"foreignKey:BandwidthProfileID"`
	Approver         *User             `gorm:"foreignKey:ApprovedBy"`
	Customer         *Customer         `gorm:"foreignKey:CustomerID"`
}

func (CustomerRegistration) TableName() string {
	return "customer_registrations"
}
