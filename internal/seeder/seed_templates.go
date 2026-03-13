package seeder

import "context"

func (s *Seeder) seedMessageTemplates(ctx context.Context) error {
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
		// Use the same query structure but with NULL for empty subject to avoid type inference issues
		_, err := s.db.ExecContext(ctx, `
			INSERT INTO message_templates (event, channel, subject, body, is_active)
			VALUES ($1, $2, NULLIF($3, ''), $4, true)
			ON CONFLICT (event, channel) DO NOTHING
		`, t.event, t.channel, t.subject, t.body)
		if err != nil {
			return err
		}
	}

	return nil
}
