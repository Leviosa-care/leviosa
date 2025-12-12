package twilio

import (
	"context"
	"fmt"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"

	"github.com/Leviosa-care/leviosa/backend/internal/notification/domain"
	"github.com/Leviosa-care/leviosa/backend/internal/notification/ports"
)

type TwilioClient struct {
	client *twilio.RestClient
	from   string
}

func NewTwilioClient(accountSid, authToken, fromNumber string) ports.SMSService {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &TwilioClient{
		client: client,
		from:   fromNumber,
	}
}

func (c *TwilioClient) SendOTP(ctx context.Context, req domain.OTPSMSRequest) error {
	// Validate and normalize phone number
	smsReq, err := domain.NewSMSRequest(req.PhoneNumber, fmt.Sprintf("Your OTP code is: %s", req.OTP))
	if err != nil {
		return fmt.Errorf("invalid phone number for OTP SMS: %w", err)
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(smsReq.Phone)
	params.SetFrom(c.from)
	params.SetBody(smsReq.Message)

	_, err = c.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send OTP SMS via Twilio: %w", err)
	}

	return nil
}

func (c *TwilioClient) SendSMS(ctx context.Context, req domain.GenericSMSRequest) error {
	// Validate and normalize phone number
	smsReq, err := domain.NewSMSRequest(req.PhoneNumber, req.Message)
	if err != nil {
		return fmt.Errorf("invalid phone number for SMS: %w", err)
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(smsReq.Phone)
	params.SetFrom(c.from)
	params.SetBody(smsReq.Message)

	_, err = c.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS via Twilio: %w", err)
	}

	return nil
}
