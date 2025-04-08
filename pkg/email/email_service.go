package email

type EmailService interface {
	SendVerificationEmail(to, link string) error
}
