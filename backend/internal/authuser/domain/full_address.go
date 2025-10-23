package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Response structure from the adresse.data.gouv.fr API
type AdresseResult struct {
	Features []struct {
		Properties struct {
			Label    string  `json:"label"`
			Score    float64 `json:"score"`
			City     string  `json:"city"`
			Postcode string  `json:"postcode"`
		} `json:"properties"`
	} `json:"features"`
}

func validateFullAddress(city, postalCode, address1, address2 string) error {
	fullAddress := fmt.Sprintf("%s %s %s", address1, address2, postalCode+" "+city)
	query := url.QueryEscape(fullAddress)
	apiURL := fmt.Sprintf("https://api-adresse.data.gouv.fr/search/?q=%s&limit=1", query)

	resp, err := http.Get(apiURL)
	if err != nil {
		return fmt.Errorf("failed to call address validation API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("address validation API returned status: %d", resp.StatusCode)
	}

	var result AdresseResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	if len(result.Features) == 0 {
		return errors.New("address not found")
	}

	props := result.Features[0].Properties

	// Basic confidence score check (tune this as needed)
	if props.Score < 0.7 {
		return fmt.Errorf("low confidence score for address: %.2f", props.Score)
	}

	// Check city and postal code match
	if !strings.EqualFold(props.City, city) {
		return fmt.Errorf("city mismatch: expected %s, got %s", city, props.City)
	}
	if props.Postcode != postalCode {
		return fmt.Errorf("postal code mismatch: expected %s, got %s", postalCode, props.Postcode)
	}

	return nil
}
