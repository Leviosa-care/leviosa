package rabbitmq

const (
	// settings change tasks constants
	SettingsExchangeName  = "settings.exchange"
	MailSettingsQueueName = "mail.settings.queue"
	OTPSettingsQueueName  = "otp.settings.queue"
	SettingsRoutingKey    = "settings.updated" // Generic routing key for settings updates
)
