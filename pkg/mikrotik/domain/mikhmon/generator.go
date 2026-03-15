package mikhmon

// CharSet constants for voucher generation
const (
	CharSetLower   = "lower"   // a-z
	CharSetUpper   = "upper"   // A-Z
	CharSetUpplow  = "upplow"  // a-zA-Z
	CharSetLower1  = "lower1"  // a-z0-9
	CharSetUpper1  = "upper1"  // A-Z0-9
	CharSetUpplow1 = "upplow1" // a-zA-Z0-9
	CharSetMix     = "mix"     // a-zA-Z0-9 (alias)
	CharSetMix1    = "mix1"    // a-zA-Z0-9 with special
	CharSetMix2    = "mix2"    // Extended mix
	CharSetNumeric = "num"     // 0-9
)

// GeneratorConfig represents configuration for user/voucher generator
type GeneratorConfig struct {
	Length  int    `json:"length" validate:"min=3,max=12"`
	Prefix  string `json:"prefix,omitempty"`
	CharSet string `json:"charSet" validate:"required"`
}

// GeneratorResult represents result of generation
type GeneratorResult struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

// VoucherMode constants
const (
	VoucherModeVoucher      = "vc" // Voucher (username = password)
	VoucherModeUserPassword = "up" // User/Password (username != password)
)

// TimeLimitFormat represents time limit format helpers
// Format: 1d2h3m4s (days, hours, minutes, seconds)
type TimeLimitFormat struct {
	Days    int `json:"days,omitempty"`
	Hours   int `json:"hours,omitempty"`
	Minutes int `json:"minutes,omitempty"`
	Seconds int `json:"seconds,omitempty"`
}

// DataLimitFormat represents data limit format helpers
// Format: 1G, 500M, 1024K
type DataLimitFormat struct {
	Value int64  `json:"value"`
	Unit  string `json:"unit"` // G, M, K
}
