package rabbitmq

// UpdatePayload defines the structure of the message sent for settings updates.
type UpdatePayload struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
