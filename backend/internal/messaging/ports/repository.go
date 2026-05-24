package ports

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/messaging/domain"
	"github.com/google/uuid"
)

// MessageRepository defines the interface for messaging data persistence.
type MessageRepository interface {
	// CreateThread creates a new thread with two participants atomically.
	CreateThread(ctx context.Context, thread *domain.Thread, participantA, participantB uuid.UUID) error

	// FindThreadByParticipants finds an existing thread shared by exactly these two users.
	FindThreadByParticipants(ctx context.Context, userA, userB uuid.UUID) (*domain.Thread, error)

	// GetThreadsForUser returns raw thread summaries for a user (body still encrypted).
	GetThreadsForUser(ctx context.Context, userID uuid.UUID) ([]domain.ThreadSummaryRaw, error)

	// GetThreadByID retrieves a thread, verifying the user is a participant.
	GetThreadByID(ctx context.Context, threadID, userID uuid.UUID) (*domain.Thread, error)

	// IsParticipant returns whether userID belongs to the given thread.
	IsParticipant(ctx context.Context, threadID, userID uuid.UUID) (bool, error)

	// CreateMessage persists an encrypted message.
	CreateMessage(ctx context.Context, message *domain.MessageEncx) error

	// GetMessagesByThread returns encrypted messages in descending creation order.
	// before + beforeID together form the composite cursor for stable pagination.
	GetMessagesByThread(ctx context.Context, threadID uuid.UUID, limit int, before *time.Time, beforeID *uuid.UUID) ([]domain.MessageEncx, error)

	// MarkThreadAsRead marks all unread messages sent by others as read.
	MarkThreadAsRead(ctx context.Context, threadID, userID uuid.UUID) error

	// GetUnreadCount returns the total number of unread messages for a user.
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}
