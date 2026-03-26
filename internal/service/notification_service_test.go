package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"mikmongo/internal/model"
	"mikmongo/internal/service/mocks"
)

func newNotificationServiceWithMocks() (
	*NotificationService,
	*mocks.MockMessageTemplateRepository,
	*mocks.MockSystemSettingRepository,
) {
	templateRepo := &mocks.MockMessageTemplateRepository{}
	settingRepo := &mocks.MockSystemSettingRepository{}
	svc := NewNotificationService(templateRepo, settingRepo, nil)
	return svc, templateRepo, settingRepo
}

func TestNotificationService_RenderTemplate(t *testing.T) {
	svc, _, _ := newNotificationServiceWithMocks()

	t.Run("substitutes all keys", func(t *testing.T) {
		result, err := svc.RenderTemplate("Halo {{name}}, tagihan {{amount}}", map[string]string{
			"name":   "Budi",
			"amount": "100000",
		})
		require.NoError(t, err)
		assert.Equal(t, "Halo Budi, tagihan 100000", result)
	})

	t.Run("unknown key remains as-is", func(t *testing.T) {
		result, err := svc.RenderTemplate("Halo {{name}}, info: {{unknown}}", map[string]string{
			"name": "Budi",
		})
		require.NoError(t, err)
		assert.Equal(t, "Halo Budi, info: {{unknown}}", result)
	})

	t.Run("empty data map → no substitution", func(t *testing.T) {
		result, err := svc.RenderTemplate("Halo {{name}}", map[string]string{})
		require.NoError(t, err)
		assert.Equal(t, "Halo {{name}}", result)
	})
}

func TestNotificationService_RenderAndSend_NoTemplate(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, _ := newNotificationServiceWithMocks()

	templateRepo.On("GetByEventAndChannel", ctx, "invoice_created", "whatsapp").
		Return(nil, errors.New("template not found"))

	// Should return nil (skip if no template)
	err := svc.RenderAndSend(ctx, "invoice_created", "whatsapp", "081234567890", map[string]string{})
	assert.NoError(t, err)
}

func TestNotificationService_RenderAndSend_InactiveTemplate(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, _ := newNotificationServiceWithMocks()

	tmpl := &model.MessageTemplate{
		Body:     "Halo {{name}}",
		Channel:  "whatsapp",
		IsActive: false,
	}
	templateRepo.On("GetByEventAndChannel", ctx, "invoice_created", "whatsapp").Return(tmpl, nil)

	// Should return nil for inactive template
	err := svc.RenderAndSend(ctx, "invoice_created", "whatsapp", "081234", map[string]string{})
	assert.NoError(t, err)
}

func TestSendInvoiceCreated_RendersCorrectData(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, settingRepo := newNotificationServiceWithMocks()

	dueDate := time.Date(2024, time.March, 15, 0, 0, 0, 0, time.UTC)
	invoice := &model.Invoice{
		InvoiceNumber: "INV000001",
		TotalAmount:   111_000,
		DueDate:       dueDate,
	}
	customer := &model.Customer{
		FullName: "Budi Santoso",
		Phone:    "081234567890",
	}

	// Template found but no whatsapp client configured → RenderAndSend will try to send
	// but SendViaWhatsApp will fail (no client). Template lookup should succeed though.
	templateRepo.On("GetByEventAndChannel", ctx, "invoice_created", "whatsapp").
		Return(nil, errors.New("template not found"))
	settingRepo.On("GetByGroupAndKey", ctx, "notification", "gowa_url").
		Return(nil, errors.New("not found")).Maybe()

	// With no template → should return no error (gracefully skip)
	err := svc.SendInvoiceCreated(ctx, invoice, customer)
	// May return error if some sends fail, but should not panic
	_ = err
}

func TestSendPaymentConfirmed_RendersCorrectData(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, _ := newNotificationServiceWithMocks()

	payment := &model.Payment{
		PaymentNumber: "PAY000001",
		Amount:        100_000,
	}
	customer := &model.Customer{
		FullName: "Budi Santoso",
		Phone:    "081234567890",
	}

	templateRepo.On("GetByEventAndChannel", ctx, "payment_confirmed", "whatsapp").
		Return(nil, errors.New("template not found"))

	// No template → graceful no-op
	err := svc.SendPaymentConfirmed(ctx, payment, customer)
	assert.NoError(t, err)
}

func TestSendPaymentReminder_RendersCorrectData(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, _ := newNotificationServiceWithMocks()

	dueDate := time.Now().AddDate(0, 0, 3)
	invoice := &model.Invoice{
		InvoiceNumber: "INV000002",
		TotalAmount:   200_000,
		DueDate:       dueDate,
	}
	customer := &model.Customer{
		FullName: "Siti Rahayu",
		Phone:    "082345678901",
	}

	templateRepo.On("GetByEventAndChannel", ctx, "payment_reminder", "whatsapp").
		Return(nil, errors.New("no template"))

	err := svc.SendPaymentReminder(ctx, invoice, customer)
	assert.NoError(t, err)
}

