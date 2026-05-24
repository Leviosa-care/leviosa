package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrThreadNotFound        = errors.New("thread not found")
	ErrNotThreadParticipant  = errors.New("user is not a participant of this thread")
	ErrEmptyMessageBody      = errors.New("message body cannot be empty")
	ErrCannotInitiateThread  = errors.New("this user role cannot initiate threads")
	ErrThreadAlreadyExists   = errors.New("a thread between these users already exists")
	ErrInvalidParticipantID  = errors.New("invalid participant ID")
	ErrNoBookingRelationship = errors.New("partner has no booking with this user")
)

// Thread represents a conversation between two users
type Thread struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

// ThreadParticipant links a user to a thread
type ThreadParticipant struct {
	ThreadID uuid.UUID `json:"thread_id"`
	UserID   uuid.UUID `json:"user_id"`
}

// Message represents a single message in a thread
type Message struct {
	ID        uuid.UUID  `json:"id"`
	ThreadID  uuid.UUID  `json:"thread_id"`
	SenderID  uuid.UUID  `json:"sender_id"`
	Body      string     `json:"body"`
	CreatedAt time.Time  `json:"created_at"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
}

// NewThread creates a new thread between two participants
func NewThread(participantA, participantB uuid.UUID) (*Thread, error) {
	if participantA == uuid.Nil || participantB == uuid.Nil {
		return nil, ErrInvalidParticipantID
	}
	if participantA == participantB {
		return nil, ErrInvalidParticipantID
	}

	return &Thread{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
	}, nil
}

// NewMessage creates a new message in a thread
func NewMessage(threadID, senderID uuid.UUID, body string) (*Message, error) {
	if threadID == uuid.Nil {
		return nil, ErrThreadNotFound
	}
	if senderID == uuid.Nil {
		return nil, ErrInvalidParticipantID
	}
	if body == "" {
		return nil, ErrEmptyMessageBody
	}

	return &Message{
		ID:        uuid.New(),
		ThreadID:  threadID,
		SenderID:  senderID,
		Body:      body,
		CreatedAt: time.Now(),
	}, nil
}
