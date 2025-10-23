package twilio

import (
	"context"
	"fmt"

	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/Leviosa-care/notification/internal/domain"
	"github.com/Leviosa-care/notification/internal/ports"

	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioClient struct {
	client     *openapi.ApiService
	sender     string
	accountSid string
}

func New(accountSid, authToken, senderPhone string) (ports.SMSService, error) {
	if accountSid == "" {
		return nil, errs.NewInvalidValueErr(fmt.Errorf("Twilio account SID cannot be empty"))
	}

	if authToken == "" {
		return nil, errs.NewInvalidValueErr(fmt.Errorf("Twilio auth token cannot be empty"))
	}

	if senderPhone == "" {
		return nil, errs.NewInvalidValueErr(fmt.Errorf("Twilio sender phone cannot be empty"))
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})

	return &TwilioClient{
		client:     client.Api,
		sender:     senderPhone,
		accountSid: accountSid,
	}, nil
}

func (c *TwilioClient) SendSMS(ctx context.Context, request *domain.SMSRequest) error {
	if request == nil {
		return errs.NewInvalidValueErr(fmt.Errorf("SMS request cannot be nil"))
	}

	if request.Phone == "" {
		return errs.NewInvalidValueErr(fmt.Errorf("recipient phone number cannot be empty"))
	}

	if request.Message == "" {
		return errs.NewInvalidValueErr(fmt.Errorf("SMS message cannot be empty"))
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(request.Phone)
	params.SetFrom(c.sender)
	params.SetBody(request.Message)

	_, err := c.client.CreateMessage(params)
	if err != nil {
		return errs.NewExternalServiceErr(fmt.Errorf("send SMS via Twilio: %w", err))
	}

	return nil
}

