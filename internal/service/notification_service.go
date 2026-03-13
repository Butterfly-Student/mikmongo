package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"text/template"

	"mikmongo/internal/model"
	"mikmongo/internal/notification"
	"mikmongo/internal/repository"
)

// NotificationService handles sending notifications via WhatsApp and Email
type NotificationService struct {
	templateRepo repository.MessageTemplateRepository
	settingRepo  repository.SystemSettingRepository
	gowaClient   *notification.GoWAClient
	emailClient  *notification.EmailClient
}

// NewNotificationService creates a new notification service
func NewNotificationService(
	templateRepo repository.MessageTemplateRepository,
	settingRepo repository.SystemSettingRepository,
) *NotificationService {
	return &NotificationService{
		templateRepo: templateRepo,
		settingRepo:  settingRepo,
	}
}

// initClients lazily initializes GoWA and Email clients from system_settings
func (s *NotificationService) initClients(ctx context.Context) {
	if s.gowaClient == nil {
		gowaURL := s.getSetting(ctx, "notification", "gowa_url")
		gowaSender := s.getSetting(ctx, "notification", "gowa_sender")
		gowaKey := s.getSetting(ctx, "notification", "gowa_auth_key")
		if gowaURL != "" {
			s.gowaClient = notification.NewGoWAClient(gowaURL, gowaSender, gowaKey)
		}
	}
	if s.emailClient == nil {
		smtpHost := s.getSetting(ctx, "notification", "smtp_host")
		smtpPort := s.getSetting(ctx, "notification", "smtp_port")
		smtpUser := s.getSetting(ctx, "notification", "smtp_user")
		smtpPass := s.getSetting(ctx, "notification", "smtp_password")
		smtpFrom := s.getSetting(ctx, "notification", "smtp_from")
		if smtpHost != "" {
			if smtpPort == "" {
				smtpPort = "587"
			}
			s.emailClient = notification.NewEmailClient(smtpHost, smtpPort, smtpUser, smtpPass, smtpFrom)
		}
	}
}

func (s *NotificationService) getSetting(ctx context.Context, group, key string) string {
	setting, err := s.settingRepo.GetByGroupAndKey(ctx, group, key)
	if err != nil || setting.Value == nil {
		return ""
	}
	return *setting.Value
}

// RenderTemplate renders a message template with data using {{key}} substitution
func (s *NotificationService) RenderTemplate(tmplBody string, data map[string]string) (string, error) {
	result := tmplBody
	for k, v := range data {
		result = strings.ReplaceAll(result, "{{"+k+"}}", v)
	}
	return result, nil
}

// renderGoTemplate renders using text/template with {{.key}} notation (Go style)
func renderGoTemplate(tmplBody string, data map[string]string) (string, error) {
	t, err := template.New("msg").Parse(tmplBody)
	if err != nil {
		return tmplBody, nil // fallback to raw
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return tmplBody, nil
	}
	return buf.String(), nil
}

// SendViaWhatsApp sends a WhatsApp message
func (s *NotificationService) SendViaWhatsApp(ctx context.Context, phone, message string) error {
	s.initClients(ctx)
	if s.gowaClient == nil {
		return fmt.Errorf("WhatsApp client not configured")
	}
	return s.gowaClient.SendMessage(ctx, phone, message)
}

// SendViaEmail sends an email
func (s *NotificationService) SendViaEmail(ctx context.Context, to, subject, body string) error {
	s.initClients(ctx)
	if s.emailClient == nil {
		return fmt.Errorf("Email client not configured")
	}
	return s.emailClient.SendEmail(ctx, to, subject, body)
}

