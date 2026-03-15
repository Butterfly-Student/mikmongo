package mikhmon

// ProfileConfig represents Mikhmon profile configuration with on-login script
type ProfileConfig struct {
	Name          string `json:"name" validate:"required"`
	AddressPool   string `json:"addressPool,omitempty"`
	RateLimit     string `json:"rateLimit,omitempty"` // e.g., "1M/2M"
	SharedUsers   int    `json:"sharedUsers,omitempty"`
	ParentQueue   string `json:"parentQueue,omitempty"`
	Price         int64  `json:"price,omitempty"`         // Price in currency
	SellingPrice  int64  `json:"sellingPrice,omitempty"`  // Selling price
	Validity      string `json:"validity,omitempty"`      // e.g., "30d", "1d"
	ExpireMode    string `json:"expireMode,omitempty"`    // "rem", "ntf", "remc", "ntfc", "0"
	LockUser      bool   `json:"lockUser,omitempty"`      // MAC address locking
	LockServer    bool   `json:"lockServer,omitempty"`    // Server locking
	OnLoginScript string `json:"onLoginScript,omitempty"` // Generated on-login script
}

// ExpireMode constants
const (
	ExpireModeRemove       = "rem"  // Remove user on expire
	ExpireModeNotify       = "ntf"  // Disable user on expire (limit-uptime=1s)
	ExpireModeRemoveRecord = "remc" // Remove + record to report
	ExpireModeNotifyRecord = "ntfc" // Disable + record to report
	ExpireModeNoExpire     = "0"    // No expiration
)

// OnLoginScriptData represents data for generating on-login script
type OnLoginScriptData struct {
	Mode         string `json:"mode"` // rem, ntf, remc, ntfc, 0
	Price        int64  `json:"price"`
	Validity     string `json:"validity"` // e.g., "30d"
	SellingPrice int64  `json:"sellingPrice"`
	NoExp        bool   `json:"noExp"`       // No expiration flag
	LockUser     string `json:"lockUser"`    // "Enable" or "Disable"
	LockServer   string `json:"lockServer"`  // "Enable" or "Disable"
	ProfileName  string `json:"profileName"` // Profile name for recording script
}

// ProfileRequest represents a request to add/update profile with Mikhmon config
type ProfileRequest struct {
	Name        string        `json:"name" validate:"required"`
	AddressPool string        `json:"addressPool,omitempty"`
	RateLimit   string        `json:"rateLimit,omitempty"`
	SharedUsers int           `json:"sharedUsers,omitempty"`
	ParentQueue string        `json:"parentQueue,omitempty"`
	Config      ProfileConfig `json:"config"`
}
