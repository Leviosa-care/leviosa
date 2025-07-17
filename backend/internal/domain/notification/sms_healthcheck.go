package notification

import (
	"context"
	"fmt"
)

func (c *smsClient) Healthcheck(ctx context.Context) error {
	_, err := c.FetchAccount(c.accountSid)
	if err != nil {
		return fmt.Errorf("twilio healthcheck failed: %w", err)
	}
	return nil
}
