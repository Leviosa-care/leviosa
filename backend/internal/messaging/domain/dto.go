package domain

import (
	"time"

	"github.com/google/uuid"
)

// ThreadSummary is the DTO returned when listing threads for a user.
type ThreadSummary struct {
	ThreadID        uuid.UUID `json:"thread_id"`
	ParticipantID   uuid.UUID `json:"participant_id"`
	ParticipantName string    `json:"participant_name"`
	LastMessage     string    `json:"last_message"`
	LastMessageAt   time.Time `json:"last_message_at"`
	UnreadCount     int       `json:"unread_count"`
}

// ThreadSummaryRaw is the repository-layer view of a thread, before names are
// resolved and the last-message body is decrypted.
type ThreadSummaryRaw struct {
	ThreadID          uuid.UUID
	ParticipantID     uuid.UUID
	LastBodyEncrypted []byte     // nil when the thread has no messages yet
	LastDEKEncrypted  []byte
	LastKeyVersion    int
	LastMessageAt     *time.Time // nil when the thread has no messages yet
	UnreadCount       int
}

// MessageResponse is the DTO returned for individual messages
type MessageResponse struct {
	ID        uuid.UUID  `json:"id"`
	ThreadID  uuid.UUID  `json:"thread_id"`
	SenderID  uuid.UUID  `json:"sender_id"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

// CreateThreadRequest is the request body for creating a thread
type CreateThreadRequest struct {
	ParticipantID uuid.UUID `json:"participant_id"`
}

// SendMessageRequest is the request body for sending a message
type SendMessageRequest struct {
	Body string `json:"body"`
}

// ThreadMessagesResponse is the paginated response for thread messages
type ThreadMessagesResponse struct {
	Messages []MessageResponse `json:"messages"`
	HasMore  bool              `json:"has_more"`
	Cursor   *string           `json:"cursor,omitempty"`
}

// UnreadCountResponse is the response for total unread count
type UnreadCountResponse struct {
	UnreadCount int `json:"unread_count"`
}
