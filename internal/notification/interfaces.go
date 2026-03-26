package notification

import "context"

// WhatsAppSender defines the interface for sending WhatsApp messages.
type WhatsAppSender interface {
	SendMessage(ctx context.Context, phone, message string) error
	SendGroupMessage(ctx context.Context, message string) error
}

// EmailSender defines the interface for sending emails.
type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}
