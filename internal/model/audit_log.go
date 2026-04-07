package model

import (
	"time"
)

// AuditLog represents audit trail of all important system changes
type AuditLog struct {
	ID         string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AdminID    *string   `gorm:"type:uuid;index" json:"admin_id"`
	Action     string    `gorm:"type:varchar(100);not null;index" json:"action"`
	EntityType string    `gorm:"type:varchar(50);not null;index:idx_audit_logs_entity" json:"entity_type"`
	EntityID   string    `gorm:"type:uuid;not null;index:idx_audit_logs_entity" json:"entity_id"`
	OldValue   *string   `gorm:"type:jsonb" json:"old_value"`
	NewValue   *string   `gorm:"type:jsonb" json:"new_value"`
	IPAddress  *string   `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent  *string   `gorm:"type:varchar(255)" json:"user_agent"`
	Notes      *string   `gorm:"type:text" json:"notes"`
	CreatedAt  time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP;index:,sort:desc" json:"created_at"`

	// Relationships
	Admin *User `gorm:"foreignKey:AdminID"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
