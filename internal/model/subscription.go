package model

import (
	"time"

	"gorm.io/gorm"
)

// Subscription represents active customer services linking customers to plans and routers
type Subscription struct {
	ID              string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CustomerID      string         `gorm:"type:uuid;not null;index:,where:deleted_at IS NULL"`
	PlanID          string         `gorm:"type:uuid;not null;index:,where:deleted_at IS NULL"`
	RouterID        string         `gorm:"type:uuid;not null;index:,where:deleted_at IS NULL"`
	Username        string         `gorm:"type:varchar(100);unique;not null;index:,where:deleted_at IS NULL"`
	Password        string         `gorm:"type:varchar(255);not null"`
	StaticIP        *string        `gorm:"type:varchar(45)"`
	Gateway         *string        `gorm:"type:varchar(15)"`
	MACAddress      *string        `gorm:"type:varchar(17)"`
	Status          string         `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending', 'active', 'suspended', 'isolated', 'expired', 'terminated');index:,where:deleted_at IS NULL"`
	ActivatedAt     *time.Time     `gorm:"type:timestamptz"`
	ExpiryDate      *time.Time     `gorm:"type:date;index:,where:deleted_at IS NULL"`
	BillingDay      *int           `gorm:"type:integer" json:"billing_day,omitempty"`
	AutoIsolate     bool           `gorm:"type:boolean;default:true" json:"auto_isolate"`
	GracePeriodDays *int           `gorm:"type:integer" json:"grace_period_days,omitempty"`
	SuspendReason   *string        `gorm:"type:text"`
	TerminatedAt    *time.Time     `gorm:"type:timestamptz"`
	PreviousPlanID  *string        `gorm:"type:uuid"`
	Notes           *string        `gorm:"type:text"`
	CreatedBy       *string        `gorm:"type:uuid"`
	CreatedAt       time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Relationships
	Customer     Customer          `gorm:"foreignKey:CustomerID"`
	Plan         BandwidthProfile  `gorm:"foreignKey:PlanID"`
	Router       MikrotikRouter    `gorm:"foreignKey:RouterID"`
	PreviousPlan *BandwidthProfile `gorm:"foreignKey:PreviousPlanID"`
	Creator      *User             `gorm:"foreignKey:CreatedBy"`
	Invoices     []Invoice         `gorm:"foreignKey:SubscriptionID"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
