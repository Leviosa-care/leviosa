package pricePayment_test

import (
	"context"
	"log"
	"os"
	"testing"

	pricePayment "github.com/Leviosa-care/leviosa/backend/internal/catalog/infrastructure/stripe/price"
	"github.com/Leviosa-care/leviosa/backend/internal/catalog/ports"

	tu "github.com/Leviosa-care/leviosa/backend/internal/common/testutils"
)

var (
	stripeContainer *tu.StripeMockContainer
	stripeService   ports.PricePaymentGateway
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var err error

	// Setup Stripe mock container
	log.Println("Setting up Stripe mock container...")
	stripeContainer, err = tu.SetupStripeMock(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to setup stripe mock container: %v", err)
	}
	defer tu.TeardownStripeMock(ctx, nil, stripeContainer)
	log.Printf("Stripe mock container started at %s", stripeContainer.URL)

	// Initialize Stripe service with mock URL
	stripeService = pricePayment.NewPrice("sk_test_123456789012345678901234", stripeContainer.URL)
	log.Println("Stripe price service initialized with mock container")

	// Run tests
	code := m.Run()

	log.Println("Tests completed, cleaning up...")
	os.Exit(code)
}