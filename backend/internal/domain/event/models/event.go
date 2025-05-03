package models

import (
	"context"
	"fmt"
	"time"

	"github.com/hengadev/leviosa/internal/domain/stripe"

	"github.com/google/uuid"
)

type Event struct {
	ID                  string    `json:"id"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	City                string    `json:"city" encx:"encrypt"`
	CityEncrypted       []byte    `json:"-"`
	PostalCode          string    `json:"postal_code" encx:"encrypt"`
	PostalCodeEncrypted []byte    `json:"-"`
	Address1            string    `json:"address1" encx:"encrypt"`
	Address1Encrypted   []byte    `json:"-"`
	Address2            string    `json:"address2" encx:"encrypt"`
	Address2Encrypted   []byte    `json:"-"`
	PlaceCount          int       `json:"place_count"`
	FreePlace           int       `json:"free_place"`
	BeginAt             time.Time `json:"begin_at"`
	EndAt               time.Time `json:"end_at"`
	Products            []string  `json:"products"`
	Offers              []string  `json:"offers"`
	PriceID             string    `json:"-" encx:"encrypt"`
	PriceIDEncrypted    []byte    `json:"-"`
	Day                 int       `json:"day"`
	Month               int       `json:"month"`
	Year                int       `json:"year"`
}

func NewEvent(
	location string,
	placecount int,
	beginat time.Time,
	sessionduration time.Duration,
	day int,
	month int,
	year int,
) *Event {
	return &Event{
		ID:         uuid.NewString(),
		City:       location,
		BeginAt:    beginat,
		PlaceCount: placecount,
		FreePlace:  placecount,
		Day:        day,
		Month:      month,
		Year:       year,
	}
}

const description = "1 X Pass valuable for all the event."
const price = 1999

func (e *Event) getProductName() string {
	return fmt.Sprintf("Ticket pour l'evenement : %s", e.Title)
}

func (e *Event) GetPaymentInfo(ctx context.Context) *stripeService.PaymentProductInfo {
	return &stripeService.PaymentProductInfo{
		ID:          e.ID,
		Name:        e.getProductName(),
		Description: description,
		Price:       price,
	}
}

func (e Event) AssertComparable() {}

func (e Event) GetSQLColumnMapping() map[string]string {
	return map[string]string{
		"ID":               "id",
		"Title":            "encrypted_title",
		"Description":      "encrypted_description",
		"City":             "encrypted_city",
		"PostalCode":       "encrypted_postal_code",
		"Address1":         "encrypted_address1",
		"Address2":         "encrypted_address2",
		"PlaceCount":       "placecount",
		"FreePlace":        "freeplace",
		"EncryptedBeginAt": "encrypted_begin_at",
		"EncryptedEndAt":   "encrypted_end_at",
		"PriceID":          "encrypted_price_id",
		"Day":              "day",
		"Month":            "month",
		"Year":             "year",
	}
}

func (e Event) GetProhibitedFields() []string {
	return []string{"ID", "PriceID"}
}
