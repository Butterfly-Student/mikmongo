package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() { goose.AddMigrationContext(up020, down020) }

func up020(ctx context.Context, tx *sql.Tx) error {
	// Seed system_settings
	seeds := []struct {
		group, key, value, typ, label string
	}{
		{"isolate", "pppoe_profile", "isolate", "string", "PPPoE Isolate Profile Name"},
		{"notification", "gowa_url", "", "string", "GoWA REST API URL"},
		{"notification", "gowa_sender", "", "string", "WhatsApp Sender Number"},
		{"billing", "due_days", "10", "integer", "Invoice Due Days"},
		{"billing", "reminder_intervals", "3,7,1", "string", "Reminder Intervals (days before due)"},
		{"billing", "late_fee_after_days", "0", "integer", "Late Fee After Days Overdue"},
	}

	for _, s := range seeds {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO system_settings (group_name, key_name, value, type, label)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (group_name, key_name) DO NOTHING
		`, s.group, s.key, s.value, s.typ, s.label)
		if err != nil {
			return err
		}
	}

	// Seed sequence_counters for invoice, payment, customer_code
	counters := []struct {
		name, prefix string
		padding      int
	}{
		{"invoice_number", "INV", 6},
		{"payment_number", "PAY", 6},
		{"customer_code", "CST", 5},
	}
	for _, c := range counters {
		_, err := tx.ExecContext(ctx, `
			INSERT INTO sequence_counters (name, prefix, padding, last_number)
			VALUES ($1, $2, $3, 0)
			ON CONFLICT (name) DO NOTHING
		`, c.name, c.prefix, c.padding)
		if err != nil {
			return err
		}
	}

	// Seed message_templates
	templates := []struct {
		event, channel, subject, body string
	}{
		{"invoice_created", "whatsapp", "", "Halo {{name}},\n\nInvoice {{invoice_no}} sebesar Rp{{amount}} telah dibuat.\nJatuh tempo: {{due_date}}\n\nSilakan lakukan pembayaran sebelum tanggal jatuh tempo.\n\nTerima kasih."},
		{"invoice_created", "email", "Invoice {{invoice_no}} - Rp{{amount}}", "Halo {{name}},\n\nInvoice {{invoice_no}} sebesar Rp{{amount}} telah dibuat.\nJatuh tempo: {{due_date}}\n\nSilakan lakukan pembayaran sebelum tanggal jatuh tempo.\n\nTerima kasih."},
		{"payment_reminder", "whatsapp", "", "Halo {{name}},\n\nPengingat: Invoice {{invoice_no}} sebesar Rp{{amount}} akan jatuh tempo pada {{due_date}}.\n\nSegera lakukan pembayaran untuk menghindari isolasi layanan.\n\nTerima kasih."},
		{"payment_confirmed", "whatsapp", "", "Halo {{name}},\n\nPembayaran Rp{{amount}} untuk invoice {{invoice_no}} telah dikonfirmasi.\n\nTerima kasih atas pembayaran Anda."},
		{"isolation_notice", "whatsapp", "", "Halo {{name}},\n\nLayanan internet Anda telah dibatasi karena tagihan {{invoice_no}} belum dibayar.\n\nSilakan segera lakukan pembayaran untuk memulihkan layanan.\n\nTerima kasih."},
		{"registration_approved", "whatsapp", "", "Halo {{name}},\n\nPendaftaran Anda telah disetujui!\nUsername: {{username}}\nPassword: {{password}}\n\nSelamat menikmati layanan kami."},
		{"registration_rejected", "whatsapp", "", "Halo {{name}},\n\nMohon maaf, pendaftaran Anda ditolak.\nAlasan: {{reason}}\n\nSilakan hubungi kami untuk informasi lebih lanjut."},
		{"suspension_warning", "whatsapp", "", "Halo {{name}},\n\nLayanan Anda akan dinonaktifkan mulai {{date}} karena: {{reason}}.\n\nHubungi kami untuk informasi lebih lanjut."},
	}
	for _, t := range templates {
		if t.subject == "" {
			_, err := tx.ExecContext(ctx, `
				INSERT INTO message_templates (event, channel, body, is_active)
				VALUES ($1, $2, $3, true)
				ON CONFLICT (event, channel) DO NOTHING
			`, t.event, t.channel, t.body)
			if err != nil {
				return err
			}
		} else {
			_, err := tx.ExecContext(ctx, `
				INSERT INTO message_templates (event, channel, subject, body, is_active)
				VALUES ($1, $2, $3, $4, true)
				ON CONFLICT (event, channel) DO NOTHING
			`, t.event, t.channel, t.subject, t.body)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func down020(ctx context.Context, tx *sql.Tx) error {
	// Seeder migration - no rollback needed
	return nil
}