// RenderAndSend renders template and sends via specified channel
func (s *NotificationService) RenderAndSend(ctx context.Context, event, channel, to string, data map[string]string) error {
	tmpl, err := s.templateRepo.GetByEventAndChannel(ctx, event, channel)
	if err != nil || !tmpl.IsActive {
		return nil // skip if no template
	}

	body, _ := s.RenderTemplate(tmpl.Body, data)

	switch channel {
	case "whatsapp":
		return s.SendViaWhatsApp(ctx, to, body)
	case "email":
		subject := ""
		if tmpl.Subject != nil {
			subj, _ := s.RenderTemplate(*tmpl.Subject, data)
			subject = subj
		}
		return s.SendViaEmail(ctx, to, subject, body)
	}
	return nil
}

// SendInvoiceCreated sends invoice created notification
func (s *NotificationService) SendInvoiceCreated(ctx context.Context, invoice *model.Invoice, customer *model.Customer) error {
	data := map[string]string{
		"name":       customer.FullName,
		"invoice_no": invoice.InvoiceNumber,
		"amount":     fmt.Sprintf("%.0f", invoice.TotalAmount),
		"due_date":   invoice.DueDate.Format("02-01-2006"),
	}
	phone := customer.Phone
	var errs []string
	if err := s.RenderAndSend(ctx, "invoice_created", "whatsapp", phone, data); err != nil {
		errs = append(errs, err.Error())
	}
	if customer.Email != nil {
		if err := s.RenderAndSend(ctx, "invoice_created", "email", *customer.Email, data); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("notification errors: %s", strings.Join(errs, "; "))
	}
	return nil
}

// SendPaymentConfirmed sends payment confirmed notification
func (s *NotificationService) SendPaymentConfirmed(ctx context.Context, payment *model.Payment, customer *model.Customer) error {
	data := map[string]string{
		"name":           customer.FullName,
		"amount":         fmt.Sprintf("%.0f", payment.Amount),
		"payment_number": payment.PaymentNumber,
	}
	return s.RenderAndSend(ctx, "payment_confirmed", "whatsapp", customer.Phone, data)
}

// SendPaymentReminder sends payment reminder notification
func (s *NotificationService) SendPaymentReminder(ctx context.Context, invoice *model.Invoice, customer *model.Customer) error {
	data := map[string]string{
		"name":       customer.FullName,
		"invoice_no": invoice.InvoiceNumber,
		"amount":     fmt.Sprintf("%.0f", invoice.TotalAmount),
		"due_date":   invoice.DueDate.Format("02-01-2006"),
	}
	return s.RenderAndSend(ctx, "payment_reminder", "whatsapp", customer.Phone, data)
}

// SendIsolationNotice sends isolation notice notification
func (s *NotificationService) SendIsolationNotice(ctx context.Context, customer *model.Customer, invoiceNo string) error {
	data := map[string]string{
		"name":       customer.FullName,
		"invoice_no": invoiceNo,
	}
	return s.RenderAndSend(ctx, "isolation_notice", "whatsapp", customer.Phone, data)
}

// SendRegistrationApproved sends registration approved notification
func (s *NotificationService) SendRegistrationApproved(ctx context.Context, customer *model.Customer, username, password string) error {
	data := map[string]string{
		"name":     customer.FullName,
		"username": username,
		"password": password,
	}
	return s.RenderAndSend(ctx, "registration_approved", "whatsapp", customer.Phone, data)
}

// SendRegistrationRejected sends registration rejected notification
func (s *NotificationService) SendRegistrationRejected(ctx context.Context, phone, name, reason string) error {
	data := map[string]string{
		"name":   name,
		"reason": reason,
	}
	return s.RenderAndSend(ctx, "registration_rejected", "whatsapp", phone, data)
}

// SendSuspensionWarning sends suspension warning notification
func (s *NotificationService) SendSuspensionWarning(ctx context.Context, customer *model.Customer, date, reason string) error {
	data := map[string]string{
		"name":   customer.FullName,
		"date":   date,
		"reason": reason,
	}
	return s.RenderAndSend(ctx, "suspension_warning", "whatsapp", customer.Phone, data)
}
