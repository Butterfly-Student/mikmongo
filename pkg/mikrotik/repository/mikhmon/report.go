package mikhmon

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Butterfly-Student/go-ros/client"
	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
	"github.com/Butterfly-Student/go-ros/repository/system"
)

// reportRepository implements ReportRepository interface
type reportRepository struct {
	client     *client.Client
	systemRepo system.Repository
}

// NewReportRepository creates a new report repository
func NewReportRepository(c *client.Client, sr system.Repository) ReportRepository {
	return &reportRepository{
		client:     c,
		systemRepo: sr,
	}
}

// AddReport adds a sales report to MikroTik /system/script
func (r *reportRepository) AddReport(ctx context.Context, req *mikhmonDomain.SalesReportRequest) error {
	// Get current date and time
	now := time.Now()
	date := now.Format("jan/02/2006")
	timeStr := now.Format("15:04:05")
	month := now.Format("jan2006")

	// Format script name: date-|-time-|-user-|-price-|-ip-|-mac-|-validity-|-profile-|-comment
	scriptName := fmt.Sprintf("%s-|-%s-|-%s-|-%d-|-%s-|-%s-|-%s-|-%s-|-%s",
		date, timeStr, req.User, req.Price, req.IP, req.MAC, req.Validity, req.Profile, req.Comment)

	// Add script to MikroTik
	_, err := r.client.RunContext(ctx,
		"/system/script/add",
		"=name="+scriptName,
		"=owner="+month,
		"=source="+date,
		"=comment=mikhmon",
	)

	if err != nil {
		return fmt.Errorf("failed to add sales report: %w", err)
	}

	return nil
}

// GetReportsByOwner retrieves reports by owner (month format: jan2024)
func (r *reportRepository) GetReportsByOwner(ctx context.Context, owner string) ([]*mikhmonDomain.SalesReport, error) {
	reply, err := r.client.RunContext(ctx, "/system/script/print", "?owner="+owner)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports by owner: %w", err)
	}

	reports := make([]*mikhmonDomain.SalesReport, 0, len(reply.Re))
	for _, re := range reply.Re {
		report := r.parseSalesReport(re.Map)
		if report != nil {
			reports = append(reports, report)
		}
	}

	return reports, nil
}

// GetReportsByDay retrieves reports by day (format: jan/01/2024)
func (r *reportRepository) GetReportsByDay(ctx context.Context, day string) ([]*mikhmonDomain.SalesReport, error) {
	reply, err := r.client.RunContext(ctx, "/system/script/print", "?source="+day)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports by day: %w", err)
	}

	reports := make([]*mikhmonDomain.SalesReport, 0, len(reply.Re))
	for _, re := range reply.Re {
		report := r.parseSalesReport(re.Map)
		if report != nil {
			reports = append(reports, report)
		}
	}

	return reports, nil
}

// GetReportSummary gets summary statistics for reports
func (r *reportRepository) GetReportSummary(ctx context.Context, filter *mikhmonDomain.ReportFilter) (*mikhmonDomain.ReportSummary, error) {
	var reports []*mikhmonDomain.SalesReport
	var err error

	if filter.Owner != "" {
		reports, err = r.GetReportsByOwner(ctx, filter.Owner)
	} else if filter.Day != "" {
		reports, err = r.GetReportsByDay(ctx, filter.Day)
	} else {
		// Get all reports with mikhmon comment
		reply, err := r.client.RunContext(ctx, "/system/script/print", "?comment=mikhmon")
		if err != nil {
			return nil, fmt.Errorf("failed to get reports: %w", err)
		}

		reports = make([]*mikhmonDomain.SalesReport, 0, len(reply.Re))
		for _, re := range reply.Re {
			report := r.parseSalesReport(re.Map)
			if report != nil {
				reports = append(reports, report)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	// Apply limit
	if filter.Limit > 0 && len(reports) > filter.Limit {
		reports = reports[:filter.Limit]
	}

	// Calculate summary
	summary := &mikhmonDomain.ReportSummary{
		TotalCount: len(reports),
	}

	for _, report := range reports {
		summary.TotalSales++
		summary.TotalRevenue += report.Price
	}

	return summary, nil
}

// parseSalesReport parses a script entry into SalesReport
func (r *reportRepository) parseSalesReport(m map[string]string) *mikhmonDomain.SalesReport {
	name := m["name"]
	if name == "" {
		return nil
	}

	// Parse script name format: date-|-time-|-user-|-price-|-ip-|-mac-|-validity-|-profile-|-comment
	parts := strings.Split(name, "-|-")
	if len(parts) != 9 {
		return nil
	}

	price, _ := strconv.ParseInt(parts[3], 10, 64)

	return &mikhmonDomain.SalesReport{
		ID:       m[".id"],
		Name:     name,
		Date:     parts[0],
		Time:     parts[1],
		User:     parts[2],
		Price:    price,
		IP:       parts[4],
		MAC:      parts[5],
		Validity: parts[6],
		Profile:  parts[7],
		Comment:  parts[8],
		Owner:    m["owner"],
		Source:   m["source"],
	}
}
