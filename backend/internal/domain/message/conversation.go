package messageService

import "time"

type Conversation struct {
	ID        string
	UserID    string
	PartnerID string
	CreatedAt time.Time
}
