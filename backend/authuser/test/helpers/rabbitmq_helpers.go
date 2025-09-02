package helpers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Leviosa-care/authuser/internal/domain"

	mq "github.com/Leviosa-care/core/contracts/rabbitmq"
	"github.com/Leviosa-care/core/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// SetupOTPQueue declares and binds the OTP notification queue for testing
func SetupOTPQueue(t *testing.T, ch *amqp.Channel) {
	t.Helper()

	// Declare exchange
	err := ch.ExchangeDeclare(
		mq.OTPNotificationExchangeName,
		"direct",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	require.NoError(t, err, "Failed to declare OTP exchange")

	// Declare queue
	queue, err := ch.QueueDeclare(
		"test_otp_notifications", // name
		false,                    // durable
		true,                     // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		nil,                      // arguments
	)
	require.NoError(t, err, "Failed to declare OTP queue")

	// Bind queue to exchange
	err = ch.QueueBind(
		queue.Name,                     // queue name
		mq.OTPEmailRoutingKey,          // routing key
		mq.OTPNotificationExchangeName, // exchange
		false,                          // no-wait
		nil,                            // arguments
	)
	require.NoError(t, err, "Failed to bind OTP queue")
}

// ConsumeOTPMessages starts consuming OTP messages from the test queue
func ConsumeOTPMessages(t *testing.T, ch *amqp.Channel) <-chan amqp.Delivery {
	t.Helper()

	msgs, err := ch.Consume(
		"test_otp_notifications", // queue
		"test_consumer",          // consumer
		false,                    // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	require.NoError(t, err, "Failed to start consuming messages")

	return msgs
}

// VerifyOTPMessage verifies that an OTP message was published correctly
func VerifyOTPMessage(t *testing.T, delivery amqp.Delivery, expectedEmail string, expectedCode string) {
	t.Helper()

	// Parse message payload
	var payload rabbitmq.UpdatePayload
	err := json.Unmarshal(delivery.Body, &payload)
	require.NoError(t, err, "Failed to unmarshal OTP message")

	// Verify key (should be email)
	assert.Equal(t, expectedEmail, payload.Key, "OTP message key mismatch")

	// Parse value as OTP event
	valueBytes, err := json.Marshal(payload.Value)
	require.NoError(t, err, "Failed to marshal payload value")

	var otpEvent domain.OTPSentEvent
	err = json.Unmarshal(valueBytes, &otpEvent)
	require.NoError(t, err, "Failed to unmarshal OTP event")

	// Verify OTP event data
	assert.Equal(t, expectedEmail, otpEvent.Email, "OTP event email mismatch")
	assert.Equal(t, expectedCode, otpEvent.Code, "OTP event code mismatch")
	assert.True(t, otpEvent.ExpiresAt.After(time.Now()), "OTP should not be expired")

	// Acknowledge the message
	delivery.Ack(false)
}

// WaitForOTPMessage waits for an OTP message with timeout
func WaitForOTPMessage(t *testing.T, msgs <-chan amqp.Delivery, timeout time.Duration) amqp.Delivery {
	t.Helper()

	select {
	case msg := <-msgs:
		return msg
	case <-time.After(timeout):
		t.Fatal("Timeout waiting for OTP message")
		return amqp.Delivery{}
	}
}

// VerifyNoOTPMessage ensures no OTP message was sent (for failure cases)
func VerifyNoOTPMessage(t *testing.T, msgs <-chan amqp.Delivery, waitTime time.Duration) {
	t.Helper()

	select {
	case msg := <-msgs:
		msg.Nack(false, true) // Reject and requeue
		t.Fatal("Expected no OTP message but received one")
	case <-time.After(waitTime):
		// Expected - no message received
	}
}

// PurgeOTPQueue clears all messages from the test OTP queue
func PurgeOTPQueue(t *testing.T, ch *amqp.Channel) {
	t.Helper()

	_, err := ch.QueuePurge("test_otp_notifications", false)
	require.NoError(t, err, "Failed to purge OTP test queue")
}
