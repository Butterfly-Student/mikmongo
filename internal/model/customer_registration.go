package model

import (
	"time"

	"gorm.io/gorm"
)

// CustomerRegistration represents new customer registration requests
type CustomerRegistration struct {
	ID                 string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FullName           string         `gorm:"type:varchar(100);not null"`
	Email              *string        `gorm:"type:varchar(100)"`
	Phone              string         `gorm:"type:varchar(20);not null"`
	Address            *string        `gorm:"type:text"`
	Latitude           *float64       `gorm:"type:decimal(10,8)"`
	Longitude          *float64       `gorm:"type:decimal(11,8)"`
	Notes              *string        `gorm:"type:text"`
	BandwidthProfileID *string        `gorm:"type:uuid;index:,where:bandwidth_profile_id IS NOT NULL AND deleted_at IS NULL"`
	Status             string         `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending', 'approved', 'rejected');index:,where:deleted_at IS NULL"`
	RejectionReason    *string        `gorm:"type:text"`
	ApprovedBy         *string        `gorm:"type:uuid"`
	ApprovedAt         *time.Time     `gorm:"type:timestamptz"`
	CustomerID         *string        `gorm:"type:uuid;index:,where:deleted_at IS NULL"`
	CreatedAt          time.Time      `gorm:"type:timestamptz;not null;default:NOW()"`
	UpdatedAt          time.Time      `gorm:"type:timestamptz;not null;default:NOW()"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`

	// Relationships
	BandwidthProfile *BandwidthProfile `gorm:"foreignKey:BandwidthProfileID"`
	Approver         *User             `gorm:"foreignKey:ApprovedBy"`
	Customer         *Customer         `gorm:"foreignKey:CustomerID"`
}

func (CustomerRegistration) TableName() string {
	return "customer_registrations"
}
