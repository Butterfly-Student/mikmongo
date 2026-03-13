package report

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"mikmongo/pkg/mikrotik/client"
	"mikmongo/pkg/mikrotik/domain"
)

// Repository handles Report data access via RouterOS API (/system/script)
type Repository struct {
	client *client.Client
}

// NewRepository creates a new Report repository
func NewRepository(c *client.Client) *Repository {
	return &Repository{client: c}
}

// GetSalesReports retrieves sales reports by owner (month).
func (r *Repository) GetSalesReports(ctx context.Context, owner string) ([]*domain.SalesReport, error) {
	args := []string{"/system/script/print"}
	if owner != "" {
		args = append(args, "?owner="+owner)
	}
	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	reports := make([]*domain.SalesReport, 0, len(reply.Re))
	for _, re := range reply.Re {
		if rpt := parseSalesReport(re.Map); rpt != nil {
			reports = append(reports, rpt)
		}
	}
	return reports, nil
}

// GetSalesReportsByDay retrieves sales reports by day (source field).
func (r *Repository) GetSalesReportsByDay(ctx context.Context, day string) ([]*domain.SalesReport, error) {
	args := []string{"/system/script/print"}
	if day != "" {
		args = append(args, "?source="+day)
	}
	reply, err := r.client.RunArgsContext(ctx, args)
	if err != nil {
		return nil, err
	}
	reports := make([]*domain.SalesReport, 0, len(reply.Re))
	for _, re := range reply.Re {
		if rpt := parseSalesReport(re.Map); rpt != nil {
			reports = append(reports, rpt)
		}
	}
	return reports, nil
}

// AddSalesReport adds a sales report entry.
func (r *Repository) AddSalesReport(ctx context.Context, report *domain.SalesReport) error {
	name := fmt.Sprintf("%s-|-%s-|-%s-|-%.0f-|-%s-|-%s-|-%s-|-%s",
		report.Date,
		report.Time,
		report.Username,
		report.Price,
		report.IPAddress,
		report.MACAddress,
		report.Validity,
		report.Profile,
	)
	if report.VoucherComment != "" {
		name = name + "-|-" + report.VoucherComment
	}
	_, err := r.client.RunContext(ctx,
		"/system/script/add",
		"=name="+name,
		"=owner="+report.Owner,
		"=source="+report.Source,
		"=comment=mikhmon",
	)
	return err
}

// parseSalesReport parses a sales report from a script entry.
// Name format: $date-|-$time-|-$user-|-$price-|-$address-|-$mac-|-$validity-|-$profile-|-$comment
func parseSalesReport(data map[string]string) *domain.SalesReport {
	name := data["name"]
	if !strings.Contains(name, "-|-") {
		return nil
	}
	parts := strings.Split(name, "-|-")
	if len(parts) < 8 {
		return nil
	}
	price, _ := strconv.ParseFloat(parts[3], 64)
	return &domain.SalesReport{
		ID:             data[".id"],
		Name:           name,
		Owner:          data["owner"],
		Source:         data["source"],
		Comment:        data["comment"],
		DontReq:        data["dont-require-permissions"],
		RunCount:       data["run-count"],
		CopyOf:         data["copy-of"],
		Date:           parts[0],
		Time:           parts[1],
		Username:       parts[2],
		Price:          price,
		IPAddress:      parts[4],
		MACAddress:     parts[5],
		Validity:       parts[6],
		Profile:        parts[7],
		VoucherComment: getPart(parts, 8),
	}
}

func getPart(parts []string, index int) string {
	if index < len(parts) {
		return parts[index]
	}
	return ""
}
