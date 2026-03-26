package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockWhatsAppSender is a mock implementation of notification.WhatsAppSender
type MockWhatsAppSender struct {
	mock.Mock
}

func (m *MockWhatsAppSender) SendMessage(ctx context.Context, phone, message string) error {
	args := m.Called(ctx, phone, message)
	return args.Error(0)
}

func (m *MockWhatsAppSender) SendGroupMessage(ctx context.Context, message string) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

// MockEmailSender is a mock implementation of notification.EmailSender
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmail(ctx context.Context, to, subject, body string) error {
	args := m.Called(ctx, to, subject, body)
	return args.Error(0)
}
