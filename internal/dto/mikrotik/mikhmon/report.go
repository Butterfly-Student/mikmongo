package mikhmon

import (
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

type CreateReportRequest struct {
	User     string `json:"user" binding:"required"`
	Price    int64  `json:"price"`
	IP       string `json:"ip,omitempty"`
	MAC      string `json:"mac,omitempty"`
	Validity string `json:"validity,omitempty"`
	Profile  string `json:"profile,omitempty"`
	Comment  string `json:"comment,omitempty"`
}

type ReportFilterQuery struct {
	Owner string `form:"owner"`
	Day   string `form:"day"`
	User  string `form:"user"`
	Limit int    `form:"limit"`
}

type SalesReportResponse struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	Date     string `json:"date"`
	Time     string `json:"time"`
	User     string `json:"user"`
	Price    int64  `json:"price"`
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	Validity string `json:"validity"`
	Profile  string `json:"profile"`
	Comment  string `json:"comment"`
	Owner    string `json:"owner"`
	Source   string `json:"source"`
}

type ReportSummaryResponse struct {
	TotalCount   int   `json:"total_count"`
	TotalSales   int64 `json:"total_sales"`
	TotalRevenue int64 `json:"total_revenue"`
}

func (r *CreateReportRequest) ToDomain() *mikhmonDomain.SalesReportRequest {
	return &mikhmonDomain.SalesReportRequest{
		User:     r.User,
		Price:    r.Price,
		IP:       r.IP,
		MAC:      r.MAC,
		Validity: r.Validity,
		Profile:  r.Profile,
		Comment:  r.Comment,
	}
}

func (q *ReportFilterQuery) ToDomain() *mikhmonDomain.ReportFilter {
	return &mikhmonDomain.ReportFilter{
		Owner: q.Owner,
		Day:   q.Day,
		User:  q.User,
		Limit: q.Limit,
	}
}

func SalesReportToResponse(report *mikhmonDomain.SalesReport) SalesReportResponse {
	return SalesReportResponse{
		ID:       report.ID,
		Name:     report.Name,
		Date:     report.Date,
		Time:     report.Time,
		User:     report.User,
		Price:    report.Price,
		IP:       report.IP,
		MAC:      report.MAC,
		Validity: report.Validity,
		Profile:  report.Profile,
		Comment:  report.Comment,
		Owner:    report.Owner,
		Source:   report.Source,
	}
}

func SalesReportsToResponse(reports []*mikhmonDomain.SalesReport) []SalesReportResponse {
	result := make([]SalesReportResponse, len(reports))
	for i, r := range reports {
		result[i] = SalesReportToResponse(r)
	}
	return result
}

func ReportSummaryToResponse(s *mikhmonDomain.ReportSummary) ReportSummaryResponse {
	return ReportSummaryResponse{
		TotalCount:   s.TotalCount,
		TotalSales:   s.TotalSales,
		TotalRevenue: s.TotalRevenue,
	}
}
