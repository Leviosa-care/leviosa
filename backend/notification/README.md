# Notification Service

A standalone microservice for handling email and SMS notifications following hexagonal architecture.

## Architecture

This service follows the hexagonal architecture pattern with clear separation of concerns:

- **Domain**: Business entities and value objects (`internal/domain/`)
- **Ports**: Interface definitions (`internal/ports/`)
- **Application**: Use cases and business workflows (`internal/application/`)
- **Adapters**: Infrastructure implementations (`internal/adapters/`)

## Features

### Email Notifications
- OTP verification emails
- Welcome emails
- Password reset emails
- Event notifications
- Payment confirmations
- Vote notifications
- Registration reminders

### SMS Notifications
- SMS sending via Twilio

### Settings Synchronization
- Real-time company settings updates via RabbitMQ
- Cached company information (email, logo, address, Instagram)

## Dependencies

- **SMTP**: Gmail SMTP for email delivery
- **SMS**: Twilio for SMS delivery
- **Message Queue**: RabbitMQ for settings synchronization
- **Templates**: HTML email templates with embedded assets

## Environment Variables

Required environment variables:

```bash
GMAIL_EMAIL=your-gmail@example.com
GMAIL_PASSWORD=your-app-password
TWILIO_ACCOUNT_SID=your-twilio-account-sid
TWILIO_AUTH_TOKEN=your-twilio-auth-token
TWILIO_PHONE_NUMBER=your-twilio-phone
RABBITMQ_URL=amqp://localhost:5672
```

## Usage

Run the service:

```bash
go run main.go
```

The service will:
1. Initialize SMTP and Twilio clients
2. Load company settings from cache
3. Start RabbitMQ consumer for settings updates
4. Provide notification capabilities via the NotificationService interface

## Templates

Email templates are located in `templates/` directory and are embedded at compile time. Available templates:

- `otp.html` - OTP verification
- `welcome.html` - Welcome message
- `verify_email.html` - Password reset
- `event_notification.html` - Event notifications
- `payment.html` - Payment confirmations
- `new_vote.html` - Vote notifications

## Error Handling

The service uses the project's standardized error handling with sentinel errors from `core/errs`:

- `ErrInvalidValueErr` - Invalid input data
- `ErrConnectionFailureErr` - SMTP/RabbitMQ connection issues
- `ErrExternalServiceErr` - Twilio API errors
- `ErrInternalErr` - Internal processing errors