package stringutil

import (
	"errors"
	"time"
)

func ToDate(s string) *time.Time {
	d, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &d
}

func ToDateValidate(s string) (*time.Time, error) {
	if s == "" {
		return nil, errors.New("request_datetime is missing")
	}
	d, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, errors.New("request_datetime is wrong format")
	}
	return &d, nil
}
