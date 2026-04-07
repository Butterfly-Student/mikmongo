package model

import (
	"time"

	"gorm.io/gorm"
)

// Customer represents ISP customer identity data
type Customer struct {
	ID                 string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CustomerCode       string         `gorm:"type:varchar(50);unique;not null;index" json:"customer_code"`
	FullName           string         `gorm:"type:varchar(100);not null" json:"full_name"`
	Email              *string        `gorm:"type:varchar(100);unique;index:,where:email IS NOT NULL" json:"email"`
	Username           *string        `gorm:"type:varchar(100);uniqueIndex:,where:username IS NOT NULL AND deleted_at IS NULL" json:"username"`
	Phone              string         `gorm:"type:varchar(20);not null;index" json:"phone"`
	IDCardNumber       *string        `gorm:"type:varchar(30)" json:"id_card_number"`
	Address            *string        `gorm:"type:text" json:"address"`
	Latitude           *float64       `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude          *float64       `gorm:"type:decimal(11,8)" json:"longitude"`
	IsActive           bool           `gorm:"type:boolean;default:true;index:,where:deleted_at IS NULL" json:"is_active"`
	PortalPasswordHash *string        `gorm:"type:varchar(255)" json:"-"`
	PortalLastLogin    *time.Time     `gorm:"type:timestamptz" json:"portal_last_login"`
	Notes              *string        `gorm:"type:text" json:"notes"`
	Tags               *string        `gorm:"type:jsonb" json:"tags"`
	CreatedBy          *string        `gorm:"type:uuid" json:"created_by"`
	UpdatedBy          *string        `gorm:"type:uuid" json:"updated_by"`
	CreatedAt          time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Creator       *User                  `gorm:"foreignKey:CreatedBy" json:"-"`
	Updater       *User                  `gorm:"foreignKey:UpdatedBy" json:"-"`
	Subscriptions []Subscription         `gorm:"foreignKey:CustomerID" json:"subscriptions,omitempty"`
	Invoices      []Invoice              `gorm:"foreignKey:CustomerID" json:"invoices,omitempty"`
	Payments      []Payment              `gorm:"foreignKey:CustomerID" json:"payments,omitempty"`
	Registrations []CustomerRegistration `gorm:"foreignKey:CustomerID" json:"registrations,omitempty"`
}

func (Customer) TableName() string {
	return "customers"
}
