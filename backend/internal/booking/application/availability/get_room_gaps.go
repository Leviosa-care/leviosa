package availability

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/booking/domain"
	catalogDomain "github.com/Leviosa-care/leviosa/backend/internal/catalog/domain"
)

// GetRoomGaps finds time gaps in a room's schedule for a specific date
// and suggests products that could fit in those gaps
func (s *AvailabilityService) GetRoomGaps(ctx context.Context, request domain.GetRoomGapsRequest) (*domain.GetRoomGapsResponse, error) {
	// 1. Get room schedule for the specific date
	roomHours, err := s.roomScheduleRepo.GetRoomHoursForDate(ctx, request.RoomID, request.Date)
	if err != nil {
		return nil, fmt.Errorf("get room hours: %w", err)
	}

	// 2. Get all bookings for the specified date
	bookings, err := s.availabilityRepo.GetRoomBookingsForDate(ctx, request.RoomID, request.Date)
	if err != nil {
		return nil, fmt.Errorf("get room bookings: %w", err)
	}

	// 3. Get products to suggest for gaps
	products, err := s.productService.GetAllPublishedProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("get products: %w", err)
	}

	// 4. Calculate operating hours for the specific date
	operatingStart := combineDateTime(request.Date, roomHours.OpenTime)
	operatingEnd := combineDateTime(request.Date, roomHours.CloseTime)

	// 5. Find gaps between bookings
	gaps := s.calculateGaps(bookings, operatingStart, operatingEnd, products)

	// 6. Calculate total gap time
	totalGapMinutes := 0
	for _, gap := range gaps {
		totalGapMinutes += gap.DurationMinutes
	}

	// 7. Build response
	response := &domain.GetRoomGapsResponse{
		RoomID: request.RoomID,
		Date:   request.Date,
		OperatingHours: domain.OperatingHoursResponse{
			StartTime: operatingStart,
			EndTime:   operatingEnd,
		},
		Gaps:            gaps,
		TotalGapMinutes: totalGapMinutes,
	}

	return response, nil
}

// calculateGaps finds time gaps between bookings and suggests fitting products.
//
// Algorithm:
//  1. If no bookings exist, return entire operating period as one gap
//  2. Check gap before first booking (if first booking starts after operating hours)
//  3. Iterate between consecutive bookings to find gaps
//  4. Check gap after last booking (if last booking ends before operating hours close)
//
// For each gap found, the function:
//   - Calculates duration in minutes
//   - Finds products that fit (product.Duration + product.BufferTime <= gap)
//   - Sorts suggestions by duration (shortest first for maximum flexibility)
func (s *AvailabilityService) calculateGaps(
	bookings []*domain.AvailabilityEncx,
	operatingStart, operatingEnd time.Time,
	products []*catalogDomain.ProductRes,
) []domain.TimeGapResponse {
	gaps := []domain.TimeGapResponse{}

	// If no bookings, the entire operating period is a gap
	if len(bookings) == 0 {
		gapMinutes := int(operatingEnd.Sub(operatingStart).Minutes())
		suggestions := findFittingProducts(gapMinutes, products)
		gaps = append(gaps, domain.TimeGapResponse{
			StartTime:         operatingStart,
			EndTime:           operatingEnd,
			DurationMinutes:   gapMinutes,
			IsBookable:        len(suggestions) > 0,
			SuggestedProducts: suggestions,
		})
		return gaps
	}

	// Check gap before first booking
	firstBooking := bookings[0]
	if firstBooking.StartTime.After(operatingStart) {
		gapMinutes := int(firstBooking.StartTime.Sub(operatingStart).Minutes())
		if gapMinutes > 0 {
			suggestions := findFittingProducts(gapMinutes, products)
			gaps = append(gaps, domain.TimeGapResponse{
				StartTime:         operatingStart,
				EndTime:           firstBooking.StartTime,
				DurationMinutes:   gapMinutes,
				IsBookable:        len(suggestions) > 0,
				SuggestedProducts: suggestions,
			})
		}
	}

	// Check gaps between consecutive bookings
	for i := 0; i < len(bookings)-1; i++ {
		currentEnd := bookings[i].EndTime
		nextStart := bookings[i+1].StartTime

		if nextStart.After(currentEnd) {
			gapMinutes := int(nextStart.Sub(currentEnd).Minutes())
			if gapMinutes > 0 {
				suggestions := findFittingProducts(gapMinutes, products)
				gaps = append(gaps, domain.TimeGapResponse{
					StartTime:         currentEnd,
					EndTime:           nextStart,
					DurationMinutes:   gapMinutes,
					IsBookable:        len(suggestions) > 0,
					SuggestedProducts: suggestions,
				})
			}
		}
	}

	// Check gap after last booking
	lastBooking := bookings[len(bookings)-1]
	if lastBooking.EndTime.Before(operatingEnd) {
		gapMinutes := int(operatingEnd.Sub(lastBooking.EndTime).Minutes())
		if gapMinutes > 0 {
			suggestions := findFittingProducts(gapMinutes, products)
			gaps = append(gaps, domain.TimeGapResponse{
				StartTime:         lastBooking.EndTime,
				EndTime:           operatingEnd,
				DurationMinutes:   gapMinutes,
				IsBookable:        len(suggestions) > 0,
				SuggestedProducts: suggestions,
			})
		}
	}

	return gaps
}

// findFittingProducts finds all products that can be scheduled within a time gap.
//
// A product fits if its total required time (session + buffer) does not exceed the gap:
//
//	gap_duration >= product.Duration + product.BufferTime
//
// Returns suggestions sorted by total time (shortest first) to maximize scheduling flexibility.
// This allows practitioners to see products that leave the most room for additional bookings.
func findFittingProducts(gapMinutes int, products []*catalogDomain.ProductRes) []domain.ProductSuggestion {
	suggestions := []domain.ProductSuggestion{}

	for _, product := range products {
		totalTime := product.Duration + product.BufferTime
		if gapMinutes >= totalTime {
			suggestions = append(suggestions, domain.ProductSuggestion{
				ProductID:   product.ID,
				ProductName: product.Name,
				Duration:    product.Duration,
				BufferTime:  product.BufferTime,
				TotalTime:   totalTime,
			})
		}
	}

	// Sort suggestions by total time (shortest first for flexibility)
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].TotalTime < suggestions[j].TotalTime
	})

	return suggestions
}

// combineDateTime combines a date with a time-of-day to create a full timestamp.
//
// Takes the year/month/day from 'date' and hour/minute/second from 'timeOfDay',
// preserving the location/timezone from 'date'.
//
// Example:
//
//	date = 2025-11-24 00:00:00 UTC
//	timeOfDay = 0001-01-01 09:30:00 UTC (from Room.OperatingStartTime)
//	result = 2025-11-24 09:30:00 UTC
func combineDateTime(date time.Time, timeOfDay time.Time) time.Time {
	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		timeOfDay.Hour(),
		timeOfDay.Minute(),
		timeOfDay.Second(),
		0,
		date.Location(),
	)
}
