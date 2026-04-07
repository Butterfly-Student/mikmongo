package model

import (
	"time"
)

// SequenceCounter represents automatic numbering counters, thread-safe with SELECT FOR UPDATE
type SequenceCounter struct {
	ID           string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name         string     `gorm:"type:varchar(50);unique;not null;index" json:"name"`
	Prefix       *string    `gorm:"type:varchar(10)" json:"prefix"`
	Padding      int        `gorm:"type:integer;default:5" json:"padding"`
	LastNumber   int        `gorm:"type:integer;default:0" json:"last_number"`
	ResetMonthly bool       `gorm:"type:boolean;default:false" json:"reset_monthly"`
	ResetYearly  bool       `gorm:"type:boolean;default:false" json:"reset_yearly"`
	LastReset    *time.Time `gorm:"type:date" json:"last_reset"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (SequenceCounter) TableName() string {
	return "sequence_counters"
}
