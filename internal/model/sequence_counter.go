package model

import (
	"time"
)

// SequenceCounter represents automatic numbering counters, thread-safe with SELECT FOR UPDATE
type SequenceCounter struct {
	ID           string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name         string    `gorm:"type:varchar(50);unique;not null;index"`
	Prefix       *string   `gorm:"type:varchar(10)"`
	Padding      int       `gorm:"type:integer;default:5"`
	LastNumber   int       `gorm:"type:integer;default:0"`
	ResetMonthly bool      `gorm:"type:boolean;default:false"`
	ResetYearly  bool      `gorm:"type:boolean;default:false"`
	LastReset    *time.Time `gorm:"type:date"`
	CreatedAt    time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
}

func (SequenceCounter) TableName() string {
	return "sequence_counters"
}
