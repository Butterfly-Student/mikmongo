package model

import (
	"time"
)

// MessageTemplate represents WhatsApp/Email notification templates per system event
type MessageTemplate struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Event     string    `gorm:"type:varchar(80);not null;index;uniqueIndex:idx_message_templates_event_channel"`
	Channel   string    `gorm:"type:varchar(20);not null;check:channel IN ('whatsapp', 'email');default:'whatsapp';index;uniqueIndex:idx_message_templates_event_channel"`
	Subject   *string   `gorm:"type:varchar(200)"`
	Body      string    `gorm:"type:text;not null"`
	IsActive  bool      `gorm:"type:boolean;default:true;index"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
}

func (MessageTemplate) TableName() string {
	return "message_templates"
}
