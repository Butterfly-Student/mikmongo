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
	gowaClient   notification.WhatsAppSender
	emailClient  notification.EmailSender
}

// NewNotificationService creates a new notification service.
// gowaClient may be nil — WhatsApp sending will be skipped.
func NewNotificationService(
	templateRepo repository.MessageTemplateRepository,
	settingRepo repository.SystemSettingRepository,
	gowaClient notification.WhatsAppSender,
) *NotificationService {
	return &NotificationService{
		templateRepo: templateRepo,
		settingRepo:  settingRepo,
		gowaClient:   gowaClient,
	}
}

// NewNotificationServiceWithClients creates a notification service with pre-configured clients.
// Intended for testing: pass mock implementations of WhatsAppSender and EmailSender.
func NewNotificationServiceWithClients(
	templateRepo repository.MessageTemplateRepository,
	settingRepo repository.SystemSettingRepository,
	wa notification.WhatsAppSender,
	email notification.EmailSender,
) *NotificationService {
	return &NotificationService{
		templateRepo: templateRepo,
		settingRepo:  settingRepo,
		gowaClient:   wa,
		emailClient:  email,
	}
}

// initEmailClient lazily initializes the email client from DB settings.
func (s *NotificationService) initEmailClient(ctx context.Context) {
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
		return "", fmt.Errorf("failed to parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}

// SendViaWhatsApp sends a WhatsApp message
func (s *NotificationService) SendViaWhatsApp(ctx context.Context, phone, message string) error {
	if s.gowaClient == nil {
		return fmt.Errorf("WhatsApp client not configured")
	}
	return s.gowaClient.SendMessage(ctx, phone, message)
}

// SendViaEmail sends an email
func (s *NotificationService) SendViaEmail(ctx context.Context, to, subject, body string) error {
	s.initEmailClient(ctx)
	if s.emailClient == nil {
		return fmt.Errorf("Email client not configured")
	}
	return s.emailClient.SendEmail(ctx, to, subject, body)
}

// RenderAndSend renders template and sends via specified channel.
// Returns nil if the template is not found or inactive (not an error).
// Returns an error if the DB query fails or sending fails.
func (s *NotificationService) RenderAndSend(ctx context.Context, event, channel, to string, data map[string]string) error {
	tmpl, err := s.templateRepo.GetByEventAndChannel(ctx, event, channel)
	if err != nil {
		// Template not found is acceptable — skip silently
		return nil
	}
	if !tmpl.IsActive {
		return nil
	}

	body, err := s.RenderTemplate(tmpl.Body, data)
	if err != nil {
		return fmt.Errorf("render template body for %s/%s: %w", event, channel, err)
	}

	switch channel {
	case "whatsapp":
		return s.SendViaWhatsApp(ctx, to, body)
	case "email":
		subject := ""
		if tmpl.Subject != nil {
			subj, err := s.RenderTemplate(*tmpl.Subject, data)
			if err != nil {
				return fmt.Errorf("render template subject for %s/%s: %w", event, channel, err)
			}
			subject = subj
		}
		return s.SendViaEmail(ctx, to, subject, body)
	}
	return nil
}

// SendToGroup renders a template and sends it to the configured WhatsApp group.
func (s *NotificationService) SendToGroup(ctx context.Context, event string, data map[string]string) error {
	if s.gowaClient == nil {
		return nil
	}
	tmpl, err := s.templateRepo.GetByEventAndChannel(ctx, event, "whatsapp")
	if err != nil {
		return nil // template not found — skip
	}
	if !tmpl.IsActive {
		return nil
	}
	body, err := s.RenderTemplate(tmpl.Body, data)
	if err != nil {
		return fmt.Errorf("render group template for %s: %w", event, err)
	}
	return s.gowaClient.SendGroupMessage(ctx, body)
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

// SendAgentInvoiceCreated sends notification when an agent invoice is created
func (s *NotificationService) SendAgentInvoiceCreated(ctx context.Context, agent *model.SalesAgent, invoice *model.AgentInvoice) error {
	if agent.Phone == nil {
		return nil
	}
	data := map[string]string{
		"name":       agent.Name,
		"invoice_no": invoice.InvoiceNumber,
		"amount":     fmt.Sprintf("%.0f", invoice.TotalAmount),
		"due_date":   invoice.PeriodEnd.Format("02-01-2006"),
	}
	return s.RenderAndSend(ctx, "agent_invoice_created", "whatsapp", *agent.Phone, data)
}

// SendAgentInvoicePaid sends notification when an agent invoice is paid
func (s *NotificationService) SendAgentInvoicePaid(ctx context.Context, agent *model.SalesAgent, invoice *model.AgentInvoice) error {
	if agent.Phone == nil {
		return nil
	}
	data := map[string]string{
		"name":       agent.Name,
		"invoice_no": invoice.InvoiceNumber,
		"amount":     fmt.Sprintf("%.0f", invoice.PaidAmount),
	}
	return s.RenderAndSend(ctx, "agent_invoice_paid", "whatsapp", *agent.Phone, data)
}

// SendAgentInvoiceReminder sends payment reminder for an agent invoice
func (s *NotificationService) SendAgentInvoiceReminder(ctx context.Context, agent *model.SalesAgent, invoice *model.AgentInvoice) error {
	if agent.Phone == nil {
		return nil
	}
	data := map[string]string{
		"name":       agent.Name,
		"invoice_no": invoice.InvoiceNumber,
		"amount":     fmt.Sprintf("%.0f", invoice.TotalAmount),
		"due_date":   invoice.PeriodEnd.Format("02-01-2006"),
	}
	return s.RenderAndSend(ctx, "agent_invoice_reminder", "whatsapp", *agent.Phone, data)
}