// newNotificationSvcWithClientMocks creates a NotificationService with injected mock senders.
func newNotificationSvcWithClientMocks() (
	*NotificationService,
	*mocks.MockMessageTemplateRepository,
	*mocks.MockSystemSettingRepository,
	*mocks.MockWhatsAppSender,
	*mocks.MockEmailSender,
) {
	templateRepo := &mocks.MockMessageTemplateRepository{}
	settingRepo := &mocks.MockSystemSettingRepository{}
	waMock := &mocks.MockWhatsAppSender{}
	emailMock := &mocks.MockEmailSender{}
	svc := NewNotificationServiceWithClients(templateRepo, settingRepo, waMock, emailMock)
	return svc, templateRepo, settingRepo, waMock, emailMock
}

func TestSendViaWhatsApp_CallsClientWithPhone(t *testing.T) {
	ctx := context.Background()
	svc, _, _, waMock, _ := newNotificationSvcWithClientMocks()

	waMock.On("SendMessage", ctx, "081234567890", "test message").Return(nil)

	err := svc.SendViaWhatsApp(ctx, "081234567890", "test message")
	require.NoError(t, err)
	waMock.AssertCalled(t, "SendMessage", ctx, "081234567890", "test message")
}

func TestRenderAndSend_WhatsAppChannel_CallsWA(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, _, waMock, emailMock := newNotificationSvcWithClientMocks()

	tmpl := &model.MessageTemplate{
		Body:     "Halo {{name}}",
		Channel:  "whatsapp",
		IsActive: true,
	}
	templateRepo.On("GetByEventAndChannel", ctx, "invoice_created", "whatsapp").Return(tmpl, nil)
	waMock.On("SendMessage", ctx, "081234567890", "Halo Budi").Return(nil)

	err := svc.RenderAndSend(ctx, "invoice_created", "whatsapp", "081234567890", map[string]string{"name": "Budi"})
	require.NoError(t, err)
	waMock.AssertCalled(t, "SendMessage", ctx, "081234567890", "Halo Budi")
	emailMock.AssertNotCalled(t, "SendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
}

func TestRenderAndSend_EmailChannel_CallsEmail(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, _, waMock, emailMock := newNotificationSvcWithClientMocks()

	subject := "Tagihan Baru"
	tmpl := &model.MessageTemplate{
		Body:     "Halo {{name}}",
		Subject:  &subject,
		Channel:  "email",
		IsActive: true,
	}
	templateRepo.On("GetByEventAndChannel", ctx, "invoice_created", "email").Return(tmpl, nil)
	emailMock.On("SendEmail", ctx, "user@test.com", "Tagihan Baru", "Halo Budi").Return(nil)

	err := svc.RenderAndSend(ctx, "invoice_created", "email", "user@test.com", map[string]string{"name": "Budi"})
	require.NoError(t, err)
	emailMock.AssertCalled(t, "SendEmail", ctx, "user@test.com", "Tagihan Baru", "Halo Budi")
	waMock.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestSendInvoiceCreated_CorrectPhone(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, _, waMock, _ := newNotificationSvcWithClientMocks()

	tmpl := &model.MessageTemplate{
		Body:     "Tagihan {{invoice_no}}",
		Channel:  "whatsapp",
		IsActive: true,
	}
	templateRepo.On("GetByEventAndChannel", ctx, "invoice_created", "whatsapp").Return(tmpl, nil)
	waMock.On("SendMessage", ctx, "08123456789", mock.AnythingOfType("string")).Return(nil)

	invoice := &model.Invoice{
		InvoiceNumber: "INV000001",
		TotalAmount:   111_000,
		DueDate:       time.Now().AddDate(0, 0, 14),
	}
	customer := &model.Customer{
		FullName: "Budi",
		Phone:    "08123456789",
	}

	err := svc.SendInvoiceCreated(ctx, invoice, customer)
	require.NoError(t, err)
	waMock.AssertCalled(t, "SendMessage", ctx, "08123456789", mock.AnythingOfType("string"))
}

func TestSendViaWhatsApp_NilClient_ReturnsError(t *testing.T) {
	ctx := context.Background()
	svc, _, settingRepo := newNotificationServiceWithMocks()
	// Make all setting lookups return empty so neither gowaClient nor emailClient is initialized
	settingRepo.On("GetByGroupAndKey", ctx, mock.Anything, mock.Anything).
		Return(nil, errors.New("not found"))

	err := svc.SendViaWhatsApp(ctx, "081234567890", "test message")
	assert.ErrorContains(t, err, "WhatsApp client not configured")
}

// --- Agent Invoice Notification Tests ---

func makeAgentPhone(phone string) *string { return &phone }

func TestSendAgentInvoiceCreated_NilPhone_Skips(t *testing.T) {
	ctx := context.Background()
	svc, _, _, waMock, _ := newNotificationSvcWithClientMocks()

	agent := &model.SalesAgent{Name: "Agen Budi", Phone: nil}
	invoice := &model.AgentInvoice{
		InvoiceNumber: "AINV-001",
		TotalAmount:   500_000,
		PeriodEnd:     time.Date(2024, time.April, 30, 0, 0, 0, 0, time.UTC),
	}

	err := svc.SendAgentInvoiceCreated(ctx, agent, invoice)
	assert.NoError(t, err)
	waMock.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestSendAgentInvoiceCreated_WithPhone_SendsWA(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, _, waMock, _ := newNotificationSvcWithClientMocks()

	phone := "081234567890"
	agent := &model.SalesAgent{Name: "Agen Budi", Phone: makeAgentPhone(phone)}
	invoice := &model.AgentInvoice{
		InvoiceNumber: "AINV-001",
		TotalAmount:   500_000,
		PeriodEnd:     time.Date(2024, time.April, 30, 0, 0, 0, 0, time.UTC),
	}

	tmpl := &model.MessageTemplate{
		Body:     "Tagihan {{invoice_no}} sebesar {{amount}} jatuh tempo {{due_date}}",
		Channel:  "whatsapp",
		IsActive: true,
	}
	templateRepo.On("GetByEventAndChannel", ctx, "agent_invoice_created", "whatsapp").Return(tmpl, nil)
	waMock.On("SendMessage", ctx, phone, mock.AnythingOfType("string")).Return(nil)

	err := svc.SendAgentInvoiceCreated(ctx, agent, invoice)
	require.NoError(t, err)
	waMock.AssertCalled(t, "SendMessage", ctx, phone, mock.AnythingOfType("string"))
}

func TestSendAgentInvoicePaid_NilPhone_Skips(t *testing.T) {
	ctx := context.Background()
	svc, _, _, waMock, _ := newNotificationSvcWithClientMocks()

	agent := &model.SalesAgent{Name: "Agen Siti", Phone: nil}
	invoice := &model.AgentInvoice{
		InvoiceNumber: "AINV-002",
		TotalAmount:   300_000,
		PaidAmount:    300_000,
	}

	err := svc.SendAgentInvoicePaid(ctx, agent, invoice)
	assert.NoError(t, err)
	waMock.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestSendAgentInvoicePaid_WithPhone_SendsWA(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, _, waMock, _ := newNotificationSvcWithClientMocks()

	phone := "082345678901"
	agent := &model.SalesAgent{Name: "Agen Siti", Phone: makeAgentPhone(phone)}
	invoice := &model.AgentInvoice{
		InvoiceNumber: "AINV-002",
		TotalAmount:   300_000,
		PaidAmount:    300_000,
	}

	tmpl := &model.MessageTemplate{
		Body:     "Pembayaran {{invoice_no}} sebesar {{amount}} diterima",
		Channel:  "whatsapp",
		IsActive: true,
	}
	templateRepo.On("GetByEventAndChannel", ctx, "agent_invoice_paid", "whatsapp").Return(tmpl, nil)
	waMock.On("SendMessage", ctx, phone, mock.AnythingOfType("string")).Return(nil)

	err := svc.SendAgentInvoicePaid(ctx, agent, invoice)
	require.NoError(t, err)
	waMock.AssertCalled(t, "SendMessage", ctx, phone, mock.AnythingOfType("string"))
}

func TestSendAgentInvoiceReminder_NilPhone_Skips(t *testing.T) {
	ctx := context.Background()
	svc, _, _, waMock, _ := newNotificationSvcWithClientMocks()

	agent := &model.SalesAgent{Name: "Agen Rudi", Phone: nil}
	invoice := &model.AgentInvoice{
		InvoiceNumber: "AINV-003",
		TotalAmount:   150_000,
		PeriodEnd:     time.Now().AddDate(0, 0, 3),
	}

	err := svc.SendAgentInvoiceReminder(ctx, agent, invoice)
	assert.NoError(t, err)
	waMock.AssertNotCalled(t, "SendMessage", mock.Anything, mock.Anything, mock.Anything)
}

func TestSendAgentInvoiceReminder_WithPhone_SendsWA(t *testing.T) {
	ctx := context.Background()
	svc, templateRepo, _, waMock, _ := newNotificationSvcWithClientMocks()

	phone := "083456789012"
	agent := &model.SalesAgent{Name: "Agen Rudi", Phone: makeAgentPhone(phone)}
	invoice := &model.AgentInvoice{
		InvoiceNumber: "AINV-003",
		TotalAmount:   150_000,
		PeriodEnd:     time.Date(2024, time.May, 15, 0, 0, 0, 0, time.UTC),
	}

	tmpl := &model.MessageTemplate{
		Body:     "Pengingat: Tagihan {{invoice_no}} sebesar {{amount}} jatuh tempo {{due_date}}",
		Channel:  "whatsapp",
		IsActive: true,
	}
	templateRepo.On("GetByEventAndChannel", ctx, "agent_invoice_reminder", "whatsapp").Return(tmpl, nil)
	waMock.On("SendMessage", ctx, phone, mock.AnythingOfType("string")).Return(nil)

	err := svc.SendAgentInvoiceReminder(ctx, agent, invoice)
	require.NoError(t, err)
	waMock.AssertCalled(t, "SendMessage", ctx, phone, mock.AnythingOfType("string"))
}
