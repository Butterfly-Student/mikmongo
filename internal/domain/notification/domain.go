// Package notification contains notification domain logic
package notification

import (
	"errors"
	"mikmongo/internal/model"
	"strings"
)

// Domain represents notification business logic
type Domain struct{}

// NewDomain creates a new notification domain
func NewDomain() *Domain {
	return &Domain{}
}

// RenderTemplate replaces all {{key}} placeholders in Body with values from data
func (d *Domain) RenderTemplate(tmpl *model.MessageTemplate, data map[string]string) (string, error) {
	if tmpl.Body == "" {
		return "", errors.New("template body is empty")
	}
	body := tmpl.Body
	for k, v := range data {
		body = strings.ReplaceAll(body, "{{"+k+"}}", v)
	}
	return body, nil
}

// RenderSubject replaces {{key}} placeholders in Subject; returns "" if Subject is nil
func (d *Domain) RenderSubject(tmpl *model.MessageTemplate, data map[string]string) string {
	if tmpl.Subject == nil {
		return ""
	}
	subject := *tmpl.Subject
	for k, v := range data {
		subject = strings.ReplaceAll(subject, "{{"+k+"}}", v)
	}
	return subject
}

// ShouldSend returns true if template is active and channel matches
func (d *Domain) ShouldSend(tmpl *model.MessageTemplate, channel string) bool {
	return tmpl.IsActive && tmpl.Channel == channel
}

// ExtractPlaceholders finds all unique {{key}} placeholders in body
func (d *Domain) ExtractPlaceholders(body string) []string {
	seen := make(map[string]struct{})
	var result []string
	remaining := body
	for {
		start := strings.Index(remaining, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(remaining[start:], "}}")
		if end == -1 {
			break
		}
		key := remaining[start+2 : start+end]
		if _, exists := seen[key]; !exists {
			seen[key] = struct{}{}
			result = append(result, key)
		}
		remaining = remaining[start+end+2:]
	}
	return result
}

// ValidateTemplate validates that Body is set and Channel is whatsapp or email
func (d *Domain) ValidateTemplate(tmpl *model.MessageTemplate) error {
	if strings.TrimSpace(tmpl.Body) == "" {
		return errors.New("template body is required")
	}
	if tmpl.Channel != "whatsapp" && tmpl.Channel != "email" {
		return errors.New("channel must be whatsapp or email")
	}
	return nil
}
