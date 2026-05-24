package ports

import (
	"context"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/messaging/domain"
	"github.com/google/uuid"
)

// MessagingService defines the business logic interface for messaging.
type MessagingService interface {
	ListThreads(ctx context.Context, userID uuid.UUID) ([]domain.ThreadSummary, error)
	CreateThread(ctx context.Context, currentUserID, participantID uuid.UUID, currentUserRole identity.Role) (*domain.Thread, error)
	GetMessages(ctx context.Context, threadID, userID uuid.UUID, limit int, cursor string) (*domain.ThreadMessagesResponse, error)
	SendMessage(ctx context.Context, threadID, userID uuid.UUID, body string) (*domain.Message, error)
	MarkAsRead(ctx context.Context, threadID, userID uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}

// BookingChecker verifies whether a partner has an existing booking with a given client.
// Defined here so the messaging domain does not import the booking domain directly.
type BookingChecker interface {
	HasBookingRelationship(ctx context.Context, partnerID, clientID uuid.UUID) (bool, error)
}

// UserNameFetcher resolves a user ID to a display name.
// Defined here so the messaging domain does not import the authuser domain directly.
type UserNameFetcher interface {
	FetchName(ctx context.Context, userID uuid.UUID) (string, error)
}
