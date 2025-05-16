package rabbitmq

// SettingsUpdatePayload defines the structure of the message sent for settings updates.
type SettingsUpdatePayload struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
