package header

import (
	"time"
)

// Header interface groups all functions together
type Header interface {
	GetAppID() string
	GetRequestID() string
	GetRequestDate() (*time.Time, error)
}

// ValidateHeader validates header function
func ValidateHeader(header Header) (err error) {
	return nil
}
