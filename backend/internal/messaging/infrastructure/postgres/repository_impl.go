package postgres

import (
	"context"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/leviosa/backend/internal/messaging/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *repository) CreateThread(ctx context.Context, thread *domain.Thread, participantA, participantB uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return errs.NewDatabaseErr(err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO messaging.threads (id, created_at) VALUES ($1, $2)`,
		thread.ID, thread.CreatedAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create thread", err)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO messaging.thread_participants (thread_id, user_id) VALUES ($1, $2), ($1, $3)`,
		thread.ID, participantA, participantB,
	)
	if err != nil {
		return errs.ClassifyPgError("create thread participants", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return errs.NewDatabaseErr(err)
	}
	return nil
}

func (r *repository) FindThreadByParticipants(ctx context.Context, userA, userB uuid.UUID) (*domain.Thread, error) {
	query := `
		SELECT t.id, t.created_at
		FROM messaging.threads t
		WHERE EXISTS (
			SELECT 1 FROM messaging.thread_participants tp
			WHERE tp.thread_id = t.id AND tp.user_id = $1
		)
		AND EXISTS (
			SELECT 1 FROM messaging.thread_participants tp
			WHERE tp.thread_id = t.id AND tp.user_id = $2
		)
		LIMIT 1
	`

	thread := &domain.Thread{}
	err := r.pool.QueryRow(ctx, query, userA, userB).Scan(&thread.ID, &thread.CreatedAt)
	if err != nil {
		return nil, errs.ClassifyPgError("find thread by participants", err)
	}
	return thread, nil
}

func (r *repository) GetThreadsForUser(ctx context.Context, userID uuid.UUID) ([]domain.ThreadSummaryRaw, error) {
	// Participant names are encrypted in auth.users; name resolution happens in the service layer.
	query := `
		SELECT
			t.id,
			tp2.user_id AS participant_id,
			(SELECT m.body_encrypted  FROM messaging.messages m WHERE m.thread_id = t.id ORDER BY m.created_at DESC, m.id DESC LIMIT 1) AS last_body_encrypted,
			(SELECT m.dek_encrypted   FROM messaging.messages m WHERE m.thread_id = t.id ORDER BY m.created_at DESC, m.id DESC LIMIT 1) AS last_dek_encrypted,
			(SELECT m.key_version     FROM messaging.messages m WHERE m.thread_id = t.id ORDER BY m.created_at DESC, m.id DESC LIMIT 1) AS last_key_version,
			(SELECT m.created_at      FROM messaging.messages m WHERE m.thread_id = t.id ORDER BY m.created_at DESC, m.id DESC LIMIT 1) AS last_message_at,
			(SELECT COUNT(*) FROM messaging.messages m
			 WHERE m.thread_id = t.id AND m.sender_id != $1 AND m.read_at IS NULL) AS unread_count
		FROM messaging.threads t
		JOIN messaging.thread_participants tp1 ON tp1.thread_id = t.id AND tp1.user_id = $1
		JOIN messaging.thread_participants tp2 ON tp2.thread_id = t.id AND tp2.user_id != $1
		ORDER BY last_message_at DESC NULLS LAST
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, errs.ClassifyPgError("get threads for user", err)
	}
	defer rows.Close()

	var threads []domain.ThreadSummaryRaw
	for rows.Next() {
		var raw domain.ThreadSummaryRaw
		var lastBodyEncrypted []byte
		var lastDEKEncrypted []byte
		var lastKeyVersion *int32
		var lastMessageAt *time.Time

		if err := rows.Scan(
			&raw.ThreadID,
			&raw.ParticipantID,
			&lastBodyEncrypted,
			&lastDEKEncrypted,
			&lastKeyVersion,
			&lastMessageAt,
			&raw.UnreadCount,
		); err != nil {
			return nil, errs.ClassifyPgError("scan thread summary", err)
		}

		raw.LastBodyEncrypted = lastBodyEncrypted
		raw.LastDEKEncrypted = lastDEKEncrypted
		if lastKeyVersion != nil {
			raw.LastKeyVersion = int(*lastKeyVersion)
		}
		raw.LastMessageAt = lastMessageAt

		threads = append(threads, raw)
	}

	if threads == nil {
		threads = []domain.ThreadSummaryRaw{}
	}
	return threads, nil
}

func (r *repository) GetThreadByID(ctx context.Context, threadID, userID uuid.UUID) (*domain.Thread, error) {
	query := `
		SELECT t.id, t.created_at
		FROM messaging.threads t
		JOIN messaging.thread_participants tp ON tp.thread_id = t.id AND tp.user_id = $2
		WHERE t.id = $1
	`

	thread := &domain.Thread{}
	err := r.pool.QueryRow(ctx, query, threadID, userID).Scan(&thread.ID, &thread.CreatedAt)
	if err != nil {
		return nil, errs.ClassifyPgError("get thread by id", err)
	}
	return thread, nil
}

func (r *repository) IsParticipant(ctx context.Context, threadID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM messaging.thread_participants WHERE thread_id = $1 AND user_id = $2)`,
		threadID, userID,
	).Scan(&exists)
	if err != nil {
		return false, errs.ClassifyPgError("check participant", err)
	}
	return exists, nil
}

