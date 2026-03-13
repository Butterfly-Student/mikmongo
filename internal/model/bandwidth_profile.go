package model

import (
	"time"

	"gorm.io/gorm"
)

// BandwidthProfile represents ISP service packages/bandwidth profiles
type BandwidthProfile struct {
	ID                 string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RouterID           string         `gorm:"type:uuid;not null;index"`
	ProfileCode        string         `gorm:"type:varchar(50);not null;index"`
	Name               string         `gorm:"type:varchar(100);not null"`
	Description        *string        `gorm:"type:text"`
	DownloadSpeed      int64          `gorm:"type:bigint;not null"`
	UploadSpeed        int64          `gorm:"type:bigint;not null"`
	PriceMonthly       float64        `gorm:"type:decimal(12,2);not null"`
	TaxRate            float64        `gorm:"type:decimal(5,4);default:0.11"`
	BillingCycle       string         `gorm:"type:varchar(20);default:'monthly';check:billing_cycle IN ('daily', 'weekly', 'monthly', 'yearly')"`
	BillingDay         *int           `gorm:"type:integer"`
	IsActive           bool           `gorm:"type:boolean;default:true;index:,where:deleted_at IS NULL"`
	IsVisible          bool           `gorm:"type:boolean;default:true;index:,where:deleted_at IS NULL AND is_active = true"`
	SortOrder          int            `gorm:"type:integer;default:0"`
	GracePeriodDays    int            `gorm:"type:integer;not null;default:3"`
	IsolateProfileName *string        `gorm:"type:varchar(100)"`
	MikrotikConfig     []byte         `gorm:"type:jsonb"`
	CreatedAt          time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt          time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`

	// Relationships
	Router        MikrotikRouter         `gorm:"foreignKey:RouterID"`
	Subscriptions []Subscription         `gorm:"foreignKey:PlanID"`
	InvoiceItems  []InvoiceItem          `gorm:"foreignKey:ProfileID"`
	Registrations []CustomerRegistration `gorm:"foreignKey:BandwidthProfileID"`
}

func (BandwidthProfile) TableName() string {
	return "bandwidth_profiles"
}
