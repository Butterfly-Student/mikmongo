package model

import (
	"time"

	"gorm.io/gorm"
)

// Customer represents ISP customer identity data
type Customer struct {
	ID                 string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerCode       string         `gorm:"type:varchar(50);unique;not null;index"`
	FullName           string         `gorm:"type:varchar(100);not null"`
	Email              *string        `gorm:"type:varchar(100);unique;index:,where:email IS NOT NULL"`
	Username           *string        `gorm:"type:varchar(100);uniqueIndex:,where:username IS NOT NULL AND deleted_at IS NULL"`
	Phone              string         `gorm:"type:varchar(20);not null;index"`
	IDCardNumber       *string        `gorm:"type:varchar(30)"`
	Address            *string        `gorm:"type:text"`
	Latitude           *float64       `gorm:"type:decimal(10,8)"`
	Longitude          *float64       `gorm:"type:decimal(11,8)"`
	IsActive           bool           `gorm:"type:boolean;default:true;index:,where:deleted_at IS NULL"`
	PortalPasswordHash *string        `gorm:"type:varchar(255)"`
	PortalLastLogin    *time.Time     `gorm:"type:timestamptz"`
	Notes              *string        `gorm:"type:text"`
	Tags               *string        `gorm:"type:jsonb"`
	CreatedBy          *string        `gorm:"type:uuid"`
	UpdatedBy          *string        `gorm:"type:uuid"`
	CreatedAt          time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`

	// Relationships
	Creator       *User                  `gorm:"foreignKey:CreatedBy"`
	Updater       *User                  `gorm:"foreignKey:UpdatedBy"`
	Subscriptions []Subscription         `gorm:"foreignKey:CustomerID"`
	Invoices      []Invoice              `gorm:"foreignKey:CustomerID"`
	Payments      []Payment              `gorm:"foreignKey:CustomerID"`
	Registrations []CustomerRegistration `gorm:"foreignKey:CustomerID"`
}

func (Customer) TableName() string {
	return "customers"
}
