package messageService

import "time"

type Message struct {
	ID               string
	ConversationID   string
	SenderID         string
	Content          string `encx:"encrypt"`
	ContentEncrypted []byte
	CreatedAt        time.Time
	DEK              []byte `json:"-"`
	DEKEncrypted     []byte
	KeyVersion       int `json:"-"`
}
