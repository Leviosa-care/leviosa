package stripe

// NOTE: Payment intent status and refund reason constants have been moved to
// internal/booking/ports/payment_service.go to properly follow hexagonal architecture.
//
// The application layer should use ports.PaymentIntentStatus* constants.
// This infrastructure adapter uses the Stripe SDK types directly and converts
// them to the ports.PaymentIntentInfo type.
//
// Reference: https://stripe.com/docs/api/payment_intents/object#payment_intent_object-status
// Reference: https://stripe.com/docs/api/refunds/create#create_refund-reason

