package eventRepository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/hengadev/leviosa/internal/repository/postgres/event"

	"github.com/hengadev/test-assert"
)

func TestConvIntToStr(t *testing.T) {
	tests := []struct {
		value          int
		expectedString string
		wantErr        bool
		name           string
	}{
		{value: -8, expectedString: "", wantErr: true, name: "Invalid case with negative value"},
		{value: 3234, expectedString: "", wantErr: true, name: "Invalid case with value > 99 value"},
		{value: 5, expectedString: "05", wantErr: false, name: "Value inferior than 10"},
		{value: 12, expectedString: "12", wantErr: false, name: "Value superior than 10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			str, err := eventRepository.ExportedConvIntToStr(tt.value)
			assert.Equal(t, str, tt.expectedString)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		value          string
		expectedString string
		wantErr        bool
		name           string
	}{
		{value: "wrong format", expectedString: "", wantErr: true, name: "invalid time, not even letters"},
		{value: "08:12:57", expectedString: "08:12:57AM", wantErr: false, name: "valid morning time"},
		{value: "16:34:23", expectedString: "04:34:23PM", wantErr: false, name: "valid afternoon time"},
		{value: "33:34:23", expectedString: "", wantErr: true, name: "invalid time, hour > 24"},
		{value: "-23:34:23", expectedString: "", wantErr: true, name: "invalid time, hour < 0"},
		{value: "09:-23:23", expectedString: "", wantErr: true, name: "invalid time, minute < 0"},
		{value: "11:68:23", expectedString: "", wantErr: true, name: "invalid time, minute > 60"},
		{value: "18:15:-43", expectedString: "", wantErr: true, name: "invalid time, second < 0"},
		{value: "19:32:72", expectedString: "", wantErr: true, name: "invalid time, second > 60"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := eventRepository.ExportedFormatTime(tt.value)
			assert.Equal(t, res, tt.expectedString)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestParseBeginAt(t *testing.T) {
	var zeroValue time.Time

	v := "07/12 07:25:17AM '98 -0700"
	vParsed, _ := time.Parse(time.Layout, v)
	expectedTime := vParsed.String()

	tests := []struct {
		hour         string
		day          int
		month        int
		year         int
		expectedTime string
		wantErr      bool
		name         string
	}{
		{hour: "07:25:17", day: 45, month: 7, year: 1998, expectedTime: zeroValue.String(), wantErr: true, name: "day > 30"},
		{hour: "07:25:17", day: -32, month: 7, year: 1998, expectedTime: zeroValue.String(), wantErr: true, name: "day < 0"},
		{hour: "07:25:17", day: 12, month: 32, year: 1998, expectedTime: zeroValue.String(), wantErr: true, name: "month > 12"},
		{hour: "07:25:17", day: 12, month: 0, year: 1998, expectedTime: zeroValue.String(), wantErr: true, name: "month < 1"},
		{hour: "07:25:17", day: 12, month: 7, year: 1998, expectedTime: expectedTime, wantErr: false, name: "nominal case"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res, err := eventRepository.ExportedParseBeginAt(tt.hour, tt.day, tt.month, tt.year)
			assert.Equal(t, res.String(), tt.expectedTime)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestFormatBeginAt(t *testing.T) {
	now := time.Now().UTC()
	// nowStr := now.Format(time.Layout)
	// beginat, err := time.Parse(time.Layout, nowStr)
	// if err != nil {
	// 	t.Error("failed intial time for test")
	// }

	hour, _ := eventRepository.ExportedConvIntToStr(now.Hour())
	minute, _ := eventRepository.ExportedConvIntToStr(now.Minute())
	second, _ := eventRepository.ExportedConvIntToStr(now.Second())

	expectedbeginat := fmt.Sprintf("%s:%s:%s", hour, minute, second)

	tests := []struct {
		beginAt         time.Time
		expectedBeginat string
		wantErr         bool
		name            string
	}{
		// {beginAt: beginat, expectedBeginat: expectedbeginat, wantErr: false, name: "first test, ba mwen tan"},
		{beginAt: now, expectedBeginat: expectedbeginat, wantErr: false, name: "first test, ba mwen tan"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			formatbeginat, err := eventRepository.ExportedFormatBeginAt(tt.beginAt)
			assert.Equal(t, err != nil, tt.wantErr)
			assert.Equal(t, formatbeginat, tt.expectedBeginat)
		})
	}
}
