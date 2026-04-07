package model

import (
	"time"
)

// SystemSetting represents application configuration in key-value format per group
type SystemSetting struct {
	ID          string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GroupName   string     `gorm:"type:varchar(50);not null;index;uniqueIndex:idx_system_settings_group_key" json:"group_name"`
	KeyName     string     `gorm:"type:varchar(100);not null;uniqueIndex:idx_system_settings_group_key" json:"key_name"`
	Value       *string    `gorm:"type:text" json:"value"`
	Type        string     `gorm:"type:varchar(20);check:type IN ('string', 'integer', 'boolean', 'json', 'password');default:'string'" json:"type"`
	Label       *string    `gorm:"type:varchar(150)" json:"label"`
	Description *string    `gorm:"type:text" json:"description"`
	IsEncrypted bool       `gorm:"type:boolean;default:false" json:"is_encrypted"`
	IsPublic    bool       `gorm:"type:boolean;default:false;index:,where:is_public = true" json:"is_public"`
	UpdatedAt   *time.Time `gorm:"type:timestamptz" json:"updated_at"`
	UpdatedBy   *string    `gorm:"type:uuid" json:"updated_by"`

	// Relationships
	Updater *User `gorm:"foreignKey:UpdatedBy"`
}

func (SystemSetting) TableName() string {
	return "system_settings"
}
