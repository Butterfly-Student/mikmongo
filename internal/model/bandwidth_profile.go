package model

import (
	"time"

	"gorm.io/gorm"
)

// BandwidthProfile represents ISP service packages/bandwidth profiles
type BandwidthProfile struct {
	ID                 string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	RouterID           string         `gorm:"type:uuid;not null;index" json:"router_id"`
	ProfileCode        string         `gorm:"type:varchar(50);not null;index" json:"profile_code"`
	Name               string         `gorm:"type:varchar(100);not null" json:"name"`
	Description        *string        `gorm:"type:text" json:"description"`
	DownloadSpeed      int64          `gorm:"type:bigint;not null" json:"download_speed"`
	UploadSpeed        int64          `gorm:"type:bigint;not null" json:"upload_speed"`
	PriceMonthly       float64        `gorm:"type:decimal(12,2);not null" json:"price_monthly"`
	TaxRate            float64        `gorm:"type:decimal(5,4);default:0.11" json:"tax_rate"`
	BillingCycle       string         `gorm:"type:varchar(20);default:'monthly';check:billing_cycle IN ('daily', 'weekly', 'monthly', 'yearly')" json:"billing_cycle"`
	BillingDay         *int           `gorm:"type:integer" json:"billing_day"`
	IsActive           bool           `gorm:"type:boolean;default:true;index:,where:deleted_at IS NULL" json:"is_active"`
	IsVisible          bool           `gorm:"type:boolean;default:true;index:,where:deleted_at IS NULL AND is_active = true" json:"is_visible"`
	SortOrder          int            `gorm:"type:integer;default:0" json:"sort_order"`
	GracePeriodDays    int            `gorm:"type:integer;not null;default:3" json:"grace_period_days"`
	IsolateProfileName *string        `gorm:"type:varchar(100)" json:"isolate_profile_name"`
	RateLimit          *string        `gorm:"type:varchar(100);column:rate_limit" json:"rate_limit"`
	CreatedAt          time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Router        MikrotikRouter         `gorm:"foreignKey:RouterID"`
	Subscriptions []Subscription         `gorm:"foreignKey:PlanID"`
	InvoiceItems  []InvoiceItem          `gorm:"foreignKey:ProfileID"`
	Registrations []CustomerRegistration `gorm:"foreignKey:BandwidthProfileID"`
}

func (BandwidthProfile) TableName() string {
	return "bandwidth_profiles"
}
