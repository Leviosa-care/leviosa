-- +goose Up
-- +goose StatementBegin

-- Create messaging schema
CREATE SCHEMA IF NOT EXISTS messaging;

-- Threads: one conversation per participant pair
CREATE TABLE messaging.threads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Thread participants: exactly 2 users per thread
CREATE TABLE messaging.thread_participants (
    thread_id UUID NOT NULL REFERENCES messaging.threads(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    PRIMARY KEY (thread_id, user_id)
);

-- Messages: individual messages within a thread (body encrypted via encx)
CREATE TABLE messaging.messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    thread_id UUID NOT NULL REFERENCES messaging.threads(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL,
    body_encrypted BYTEA NOT NULL,
    dek_encrypted BYTEA NOT NULL,
    key_version INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    read_at TIMESTAMPTZ
);

-- Indexes for performance
CREATE INDEX idx_messages_thread_created ON messaging.messages(thread_id, created_at DESC, id DESC);
CREATE INDEX idx_messages_thread_unread ON messaging.messages(thread_id, read_at) WHERE read_at IS NULL;
CREATE INDEX idx_thread_participants_user ON messaging.thread_participants(user_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP SCHEMA IF EXISTS messaging CASCADE;

-- +goose StatementEnd
