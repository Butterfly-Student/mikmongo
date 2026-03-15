package model

import (
	"time"

	"gorm.io/gorm"
)

// MikrotikRouter represents managed MikroTik routers
type MikrotikRouter struct {
	ID                string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name              string         `gorm:"type:varchar(100);not null;index"`
	Address           string         `gorm:"type:varchar(100);not null;index"`
	Area              *string        `gorm:"type:varchar(100);index:,where:area IS NOT NULL AND deleted_at IS NULL"`
	APIPort           int            `gorm:"type:integer;default:8728"`
	RESTPort          int            `gorm:"type:integer;default:80"`
	Username          string         `gorm:"type:varchar(100);not null"`
	PasswordEncrypted string         `gorm:"type:text;not null"`
	UseSSL            bool           `gorm:"type:boolean;default:false"`
	IsMaster          bool           `gorm:"type:boolean;default:false"`
	IsActive          bool           `gorm:"type:boolean;default:true;index:,where:deleted_at IS NULL"`
	Status            string         `gorm:"type:varchar(20);default:'unknown';check:status IN ('online', 'offline', 'unknown');index:,where:deleted_at IS NULL"`
	LastSeenAt        *time.Time     `gorm:"type:timestamptz"`
	Notes             *string        `gorm:"type:text"`
	CreatedAt         time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt         gorm.DeletedAt `gorm:"index"`

	// Relationships
	Subscriptions []Subscription `gorm:"foreignKey:RouterID"`
}

func (MikrotikRouter) TableName() string {
	return "mikrotik_routers"
}
