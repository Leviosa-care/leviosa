package authuser

import (
	"context"
	"fmt"
	"strings"

	authuserPorts "github.com/Leviosa-care/leviosa/backend/internal/authuser/ports"
	bookingPorts "github.com/Leviosa-care/leviosa/backend/internal/booking/ports"

	"github.com/google/uuid"
)

// InProcessClient is an in-process implementation of AuthUserClient that directly
// delegates to the authuser service's PublicPartnerService interface.
//
// This implementation is used in the modular monolith architecture for efficient
// in-process communication without HTTP overhead. When migrating to microservices,
// this can be replaced with an HTTP-based implementation without changing the
// interface or business logic.
type InProcessClient struct {
	partnerService authuserPorts.PublicPartnerService
	userService    authuserPorts.UserService
}

// NewInProcessClient creates a new in-process AuthUserClient implementation.
func NewInProcessClient(partnerService authuserPorts.PublicPartnerService, userService authuserPorts.UserService) bookingPorts.AuthUserClient {
	return &InProcessClient{
		partnerService: partnerService,
		userService:    userService,
	}
}

// GetPartnerVerificationStatus checks if a partner is verified by delegating
// to the authuser service's PublicPartnerService.
func (c *InProcessClient) GetPartnerVerificationStatus(ctx context.Context, partnerID uuid.UUID) (bool, error) {
	return c.partnerService.GetPartnerVerificationStatus(ctx, partnerID)
}

// GetPartnerByUserID retrieves partner information by user ID and converts
// the response to the booking module's PartnerInfo format.
func (c *InProcessClient) GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*bookingPorts.PartnerInfo, error) {
	partnerResponse, err := c.partnerService.GetPartnerByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check verification status
	isVerified, err := c.partnerService.GetPartnerVerificationStatus(ctx, partnerResponse.ID)
	if err != nil {
		return nil, err
	}

	return &bookingPorts.PartnerInfo{
		ID:         partnerResponse.ID,
		UserID:     partnerResponse.UserID,
		IsVerified: isVerified,
	}, nil
}

// GetUserName retrieves the display name for a user by their ID.
func (c *InProcessClient) GetUserName(ctx context.Context, userID uuid.UUID) (string, error) {
	user, err := c.userService.GetUserByID(ctx, userID)
	if err != nil {
		return "", fmt.Errorf("get user by id %s: %w", userID, err)
	}
	name := strings.TrimSpace(fmt.Sprintf("%s %s", user.FirstName, user.LastName))
	return name, nil
}
