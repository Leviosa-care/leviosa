package eventRepository_test

import (
	"github.com/hengadev/leviosa/internal/domain/event/models"
	"github.com/hengadev/leviosa/internal/repository/postgres/event"
)

var (
	beginAt, _  = eventRepository.ExportedParseBeginAt("08:00:00", 12, 7, 1998)
	beginAt1, _ = eventRepository.ExportedParseBeginAt("08:00:00", 13, 7, 1998)
	beginAt2, _ = eventRepository.ExportedParseBeginAt("08:00:00", 14, 7, 1998)

	baseEvent = &models.Event{
		ID:         "ea1d74e2-1612-47ec-aee9-c6a46b65640f",
		PlaceCount: 16,
		FreePlace:  14,
		BeginAt:    beginAt,
		Day:        12,
		Month:      7,
		Year:       1998,
	}

	baseEvent1 = &models.Event{
		ID:         "b16a6f38-d2fb-428c-b97c-929b1010b951",
		PlaceCount: 23,
		FreePlace:  19,
		BeginAt:    beginAt1,
		Day:        13,
		Month:      7,
		Year:       1998,
	}

	baseEvent2 = &models.Event{
		ID:         "9a676c5d-c9ec-4266-a426-24e5d4983caf",
		PlaceCount: 22,
		FreePlace:  0,
		BeginAt:    beginAt2,
		Day:        14,
		Month:      7,
		Year:       1998,
	}

	baseEventWithPriceID = &models.Event{
		ID:         "ea1d74e2-1612-47ec-aee9-c6a46b65640f",
		PlaceCount: 16,
		FreePlace:  14,
		BeginAt:    beginAt,
		PriceID:    "4fe0vuw3ef0223",
		Day:        12,
		Month:      7,
		Year:       1998,
	}
)
