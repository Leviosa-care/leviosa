package eventRepository

import (
	"fmt"
	"time"
)

func convIntToStr(value int) (string, error) {
	if value < 0 {
		return "", fmt.Errorf("%d is invalid value", value)
	} else if value < 10 {
		return fmt.Sprintf("0%d", value), nil
	} else if value < 100 {
		return fmt.Sprintf("%d", value), nil
	} else {
		return "", fmt.Errorf("%d is invalid value", value)
	}
}

// TODO: change the API so that I can have better function like this one to restrict the API for better use.
// func convIntToStr(value, maxValue int) (string, error) {
// 	if value < 0 {
// 		return "", fmt.Errorf("%d is invalid value", value)
// 	} else if value < 10 {
// 		return fmt.Sprintf("0%d", value), nil
// 	} else if value < maxValue {
// 		return fmt.Sprintf("%d", value), nil
// 	} else {
// 		return "", fmt.Errorf("%d is invalid value", value)
// 	}
// }
// func convIntToHour(value int) (string, error) {
// 	return convIntToStr(value, 24)
// }
//
// func convIntToMinuteOrSecond(value int) (string, error) {
// 	return convIntToStr(value, 60)
// }

// helper function for the GetEventForUser function
func formatTime(hour string) (string, error) {
	res := hour
	suffix := "AM"
	timeHour, err := time.Parse(time.TimeOnly, hour)
	if err != nil {
		return "", fmt.Errorf("error parsing time: %w", err)
	}
	if timeHour.Hour() > 12 {
		suffix = "PM"
		hour, err := convIntToStr(timeHour.Hour() - 12)
		if err != nil {
			return "", fmt.Errorf("convert string to int: %w", err)
		}
		minute, err := convIntToStr(timeHour.Minute())
		if err != nil {
			return "", fmt.Errorf("convert string to int: %w", err)
		}
		second, err := convIntToStr(timeHour.Second())
		if err != nil {
			return "", fmt.Errorf("convert string to int: %w", err)
		}

		res = fmt.Sprintf("%s:%s:%s", hour, minute, second)
	}
	return res + suffix, nil
}

// helper function for the GetEventForUser function
func parseBeginAt(hour string, day, month, year int) (time.Time, error) {
	var res time.Time
	hourFormatted, err := formatTime(hour)
	if err != nil {
		return res, err
	}
	parsedDay, err := convIntToStr(day)
	if err != nil {
		return res, fmt.Errorf("convert string to int: %w", err)
	}
	parsedMonth, err := convIntToStr(month)
	if err != nil {
		return res, fmt.Errorf("convert string to int: %w", err)
	}
	parsedYear, err := convIntToStr(year % 100)
	if err != nil {
		return res, fmt.Errorf("convert string to int: %w", err)
	}
	dateFormatted := fmt.Sprintf("%s/%s %s '%s -0700", parsedMonth, parsedDay, hourFormatted, parsedYear)
	res, err = time.Parse(time.Layout, dateFormatted)
	if err != nil {
		return res, fmt.Errorf("parsing formatted date: %w", err)
	}
	return res, nil
}

// Given a time.Time return hour, minute, second in this format "08:00:00"
func formatBeginAt(beginAt time.Time) (string, error) {
	hour, _ := convIntToStr(beginAt.Hour())
	minute, _ := convIntToStr(beginAt.Minute())
	second, _ := convIntToStr(beginAt.Second())
	return fmt.Sprintf("%s:%s:%s", hour, minute, second), nil
}
