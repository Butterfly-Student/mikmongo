package model

import (
	"time"

	"gorm.io/gorm"
)

// MikrotikRouter represents managed MikroTik routers
type MikrotikRouter struct {
	ID                string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name              string         `gorm:"type:varchar(100);not null;index" json:"name"`
	Address           string         `gorm:"type:varchar(100);not null;index" json:"address"`
	Area              *string        `gorm:"type:varchar(100);index:,where:area IS NOT NULL AND deleted_at IS NULL" json:"area"`
	APIPort           int            `gorm:"type:integer;default:8728" json:"api_port"`
	RESTPort          int            `gorm:"type:integer;default:80" json:"rest_port"`
	Username          string         `gorm:"type:varchar(100);not null" json:"username"`
	PasswordEncrypted string         `gorm:"type:text;not null" json:"-"`
	UseSSL            bool           `gorm:"type:boolean;default:false" json:"use_ssl"`
	IsMaster          bool           `gorm:"type:boolean;default:false" json:"is_master"`
	IsActive          bool           `gorm:"type:boolean;default:true;index:,where:deleted_at IS NULL" json:"is_active"`
	Status            string         `gorm:"type:varchar(20);default:'unknown';check:status IN ('online', 'offline', 'unknown');index:,where:deleted_at IS NULL" json:"status"`
	LastSeenAt        *time.Time     `gorm:"type:timestamptz" json:"last_seen_at"`
	Notes             *string        `gorm:"type:text" json:"notes"`
	CreatedAt         time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt         time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Subscriptions []Subscription `gorm:"foreignKey:RouterID"`
}

func (MikrotikRouter) TableName() string {
	return "mikrotik_routers"
}
