package notification

import "context"

// WhatsAppSender defines the interface for sending WhatsApp messages.
// *GoWAClient implements this interface.
type WhatsAppSender interface {
	SendMessage(ctx context.Context, phone, message string) error
}

// EmailSender defines the interface for sending emails.
// *EmailClient implements this interface.
type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}
