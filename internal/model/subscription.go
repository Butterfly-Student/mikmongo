package model

import (
	"time"

	"gorm.io/gorm"
)

// Subscription represents active customer services linking customers to plans and routers
type Subscription struct {
	ID              string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CustomerID      string         `gorm:"type:uuid;not null;index:,where:deleted_at IS NULL" json:"customer_id"`
	PlanID          string         `gorm:"type:uuid;not null;index:,where:deleted_at IS NULL" json:"plan_id"`
	RouterID        string         `gorm:"type:uuid;not null;index:,where:deleted_at IS NULL" json:"router_id"`
	Username        string         `gorm:"type:varchar(100);unique;not null;index:,where:deleted_at IS NULL" json:"username"`
	Password        string         `gorm:"type:varchar(255);not null" json:"-"`
	StaticIP        *string        `gorm:"type:varchar(45)" json:"static_ip"`
	Gateway         *string        `gorm:"type:varchar(15)" json:"gateway"`
	MtPPPID         *string        `gorm:"type:varchar(50);column:mt_ppp_id" json:"mt_ppp_id"`
	Status          string         `gorm:"type:varchar(20);not null;default:'pending';check:status IN ('pending', 'active', 'suspended', 'isolated', 'expired', 'terminated');index:,where:deleted_at IS NULL" json:"status"`
	ActivatedAt     *time.Time     `gorm:"type:timestamptz" json:"activated_at"`
	ExpiryDate      *time.Time     `gorm:"type:date;index:,where:deleted_at IS NULL" json:"expiry_date"`
	BillingDay      *int           `gorm:"type:integer" json:"billing_day,omitempty"`
	AutoIsolate     bool           `gorm:"type:boolean;default:true" json:"auto_isolate"`
	GracePeriodDays *int           `gorm:"type:integer" json:"grace_period_days,omitempty"`
	SuspendReason   *string        `gorm:"type:text" json:"suspend_reason"`
	TerminatedAt    *time.Time     `gorm:"type:timestamptz" json:"terminated_at"`
	PreviousPlanID  *string        `gorm:"type:uuid" json:"previous_plan_id"`
	Notes           *string        `gorm:"type:text" json:"notes"`
	CreatedBy       *string        `gorm:"type:uuid" json:"created_by"`
	CreatedAt       time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Customer     Customer          `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Plan         BandwidthProfile  `gorm:"foreignKey:PlanID" json:"plan,omitempty"`
	Router       MikrotikRouter    `gorm:"foreignKey:RouterID" json:"router,omitempty"`
	PreviousPlan *BandwidthProfile `gorm:"foreignKey:PreviousPlanID" json:"previous_plan,omitempty"`
	Creator      *User             `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Invoices     []Invoice         `gorm:"foreignKey:SubscriptionID" json:"invoices,omitempty"`
}

func (Subscription) TableName() string {
	return "subscriptions"
}
