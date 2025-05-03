package vote

import (
	"context"
	"reflect"
	"time"

	"github.com/hengadev/errsx"
)

const VoteSeparator = "|"

type Vote struct {
	UserID string `json:"userid,omitempty"`
	Day    int    `json:"day,omitempty"`
	Month  int    `json:"month"`
	Year   int    `json:"year"`
}

type AvailableVote struct {
	Day   int `json:"day,omitempty"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

func NewVote(day, month, year int) *Vote {
	return &Vote{
		Day:   day,
		Month: month,
		Year:  year,
	}
}

// NOTE: the way I am going to make the tables
// vote : month year days
// vote_august_2024 : userID days

func (v Vote) Valid(ctx context.Context) error {
	var errs = make(errsx.Map)
	vf := reflect.VisibleFields(reflect.TypeOf(v))
	now := time.Now().UTC()
	for _, f := range vf {
		switch f.Name {
		case "Day":
			if int(now.Month())%2 == 0 && int(now.Month()) != 8 && v.Day > 30 {
				errs.Set("not_enough_day", "this month has 30 days")
			}
			if v.Day < 1 {
				errs.Set("day_too_small", "should be > 1")
			}
			if v.Day > 31 {
				errs.Set("day_too_large", "should be < 32")
			}
		case "Month":
			if v.Month <= int(now.Month()) && v.Year == now.Year() {
				errs.Set("past_month", "date should not be in the past")
			}
			if v.Month < 1 {
				errs.Set("month_too_small", "should be > 1")
			}
			if v.Month > 12 {
				errs.Set("month_too_large", "should be < 13")
			}
		case "Year":
			if v.Year < now.Year() {
				errs.Set("year", "should be > than current year")
			}
		default:
			continue
		}
	}
	return errs.AsError()
}
