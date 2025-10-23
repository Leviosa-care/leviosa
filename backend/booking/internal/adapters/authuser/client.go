package authuser

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Leviosa-care/booking/internal/ports"
	"github.com/Leviosa-care/leviosa/backend/internal/common/errs"
	"github.com/google/uuid"
)

// Client implements the AuthUserClient interface for HTTP communication
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string // For service-to-service authentication
}

// NewClient creates a new authuser HTTP client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
	}
}

// GetPartnerByID retrieves a partner by their ID from the authuser service
func (c *Client) GetPartnerByID(ctx context.Context, partnerID uuid.UUID) (*ports.PartnerInfo, error) {
	url := fmt.Sprintf("%s/admin/partners/%s", c.baseURL, partnerID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add service-to-service authentication header
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.NewConnectionFailure("authuser service communication failed", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var partnerResp struct {
			ID             uuid.UUID `json:"id"`
			UserID         uuid.UUID `json:"user_id"`
			Bio            string    `json:"bio"`
			Experience     string    `json:"experience"`
			Certifications []string  `json:"certifications"`
			IsVerified     bool      `json:"is_verified"`
			User           struct {
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Telephone string `json:"telephone"`
			} `json:"user"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&partnerResp); err != nil {
			return nil, fmt.Errorf("failed to decode partner response: %w", err)
		}

		return &ports.PartnerInfo{
			ID:             partnerResp.ID,
			UserID:         partnerResp.UserID,
			Bio:            partnerResp.Bio,
			Experience:     partnerResp.Experience,
			Certifications: partnerResp.Certifications,
			IsVerified:     partnerResp.IsVerified,
			Email:          partnerResp.User.Email,
			FirstName:      partnerResp.User.FirstName,
			LastName:       partnerResp.User.LastName,
			Telephone:      partnerResp.User.Telephone,
		}, nil

	case http.StatusNotFound:
		return nil, errs.NewNotFoundErr("partner", partnerID.String())

	case http.StatusUnauthorized:
		return nil, errs.NewPermissionDenied("insufficient permissions to access partner data")

	case http.StatusInternalServerError:
		return nil, errs.NewConnectionFailure("authuser service internal error", nil)

	default:
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}
}

// GetPartnerByUserID retrieves a partner by their user ID from the authuser service
func (c *Client) GetPartnerByUserID(ctx context.Context, userID uuid.UUID) (*ports.PartnerInfo, error) {
	url := fmt.Sprintf("%s/admin/partners/user/%s", c.baseURL, userID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add service-to-service authentication header
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.NewConnectionFailure("authuser service communication failed", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var partnerResp struct {
			ID             uuid.UUID `json:"id"`
			UserID         uuid.UUID `json:"user_id"`
			Bio            string    `json:"bio"`
			Experience     string    `json:"experience"`
			Certifications []string  `json:"certifications"`
			IsVerified     bool      `json:"is_verified"`
			User           struct {
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Telephone string `json:"telephone"`
			} `json:"user"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&partnerResp); err != nil {
			return nil, fmt.Errorf("failed to decode partner response: %w", err)
		}

		return &ports.PartnerInfo{
			ID:             partnerResp.ID,
			UserID:         partnerResp.UserID,
			Bio:            partnerResp.Bio,
			Experience:     partnerResp.Experience,
			Certifications: partnerResp.Certifications,
			IsVerified:     partnerResp.IsVerified,
			Email:          partnerResp.User.Email,
			FirstName:      partnerResp.User.FirstName,
			LastName:       partnerResp.User.LastName,
			Telephone:      partnerResp.User.Telephone,
		}, nil

	case http.StatusNotFound:
		return nil, errs.NewNotFoundErr("partner", fmt.Sprintf("user_id:%s", userID.String()))

	case http.StatusUnauthorized:
		return nil, errs.NewPermissionDenied("insufficient permissions to access partner data")

	case http.StatusInternalServerError:
		return nil, errs.NewConnectionFailure("authuser service internal error", nil)

	default:
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}
}

// ValidatePartnerExists checks if a partner exists and is verified
func (c *Client) ValidatePartnerExists(ctx context.Context, partnerID uuid.UUID) (bool, error) {
	partner, err := c.GetPartnerByID(ctx, partnerID)
	if err != nil {
		if errs.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	// Return true only if partner exists and is verified
	return partner.IsVerified, nil
}

// GetPartnerSpecializations retrieves a partner's specializations
func (c *Client) GetPartnerSpecializations(ctx context.Context, partnerID uuid.UUID) ([]ports.SpecializationInfo, error) {
	url := fmt.Sprintf("%s/partners/%s/specializations", c.baseURL, partnerID.String())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add service-to-service authentication header
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.NewConnectionFailure("authuser service communication failed", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var specializationsResp struct {
			Specializations []struct {
				ID          uuid.UUID `json:"id"`
				Name        string    `json:"name"`
				Description string    `json:"description"`
				IsActive    bool      `json:"is_active"`
			} `json:"specializations"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&specializationsResp); err != nil {
			return nil, fmt.Errorf("failed to decode specializations response: %w", err)
		}

		specializations := make([]ports.SpecializationInfo, len(specializationsResp.Specializations))
		for i, spec := range specializationsResp.Specializations {
			specializations[i] = ports.SpecializationInfo{
				ID:          spec.ID,
				Name:        spec.Name,
				Description: spec.Description,
				IsActive:    spec.IsActive,
			}
		}

		return specializations, nil

	case http.StatusNotFound:
		return nil, errs.NewNotFoundErr("partner", partnerID.String())

	case http.StatusUnauthorized:
		return nil, errs.NewPermissionDenied("insufficient permissions to access partner specializations")

	case http.StatusInternalServerError:
		return nil, errs.NewConnectionFailure("authuser service internal error", nil)

	default:
		return nil, fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}
}

