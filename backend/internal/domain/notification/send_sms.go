package notification

import (
	"context"
	"fmt"

	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func (c *smsClient) SendSMS(ctx context.Context, phone, message string) error {
	params := &openapi.CreateMessageParams{}
	params.SetTo(phone)      // Destination phone number : Gary
	params.SetFrom(c.sender) // Your Twilio phone number
	params.SetBody(message)
	resp, err := c.CreateMessage(params)
	_ = resp
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	if resp != nil && resp.Sid != nil {
		fmt.Printf("SMS sent with SID: %s\n", *resp.Sid)
	}
	return nil
}
