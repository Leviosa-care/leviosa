package messaging

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/identity"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/messaging/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/messaging/ports"
	"github.com/google/uuid"
	"github.com/hengadev/encx"
)

type service struct {
	repo        ports.MessageRepository
	crypto      encx.CryptoService
	bookChecker ports.BookingChecker
	nameFetcher ports.UserNameFetcher
	broker      Broker
}

// Broker publishes SSE events when a new message is created.
type Broker interface {
	Publish(threadID uuid.UUID, msg *domain.MessageResponse)
}

func New(
	repo ports.MessageRepository,
	crypto encx.CryptoService,
	bookChecker ports.BookingChecker,
	nameFetcher ports.UserNameFetcher,
	broker Broker,
) ports.MessagingService {
	return &service{
		repo:        repo,
		crypto:      crypto,
		bookChecker: bookChecker,
		nameFetcher: nameFetcher,
		broker:      broker,
	}
}

func (s *service) ListThreads(ctx context.Context, userID uuid.UUID) ([]domain.ThreadSummary, error) {
	raws, err := s.repo.GetThreadsForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list threads: %w", err)
	}

	summaries := make([]domain.ThreadSummary, len(raws))
	for i, raw := range raws {
		name, _ := s.nameFetcher.FetchName(ctx, raw.ParticipantID)

		var lastMessage string
		var lastMessageAt time.Time

		if raw.LastMessageAt != nil {
			lastMessageAt = *raw.LastMessageAt
		}

		if len(raw.LastBodyEncrypted) > 0 {
			msgEncx := &domain.MessageEncx{
				BodyEncrypted: raw.LastBodyEncrypted,
				DEKEncrypted:  raw.LastDEKEncrypted,
				KeyVersion:    raw.LastKeyVersion,
			}
			msg, decErr := domain.DecryptMessageEncx(ctx, s.crypto, msgEncx)
			if decErr == nil {
				lastMessage = msg.Body
			}
		}

		summaries[i] = domain.ThreadSummary{
			ThreadID:        raw.ThreadID,
			ParticipantID:   raw.ParticipantID,
			ParticipantName: name,
			LastMessage:     lastMessage,
			LastMessageAt:   lastMessageAt,
			UnreadCount:     raw.UnreadCount,
		}
	}

	return summaries, nil
}

func (s *service) CreateThread(ctx context.Context, currentUserID, participantID uuid.UUID, currentUserRole identity.Role) (*domain.Thread, error) {
	if currentUserID == uuid.Nil || participantID == uuid.Nil {
		return nil, errs.NewInvalidInputErr(fmt.Errorf("invalid user IDs"))
	}

	switch currentUserRole {
	case identity.Administrator:
		// allowed unconditionally
	case identity.Partner:
		ok, err := s.bookChecker.HasBookingRelationship(ctx, currentUserID, participantID)
		if err != nil {
			return nil, fmt.Errorf("check booking relationship: %w", err)
		}
		if !ok {
			return nil, domain.ErrNoBookingRelationship
		}
	default:
		return nil, domain.ErrCannotInitiateThread
	}

	existing, err := s.repo.FindThreadByParticipants(ctx, currentUserID, participantID)
	if err != nil && !errors.Is(err, errs.ErrRepositoryNotFound) {
		return nil, fmt.Errorf("check existing thread: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrThreadAlreadyExists
	}

	thread, err := domain.NewThread(currentUserID, participantID)
	if err != nil {
		return nil, fmt.Errorf("create thread domain object: %w", err)
	}

	if err := s.repo.CreateThread(ctx, thread, currentUserID, participantID); err != nil {
		return nil, fmt.Errorf("persist thread: %w", err)
	}

	return thread, nil
}

func (s *service) GetMessages(ctx context.Context, threadID, userID uuid.UUID, limit int, cursor string) (*domain.ThreadMessagesResponse, error) {
	isParticipant, err := s.repo.IsParticipant(ctx, threadID, userID)
	if err != nil {
		return nil, fmt.Errorf("check participation: %w", err)
	}
	if !isParticipant {
		return nil, domain.ErrNotThreadParticipant
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}

	var before *time.Time
	var beforeID *uuid.UUID
	if cursor != "" {
		parts := strings.SplitN(cursor, "|", 2)
		if len(parts) == 2 {
			t, err1 := time.Parse(time.RFC3339Nano, parts[0])
			id, err2 := uuid.Parse(parts[1])
			if err1 == nil && err2 == nil {
				before = &t
				beforeID = &id
			}
		}
	}

	encxMessages, err := s.repo.GetMessagesByThread(ctx, threadID, limit+1, before, beforeID)
	if err != nil {
		return nil, fmt.Errorf("get messages: %w", err)
	}

	hasMore := len(encxMessages) > limit
	if hasMore {
		encxMessages = encxMessages[:limit]
	}

	response := &domain.ThreadMessagesResponse{
		Messages: make([]domain.MessageResponse, 0, len(encxMessages)),
		HasMore:  hasMore,
	}

	for _, encxMsg := range encxMessages {
		msg, decErr := domain.DecryptMessageEncx(ctx, s.crypto, &encxMsg)
		if decErr != nil {
			return nil, fmt.Errorf("decrypt message %s: %w", encxMsg.ID, decErr)
		}
		response.Messages = append(response.Messages, domain.MessageResponse{
			ID:        msg.ID,
			ThreadID:  msg.ThreadID,
			SenderID:  msg.SenderID,
			Body:      msg.Body,
			CreatedAt: msg.CreatedAt,
			ReadAt:    msg.ReadAt,
		})
	}

	if hasMore && len(encxMessages) > 0 {
		last := encxMessages[len(encxMessages)-1]
		cursorStr := last.CreatedAt.Format(time.RFC3339Nano) + "|" + last.ID.String()
		response.Cursor = &cursorStr
	}

	return response, nil
}

func (s *service) SendMessage(ctx context.Context, threadID, userID uuid.UUID, body string) (*domain.MessageResponse, error) {
	isParticipant, err := s.repo.IsParticipant(ctx, threadID, userID)
	if err != nil {
		return nil, fmt.Errorf("check participation: %w", err)
	}
	if !isParticipant {
		return nil, domain.ErrNotThreadParticipant
	}

	message, err := domain.NewMessage(threadID, userID, body)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	msgEncx, err := domain.ProcessMessageEncx(ctx, s.crypto, message)
	if err != nil {
		return nil, fmt.Errorf("encrypt message: %w", err)
	}

	if err := s.repo.CreateMessage(ctx, msgEncx); err != nil {
		return nil, fmt.Errorf("persist message: %w", err)
	}

	response := &domain.MessageResponse{
		ID:        message.ID,
		ThreadID:  message.ThreadID,
		SenderID:  message.SenderID,
		Body:      message.Body,
		CreatedAt: message.CreatedAt,
		ReadAt:    message.ReadAt,
	}

	if s.broker != nil {
		s.broker.Publish(threadID, response)
	}

	return response, nil
}

func (s *service) MarkAsRead(ctx context.Context, threadID, userID uuid.UUID) error {
	isParticipant, err := s.repo.IsParticipant(ctx, threadID, userID)
	if err != nil {
		return fmt.Errorf("check participation: %w", err)
	}
	if !isParticipant {
		return domain.ErrNotThreadParticipant
	}
	return s.repo.MarkThreadAsRead(ctx, threadID, userID)
}

func (s *service) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.repo.GetUnreadCount(ctx, userID)
}
