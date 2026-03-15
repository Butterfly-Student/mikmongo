package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mikmongo/internal/model"
)

func TestRenderTemplate(t *testing.T) {
	d := NewDomain()

	t.Run("substitutes known keys", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: "Halo {{name}}, tagihan {{amount}}"}
		data := map[string]string{"name": "Budi", "amount": "100000"}
		got, err := d.RenderTemplate(tmpl, data)
		assert.NoError(t, err)
		assert.Equal(t, "Halo Budi, tagihan 100000", got)
	})

	t.Run("unknown key remains as-is", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: "Halo {{name}}, info: {{unknown}}"}
		data := map[string]string{"name": "Budi"}
		got, err := d.RenderTemplate(tmpl, data)
		assert.NoError(t, err)
		assert.Equal(t, "Halo Budi, info: {{unknown}}", got)
	})

	t.Run("empty data map → no substitution", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: "Halo {{name}}"}
		data := map[string]string{}
		got, err := d.RenderTemplate(tmpl, data)
		assert.NoError(t, err)
		assert.Equal(t, "Halo {{name}}", got)
	})

	t.Run("empty body → error", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: ""}
		_, err := d.RenderTemplate(tmpl, map[string]string{})
		assert.Error(t, err)
	})

	t.Run("multiple occurrences of same key all replaced", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: "{{name}} dan {{name}}"}
		data := map[string]string{"name": "Budi"}
		got, err := d.RenderTemplate(tmpl, data)
		assert.NoError(t, err)
		assert.Equal(t, "Budi dan Budi", got)
	})
}

func TestRenderSubject(t *testing.T) {
	d := NewDomain()

	t.Run("substitutes subject placeholders", func(t *testing.T) {
		subject := "Tagihan {{invoice_no}}"
		tmpl := &model.MessageTemplate{Subject: &subject}
		data := map[string]string{"invoice_no": "INV000001"}
		got := d.RenderSubject(tmpl, data)
		assert.Equal(t, "Tagihan INV000001", got)
	})

	t.Run("nil subject → empty string", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Subject: nil}
		got := d.RenderSubject(tmpl, map[string]string{})
		assert.Equal(t, "", got)
	})
}

func TestExtractPlaceholders(t *testing.T) {
	d := NewDomain()

	t.Run("extracts multiple placeholders", func(t *testing.T) {
		keys := d.ExtractPlaceholders("{{name}} dan {{phone}}")
		assert.ElementsMatch(t, []string{"name", "phone"}, keys)
	})

	t.Run("no placeholders → empty slice", func(t *testing.T) {
		keys := d.ExtractPlaceholders("Hello world")
		assert.Empty(t, keys)
	})

	t.Run("duplicate placeholders → extracted once each", func(t *testing.T) {
		keys := d.ExtractPlaceholders("{{name}} bayar {{amount}} ke {{name}}")
		// "name" appears twice but should be deduped
		nameCount := 0
		for _, k := range keys {
			if k == "name" {
				nameCount++
			}
		}
		assert.Equal(t, 1, nameCount)
		assert.Contains(t, keys, "amount")
	})

	t.Run("empty string → empty slice", func(t *testing.T) {
		keys := d.ExtractPlaceholders("")
		assert.Empty(t, keys)
	})
}

func TestValidateTemplate(t *testing.T) {
	d := NewDomain()

	t.Run("valid whatsapp template", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: "Hello {{name}}", Channel: "whatsapp"}
		assert.NoError(t, d.ValidateTemplate(tmpl))
	})

	t.Run("valid email template", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: "Hello {{name}}", Channel: "email"}
		assert.NoError(t, d.ValidateTemplate(tmpl))
	})

	t.Run("body empty → error", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: "", Channel: "whatsapp"}
		assert.Error(t, d.ValidateTemplate(tmpl))
	})

	t.Run("body whitespace only → error", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: "   ", Channel: "whatsapp"}
		assert.Error(t, d.ValidateTemplate(tmpl))
	})

	t.Run("channel not in (whatsapp, email) → error", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: "Hello", Channel: "sms"}
		assert.Error(t, d.ValidateTemplate(tmpl))
	})

	t.Run("empty channel → error", func(t *testing.T) {
		tmpl := &model.MessageTemplate{Body: "Hello", Channel: ""}
		assert.Error(t, d.ValidateTemplate(tmpl))
	})
}

func TestShouldSend(t *testing.T) {
	d := NewDomain()

	t.Run("active template, matching channel → true", func(t *testing.T) {
		tmpl := &model.MessageTemplate{IsActive: true, Channel: "whatsapp"}
		assert.True(t, d.ShouldSend(tmpl, "whatsapp"))
	})

	t.Run("active template, mismatched channel → false", func(t *testing.T) {
		tmpl := &model.MessageTemplate{IsActive: true, Channel: "whatsapp"}
		assert.False(t, d.ShouldSend(tmpl, "email"))
	})

	t.Run("inactive template, matching channel → false", func(t *testing.T) {
		tmpl := &model.MessageTemplate{IsActive: false, Channel: "whatsapp"}
		assert.False(t, d.ShouldSend(tmpl, "whatsapp"))
	})
}
