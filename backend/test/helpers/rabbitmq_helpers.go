package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"

	mq "github.com/Leviosa-care/leviosa/backend/internal/common/contracts/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockNotificationService is a thread-safe mock that records SendOTPEmail calls
// for verification in integration tests. It satisfies authuser/ports.NotificationService.
type MockNotificationService struct {
	mu    sync.Mutex
	calls []OTPEmailCall
}

// OTPEmailCall records a single SendOTPEmail invocation.
type OTPEmailCall struct {
	Email string
	OTP   string
}

// NewMockNotificationService creates a new mock notification service.
func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{}
}

func (m *MockNotificationService) SendOTPEmail(_ context.Context, email, otp string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = append(m.calls, OTPEmailCall{Email: email, OTP: otp})
	return nil
}

// GetCalls returns a copy of the recorded calls.
func (m *MockNotificationService) GetCalls() []OTPEmailCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	result := make([]OTPEmailCall, len(m.calls))
	copy(result, m.calls)
	return result
}

// Reset clears all recorded calls.
func (m *MockNotificationService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = nil
}

// Compile-time check
var _ ports.NotificationService = (*MockNotificationService)(nil)

// AssertOTPReceived asserts that an OTP email was sent to the given email address
// and returns the OTP code.
func AssertOTPReceived(t *testing.T, svc *MockNotificationService, expectedEmail string) string {
	t.Helper()
	calls := svc.GetCalls()
	require.NotEmpty(t, calls, "expected at least one OTP email to be sent")
	last := calls[len(calls)-1]
	assert.Equal(t, expectedEmail, last.Email)
	assert.NotEmpty(t, last.OTP)
	return last.OTP
}

// AssertNoOTPSent asserts that no OTP email was sent.
func AssertNoOTPSent(t *testing.T, svc *MockNotificationService) {
	t.Helper()
	calls := svc.GetCalls()
	assert.Empty(t, calls, "expected no OTP emails to be sent")
}

// AssertOTPSentCount asserts the exact number of OTP emails sent.
func AssertOTPSentCount(t *testing.T, svc *MockNotificationService, expected int) {
	t.Helper()
	calls := svc.GetCalls()
	assert.Len(t, calls, expected, "expected %d OTP emails to be sent, got %d", expected, len(calls))
}

// Settings-specific RabbitMQ helpers

// RabbitMQMessage represents a consumed message for verification
type RabbitMQMessage struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// GetRabbitMQChannel gets a RabbitMQ channel from connection
func GetRabbitMQChannel(t *testing.T, conn *amqp.Connection) *amqp.Channel {
	t.Helper()

	ch, err := conn.Channel()
	require.NoError(t, err, "Failed to open RabbitMQ channel")
	return ch
}

// PurgeSettingsQueues clears all messages from the test settings queues
func PurgeSettingsQueues(t *testing.T, ch *amqp.Channel) {
	t.Helper()

	queues := []string{mq.NotificationSettingsQueueName, mq.OTPSettingsQueueName}

	for _, queue := range queues {
		_, err := ch.QueuePurge(queue, false)
		require.NoError(t, err, "Failed to purge settings test queue: "+queue)
	}
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

// VerifySettingsUpdateMessage consumes and verifies a settings update message
func VerifySettingsUpdateMessage(t *testing.T, ch *amqp.Channel, expectedKey, expectedValue interface{}) {
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
