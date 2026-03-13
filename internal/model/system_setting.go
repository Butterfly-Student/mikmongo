package model

import (
	"time"
)

// SystemSetting represents application configuration in key-value format per group
type SystemSetting struct {
	ID          string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	GroupName   string     `gorm:"type:varchar(50);not null;index;uniqueIndex:idx_system_settings_group_key"`
	KeyName     string     `gorm:"type:varchar(100);not null;uniqueIndex:idx_system_settings_group_key"`
	Value       *string    `gorm:"type:text"`
	Type        string     `gorm:"type:varchar(20);check:type IN ('string', 'integer', 'boolean', 'json', 'password');default:'string'"`
	Label       *string    `gorm:"type:varchar(150)"`
	Description *string    `gorm:"type:text"`
	IsEncrypted bool       `gorm:"type:boolean;default:false"`
	IsPublic    bool       `gorm:"type:boolean;default:false;index:,where:is_public = true"`
	UpdatedAt   *time.Time `gorm:"type:timestamptz"`
	UpdatedBy   *string    `gorm:"type:uuid"`

	// Relationships
	Updater     *User `gorm:"foreignKey:UpdatedBy"`
}

func (SystemSetting) TableName() string {
	return "system_settings"
}
