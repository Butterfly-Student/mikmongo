package model

import (
	"time"
)

// AuditLog represents audit trail of all important system changes
type AuditLog struct {
	ID          string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	AdminID     *string   `gorm:"type:uuid;index"`
	Action      string    `gorm:"type:varchar(100);not null;index"`
	EntityType  string    `gorm:"type:varchar(50);not null;index:idx_audit_logs_entity"`
	EntityID    string    `gorm:"type:uuid;not null;index:idx_audit_logs_entity"`
	OldValue    *string   `gorm:"type:jsonb"`
	NewValue    *string   `gorm:"type:jsonb"`
	IPAddress   *string   `gorm:"type:varchar(45)"`
	UserAgent   *string   `gorm:"type:varchar(255)"`
	Notes       *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP;index:,sort:desc"`

	// Relationships
	Admin       *User `gorm:"foreignKey:AdminID"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
