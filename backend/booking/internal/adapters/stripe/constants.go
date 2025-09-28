package stripe

// Stripe Payment Intent Status Constants
// Reference: https://stripe.com/docs/api/payment_intents/object#payment_intent_object-status
const (
	// PaymentIntentStatusSucceeded indicates payment was successful
	PaymentIntentStatusSucceeded = "succeeded"

	// PaymentIntentStatusRequiresPaymentMethod indicates payment requires a payment method
	PaymentIntentStatusRequiresPaymentMethod = "requires_payment_method"

	// PaymentIntentStatusRequiresConfirmation indicates payment requires confirmation
	PaymentIntentStatusRequiresConfirmation = "requires_confirmation"

	// PaymentIntentStatusRequiresAction indicates payment requires additional action
	PaymentIntentStatusRequiresAction = "requires_action"

	// PaymentIntentStatusProcessing indicates payment is being processed
	PaymentIntentStatusProcessing = "processing"

	// PaymentIntentStatusCanceled indicates payment was canceled
	PaymentIntentStatusCanceled = "canceled"

	// PaymentIntentStatusPaymentFailed indicates payment failed
	PaymentIntentStatusPaymentFailed = "payment_failed"
)

// Stripe Refund Reason Constants
// Reference: https://stripe.com/docs/api/refunds/create#create_refund-reason
const (
	// RefundReasonDuplicate indicates duplicate charge
	RefundReasonDuplicate = "duplicate"

	// RefundReasonFraudulent indicates fraudulent charge
	RefundReasonFraudulent = "fraudulent"

	// RefundReasonRequestedByCustomer indicates customer requested refund
	RefundReasonRequestedByCustomer = "requested_by_customer"
)