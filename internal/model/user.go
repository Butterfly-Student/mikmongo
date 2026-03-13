package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents admin users and operators in the ISP system
type User struct {
	ID           string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FullName     string         `gorm:"type:varchar(100);not null"`
	Email        string         `gorm:"type:varchar(100);unique;not null"`
	Phone        string         `gorm:"type:varchar(20)"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	Role         string         `gorm:"type:varchar(20);not null;default:'cs';check:role IN ('superadmin', 'admin', 'cs', 'billing', 'technician', 'readonly')"`
	IsActive     bool           `gorm:"type:boolean;default:true"`
	LastLogin    *time.Time     `gorm:"type:timestamptz"`
	LastIP       string         `gorm:"type:varchar(45)"`
	BearerKey    string         `gorm:"type:varchar(255);unique"`
	CreatedAt    time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relationships
	CreatedCustomers []Customer `gorm:"foreignKey:CreatedBy"`
	UpdatedCustomers []Customer `gorm:"foreignKey:UpdatedBy"`
}

func (User) TableName() string {
	return "users"
}
