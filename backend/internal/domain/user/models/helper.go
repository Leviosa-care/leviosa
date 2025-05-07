package models

import (
	"fmt"
)

func convIntToMonth(month int) (string, error) {
	return convIntToStr(month, 13)
}

func convIntToDay(day, month int) (string, error) {
	var maxValue int
	if month == 2 {
		maxValue = 28
	} else if month%2 == 0 {
		maxValue = 30
	} else {
		maxValue = 31
	}
	return convIntToStr(day, maxValue)
}

func convIntToStr(value, maxValue int) (string, error) {
	if value < 0 {
		return "", fmt.Errorf("%d is negative value", value)
	} else if value < 10 {
		return fmt.Sprintf("0%d", value), nil
	} else if value < maxValue {
		return fmt.Sprintf("%d", value), nil
	} else {
		return "", fmt.Errorf("%d is too large value", value)
	}
}
