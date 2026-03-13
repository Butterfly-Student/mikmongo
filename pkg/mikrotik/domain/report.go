package domain

// SalesReport represents a sales report entry from MikroTik /system/script
// Format name: $date-|-$time-|-$user-|-$price-|-$address-|-$mac-|-$validity-|-$profile-|-$comment
type SalesReport struct {
	ID             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	Owner          string  `json:"owner,omitempty"`
	Source         string  `json:"source,omitempty"`
	Comment        string  `json:"comment,omitempty"`
	DontReq        string  `json:"dont-require-permissions,omitempty"`
	RunCount       string  `json:"run-count,omitempty"`
	CopyOf         string  `json:"copy-of,omitempty"`
	// Parsed fields from name
	Date           string  `json:"date"`
	Time           string  `json:"time"`
	Username       string  `json:"username"`
	Price          float64 `json:"price"`
	IPAddress      string  `json:"ipAddress"`
	MACAddress     string  `json:"macAddress"`
	Validity       string  `json:"validity"`
	Profile        string  `json:"profile"`
	VoucherComment string  `json:"voucherComment,omitempty"`
}

// ReportFilter represents filter parameters for reports
type ReportFilter struct {
	RouterID string `json:"routerId" validate:"required"`
	Day      string `json:"day,omitempty"`
	Month    string `json:"month,omitempty"`
	Year     string `json:"year,omitempty"`
	Profile  string `json:"profile,omitempty"`
}

// ReportSummary represents a summary of sales
type ReportSummary struct {
	TotalVouchers int                       `json:"totalVouchers"`
	TotalAmount   float64                   `json:"totalAmount"`
	ByProfile     map[string]ProfileSummary `json:"byProfile,omitempty"`
}

// ProfileSummary represents summary by profile
type ProfileSummary struct {
	Count int     `json:"count"`
	Total float64 `json:"total"`
}

// ReportResponse represents report API response
type ReportResponse struct {
	Data    []*SalesReport `json:"data"`
	Summary *ReportSummary `json:"summary,omitempty"`
	Filter  *ReportFilter  `json:"filter,omitempty"`
}

// LiveReportRequest represents a live report request
type LiveReportRequest struct {
	RouterID string `json:"routerId" validate:"required"`
	Month    string `json:"month,omitempty"`
	Day      string `json:"day,omitempty"`
}

// VoucherSaleRecord adalah data penjualan voucher untuk disimpan di database.
// Berbeda dari SalesReport (yang mikhmon-style untuk MikroTik), struct ini
// digunakan untuk persistensi DB.
type VoucherSaleRecord struct {
	RouterID   string
	SoldBy     string  // user ID admin yang menjual (opsional)
	SaleDate   string  // "2006-01-02"
	SaleTime   string  // "15:04:05"
	Username   string
	Price      float64
	IPAddress  string
	MACAddress string
	Validity   string
	Profile    string
	Comment    string
}
