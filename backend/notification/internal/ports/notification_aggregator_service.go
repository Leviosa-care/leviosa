package ports

type Notification interface {
	MailService
	SMSService
}
