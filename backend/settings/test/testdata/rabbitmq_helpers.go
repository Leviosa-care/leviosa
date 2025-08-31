package testdata

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	mq "github.com/Leviosa-care/core/contracts/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RabbitMQMessage represents a consumed message for verification
type RabbitMQMessage struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// ConsumeSettingsMessage consumes one message from a settings queue with timeout
func ConsumeSettingsMessage(t *testing.T, ch *amqp.Channel, queueName string, timeout time.Duration) (*RabbitMQMessage, error) {
	t.Helper()

	// Consume from the queue
	msgs, err := ch.Consume(
		queueName,
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume from queue %s: %w", queueName, err)
	}

	// Wait for message with timeout
	select {
	case msg := <-msgs:
		var payload RabbitMQMessage
		if err := json.Unmarshal(msg.Body, &payload); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}
		return &payload, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("timeout waiting for message in queue %s after %v", queueName, timeout)
	}
}

// VerifySettingsUpdateMessage verifies that a settings update message was published correctly
func VerifySettingsUpdateMessage(t *testing.T, ch *amqp.Channel, expectedKey string, expectedValue interface{}) {
	t.Helper()

	// Check both queues since settings updates are published to both
	queues := []string{mq.NotificationSettingsQueueName, mq.OTPSettingsQueueName}

	for _, queueName := range queues {
		msg, err := ConsumeSettingsMessage(t, ch, queueName, 2*time.Second)
		require.NoError(t, err, "Failed to consume message from queue %s", queueName)
		require.NotNil(t, msg, "No message received from queue %s", queueName)

		assert.Equal(t, expectedKey, msg.Key, "Message key mismatch in queue %s", queueName)

		// Handle type conversion for numeric values - JSON unmarshaling converts numbers to float64
		if expectedInt, ok := expectedValue.(int); ok {
			if actualFloat, ok := msg.Value.(float64); ok {
				assert.Equal(t, float64(expectedInt), actualFloat, "Message value mismatch in queue %s", queueName)
			} else {
				assert.Equal(t, expectedValue, msg.Value, "Message value mismatch in queue %s", queueName)
			}
		} else {
			assert.Equal(t, expectedValue, msg.Value, "Message value mismatch in queue %s", queueName)
		}
	}
}

// VerifyNoSettingsUpdateMessage verifies that no unexpected messages are in the queues
func VerifyNoSettingsUpdateMessage(t *testing.T, ch *amqp.Channel) {
	t.Helper()

	queues := []string{mq.NotificationSettingsQueueName, mq.OTPSettingsQueueName}

	for _, queueName := range queues {
		msg, err := ConsumeSettingsMessage(t, ch, queueName, 100*time.Millisecond)
		// We expect an error (timeout) when no messages are present
		if err == nil {
			t.Errorf("Unexpected message found in queue %s: key=%s, value=%v",
				queueName, msg.Key, msg.Value)
		}
	}
}

// PurgeSettingsQueues purges all messages from settings queues (useful for cleanup)
func PurgeSettingsQueues(t *testing.T, ch *amqp.Channel) {
	t.Helper()

	queues := []string{mq.NotificationSettingsQueueName, mq.OTPSettingsQueueName}

	for _, queueName := range queues {
		_, err := ch.QueuePurge(queueName, false)
		require.NoError(t, err, "Failed to purge queue %s", queueName)
	}
}

// GetRabbitMQChannel creates a new channel from the existing connection for testing
// This assumes the connection is available globally in the test setup
func GetRabbitMQChannel(t *testing.T, conn *amqp.Connection) *amqp.Channel {
	t.Helper()

	ch, err := conn.Channel()
	require.NoError(t, err, "Failed to create RabbitMQ channel")

	return ch
}
