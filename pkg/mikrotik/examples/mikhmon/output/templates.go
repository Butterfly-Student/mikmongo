package output

import (
	"bytes"
	"fmt"
	"text/template"

	mikhmonDomain "github.com/Butterfly-Student/go-ros/domain/mikhmon"
)

// VoucherPrintTemplate is a template for printing vouchers
type VoucherPrintTemplate struct {
	Header string
	Row    string
	Footer string
}

// DefaultVoucherTemplate returns the default voucher print template
func DefaultVoucherTemplate() VoucherPrintTemplate {
	return VoucherPrintTemplate{
		Header: "================================\nVOUCHER BATCH: {{.Batch.Code}}\nProfile: {{.Batch.Profile}}\nQuantity: {{.Batch.Quantity}}\n================================\n",
		Row:    "{{.Number}}. {{.Voucher.Name}}\n",
		Footer: "================================\n",
	}
}

// VoucherTemplateData is the data for voucher template
type VoucherTemplateData struct {
	Batch   mikhmonDomain.VoucherBatch
	Voucher mikhmonDomain.Voucher
	Number  int
}

// RenderVoucherBatch renders a voucher batch using template
func RenderVoucherBatch(batch *mikhmonDomain.VoucherBatch, tmpl VoucherPrintTemplate) (string, error) {
	var buf bytes.Buffer

	// Render header
	headerTmpl, err := template.New("header").Parse(tmpl.Header)
	if err != nil {
		return "", err
	}
	if err := headerTmpl.Execute(&buf, struct{ Batch mikhmonDomain.VoucherBatch }{Batch: *batch}); err != nil {
		return "", err
	}

	// Render rows
	rowTmpl, err := template.New("row").Parse(tmpl.Row)
	if err != nil {
		return "", err
	}

	for i, v := range batch.Vouchers {
		data := VoucherTemplateData{
			Batch:   *batch,
			Voucher: v,
			Number:  i + 1,
		}
		if err := rowTmpl.Execute(&buf, data); err != nil {
			return "", err
		}
	}

	// Render footer
	buf.WriteString(tmpl.Footer)

	return buf.String(), nil
}

// SimpleVoucherOutput returns a simple text output for vouchers
func SimpleVoucherOutput(batch *mikhmonDomain.VoucherBatch) string {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("Batch: %s\n", batch.Code))
	buf.WriteString(fmt.Sprintf("Profile: %s | Qty: %d\n", batch.Profile, batch.Quantity))
	buf.WriteString("--------------------------------\n")

	for i, v := range batch.Vouchers {
		if v.Mode == "vc" {
			buf.WriteString(fmt.Sprintf("%d. %s (same password)\n", i+1, v.Name))
		} else {
			buf.WriteString(fmt.Sprintf("%d. %s / %s\n", i+1, v.Name, v.Password))
		}
	}

	return buf.String()
}