func (r *repository) CreateMessage(ctx context.Context, message *domain.MessageEncx) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO messaging.messages
			(id, thread_id, sender_id, body_encrypted, dek_encrypted, key_version, created_at, read_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		message.ID,
		message.ThreadID,
		message.SenderID,
		message.BodyEncrypted,
		message.DEKEncrypted,
		message.KeyVersion,
		message.CreatedAt,
		message.ReadAt,
	)
	if err != nil {
		return errs.ClassifyPgError("create message", err)
	}
	return nil
}

func (r *repository) GetMessagesByThread(ctx context.Context, threadID uuid.UUID, limit int, before *time.Time, beforeID *uuid.UUID) ([]domain.MessageEncx, error) {
	var rows pgx.Rows
	var err error

	if before != nil && beforeID != nil {
		rows, err = r.pool.Query(ctx, `
			SELECT id, thread_id, sender_id, body_encrypted, dek_encrypted, key_version, created_at, read_at
			FROM messaging.messages
			WHERE thread_id = $1
			  AND (created_at < $2 OR (created_at = $2 AND id < $3))
			ORDER BY created_at DESC, id DESC
			LIMIT $4
		`, threadID, *before, *beforeID, limit)
	} else {
		rows, err = r.pool.Query(ctx, `
			SELECT id, thread_id, sender_id, body_encrypted, dek_encrypted, key_version, created_at, read_at
			FROM messaging.messages
			WHERE thread_id = $1
			ORDER BY created_at DESC, id DESC
			LIMIT $2
		`, threadID, limit)
	}

	if err != nil {
		return nil, errs.ClassifyPgError("get messages by thread", err)
	}
	defer rows.Close()

	var messages []domain.MessageEncx
	for rows.Next() {
		var m domain.MessageEncx
		if err := rows.Scan(
			&m.ID, &m.ThreadID, &m.SenderID,
			&m.BodyEncrypted, &m.DEKEncrypted, &m.KeyVersion,
			&m.CreatedAt, &m.ReadAt,
		); err != nil {
			return nil, errs.ClassifyPgError("scan message", err)
		}
		messages = append(messages, m)
	}

	if messages == nil {
		messages = []domain.MessageEncx{}
	}
	return messages, nil
}

func (r *repository) MarkThreadAsRead(ctx context.Context, threadID, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE messaging.messages SET read_at = NOW()
		 WHERE thread_id = $1 AND sender_id != $2 AND read_at IS NULL`,
		threadID, userID,
	)
	if err != nil {
		return errs.ClassifyPgError("mark thread as read", err)
	}
	return nil
}

func (r *repository) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM messaging.messages m
		JOIN messaging.thread_participants tp ON tp.thread_id = m.thread_id AND tp.user_id = $1
		WHERE m.sender_id != $1 AND m.read_at IS NULL
	`, userID).Scan(&count)
	if err != nil {
		return 0, errs.ClassifyPgError("get unread count", err)
	}
	return count, nil
}
