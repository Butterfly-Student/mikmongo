package mikhmon

// SalesReport represents a sales report entry stored in MikroTik /system/script
type SalesReport struct {
	ID       string `json:"id,omitempty"` // Script .id
	Name     string `json:"name"`         // Full script name with all data
	Date     string `json:"date"`         // Date of transaction
	Time     string `json:"time"`         // Time of transaction
	User     string `json:"user"`         // Username
	Price    int64  `json:"price"`        // Price
	IP       string `json:"ip"`           // IP address
	MAC      string `json:"mac"`          // MAC address
	Validity string `json:"validity"`     // Validity period
	Profile  string `json:"profile"`      // Profile name
	Comment  string `json:"comment"`      // Original comment
	Owner    string `json:"owner"`        // MonthYear format (e.g., "jan2024")
	Source   string `json:"source"`       // Date (same as Date field)
}

// ReportSummary represents summary statistics for reports
type ReportSummary struct {
	TotalCount   int   `json:"totalCount"`
	TotalSales   int64 `json:"totalSales"`
	TotalRevenue int64 `json:"totalRevenue"`
}

// SalesReportRequest represents a request to add a sales report
type SalesReportRequest struct {
	User     string `json:"user" validate:"required"`
	Price    int64  `json:"price"`
	IP       string `json:"ip,omitempty"`
	MAC      string `json:"mac,omitempty"`
	Validity string `json:"validity,omitempty"`
	Profile  string `json:"profile,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

// ReportFilter represents filter options for getting reports
type ReportFilter struct {
	Owner string `json:"owner,omitempty"` // Filter by month (e.g., "jan2024")
	Day   string `json:"day,omitempty"`   // Filter by day (e.g., "jan/01/2024")
	User  string `json:"user,omitempty"`  // Filter by username
	Limit int    `json:"limit,omitempty"` // Limit results
}
