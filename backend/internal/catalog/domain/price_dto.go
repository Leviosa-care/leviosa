package domain

// TODO: add all the price related struct with Input in their name

// CreatePriceRequest represents the data required to create a new price.
// This is used for inbound API requests (from client to handler/service).
type CreatePriceRequest struct {
	Amount   int               `json:"amount"`             // Price in cents
	Currency string            `json:"currency"`           // e.g., "usd", "eur"
	Interval string            `json:"interval"`           // "month", "year", "one_time"
	Nickname *string           `json:"nickname,omitempty"` // Optional nickname for the price
	Metadata map[string]string `json:"metadata,omitempty"` // Optional metadata to attach
}

// UpdatePriceRequest represents the fields that can be updated for an existing price.
// This is used for inbound API requests (from client to handler/service) for PATCH operations.
type UpdatePriceRequest struct {
	Active   *bool             `json:"active,omitempty"`   // Pointer for optional update (true/false)
	Nickname *string           `json:"nickname,omitempty"` // Pointer for optional update
	Metadata map[string]string `json:"metadata,omitempty"` // Full map replacement if provided
}
