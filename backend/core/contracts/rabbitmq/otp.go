package rabbitmq

const (
	// exchange
	OTPNotificationExchangeName = "otp.notification.exchange"

	// routing keys
	OTPEmailRoutingKey = "otp.notification.email"
	OTPSMSRoutingKey   = "otp.notification.sms"

	// queues
	OTPEmailQueueName = "otp.notification.email.queue"
	OTPSMSQueueName   = "otp.notification.sms.queue"
)
