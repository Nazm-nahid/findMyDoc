package email

import (
	"fmt"
	"net/smtp"
)

type SMTPService struct {
	SMTPHost string
	SMTPPort string
	Sender   string
	Password string
}

func NewSMTPService(host, port, sender, password string) *SMTPService {
	return &SMTPService{host, port, sender, password}
}

func (s *SMTPService) SendVerificationEmail(to, verificationCode string) error {
	auth := smtp.PlainAuth("", s.Sender, s.Password, s.SMTPHost)
	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)

	// MIME headers for HTML
	subject := "Subject: Verify your email for findMyDoc\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// HTML body
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; background-color: #f9f9f9; padding: 20px;">
			<div style="max-width: 500px; margin: auto; background: #ffffff; border-radius: 8px; padding: 30px; box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);">
				<h2 style="color: #2E86C1;">Welcome to <span style="color: #28B463;">findMyDoc</span> 👋</h2>
				<p style="font-size: 16px;">Thank you for registering! Please verify your email address to complete your registration.</p>
				<p style="font-size: 16px;">Your verification code is:</p>
				<div style="font-size: 24px; font-weight: bold; color: #333; background: #f0f0f0; padding: 10px 20px; display: inline-block; border-radius: 6px;">
					%s
				</div>
				<p style="font-size: 14px; margin-top: 20px;">If you did not sign up for findMyDoc, please ignore this email.</p>
				<p style="font-size: 14px; color: #888;">&mdash; The findMyDoc Team</p>
			</div>
		</body>
		</html>`, verificationCode)

	// Combine all
	msg := []byte(subject + mime + body)

	return smtp.SendMail(addr, auth, s.Sender, []string{to}, msg)
}
