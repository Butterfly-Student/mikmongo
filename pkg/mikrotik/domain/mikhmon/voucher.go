package mikhmon

// Voucher represents a generated hotspot voucher
type Voucher struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Password string `json:"password,omitempty"`
	Profile  string `json:"profile,omitempty"`
	Server   string `json:"server,omitempty"`
	Comment  string `json:"comment,omitempty"`
	Code     string `json:"code,omitempty"` // Voucher code (e.g., "123")
	Mode     string `json:"mode,omitempty"` // "vc" (voucher) or "up" (user/password)
	Date     string `json:"date,omitempty"` // Generation date
}

// VoucherBatch represents a batch of generated vouchers
type VoucherBatch struct {
	Code      string    `json:"code"` // Unique batch code
	Quantity  int       `json:"quantity"`
	Profile   string    `json:"profile"`
	Server    string    `json:"server"`
	TimeLimit string    `json:"timeLimit,omitempty"`
	DataLimit string    `json:"dataLimit,omitempty"`
	Vouchers  []Voucher `json:"vouchers"`
}

// VoucherGenerateRequest represents a request to generate vouchers
type VoucherGenerateRequest struct {
	Quantity   int    `json:"quantity" validate:"required,min=1,max=1000"`
	Server     string `json:"server,omitempty"`
	Profile    string `json:"profile" validate:"required"`
	Mode       string `json:"mode" validate:"required,oneof=vc up"` // vc=voucher, up=user/password
	NameLength int    `json:"nameLength" validate:"min=3,max=12"`
	Prefix     string `json:"prefix,omitempty"`
	CharSet    string `json:"charSet" validate:"required"` // lower, upper, upplow, lower1, upper1, upplow1, mix, mix1, mix2, num
	TimeLimit  string `json:"timeLimit,omitempty"`         // e.g., "1d", "2h", "30d"
	DataLimit  string `json:"dataLimit,omitempty"`         // e.g., "1G", "500M"
	Comment    string `json:"comment,omitempty"`
}

// VoucherTemplate represents a voucher print template
type VoucherTemplate struct {
	Name   string `json:"name"`
	Header string `json:"header"`
	Row    string `json:"row"`
	Footer string `json:"footer"`
}

// VoucherPrintData represents data for printing vouchers
type VoucherPrintData struct {
	Batch    VoucherBatch    `json:"batch"`
	Profile  string          `json:"profile"`
	Template VoucherTemplate `json:"template"`
}
