package vote_test

import (
	"context"
	"testing"
	"time"

	"github.com/hengadev/leviosa/internal/domain/vote"

	"github.com/hengadev/test-assert"
)

func TestValidVote(t *testing.T) {
	now := time.Now().UTC()
	day := now.Day()
	month := int(now.Month())
	year := now.Year()
	tests := []struct {
		day     int
		month   int
		year    int
		wantErr bool
		name    string
	}{
		{day: -2, month: month, year: year, wantErr: true, name: "day too small"},
		{day: 54, month: month, year: year, wantErr: true, name: "day too large"},
		{day: day, month: -4, year: year, wantErr: true, name: "month too small"},
		{day: day, month: 45, year: year, wantErr: true, name: "month too large"},
		{day: day, month: month, year: year - 3, wantErr: true, name: "year too small"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			vote := &vote.Vote{Day: tt.day, Month: tt.month, Year: tt.year}
			err := vote.Valid(ctx)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}
