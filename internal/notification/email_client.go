package notification

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

// EmailClient handles email notifications via SMTP
type EmailClient struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// NewEmailClient creates a new email client
func NewEmailClient(host, port, username, password, from string) *EmailClient {
	return &EmailClient{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// SendEmail sends an email
func (c *EmailClient) SendEmail(ctx context.Context, to, subject, body string) error {
	if c.host == "" {
		return fmt.Errorf("SMTP host not configured")
	}

	addr := net.JoinHostPort(c.host, c.port)

	msg := buildMIMEMessage(c.from, to, subject, body)

	var auth smtp.Auth
	if c.username != "" {
		auth = smtp.PlainAuth("", c.username, c.password, c.host)
	}

	// Try TLS first, fallback to plain
	if c.port == "465" {
		return c.sendTLS(addr, auth, to, msg)
	}
	return smtp.SendMail(addr, auth, c.from, []string{to}, []byte(msg))
}

func (c *EmailClient) sendTLS(addr string, auth smtp.Auth, to, msg string) error {
	tlsConfig := &tls.Config{ServerName: c.host}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("TLS dial failed: %w", err)
	}
	client, err := smtp.NewClient(conn, c.host)
	if err != nil {
		return fmt.Errorf("SMTP client failed: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP auth failed: %w", err)
		}
	}
	if err := client.Mail(c.from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = w.Write([]byte(msg))
	return err
}

func buildMIMEMessage(from, to, subject, body string) string {
	var sb strings.Builder
	sb.WriteString("From: " + from + "\r\n")
	sb.WriteString("To: " + to + "\r\n")
	sb.WriteString("Subject: " + subject + "\r\n")
	sb.WriteString("MIME-Version: 1.0\r\n")
	sb.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	sb.WriteString("\r\n")
	sb.WriteString(body)
	return sb.String()
}
