package utils

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"pmhb-redis/internal/app/models"
	"pmhb-redis/internal/kerrors"
	"pmhb-redis/internal/pkg/mapper"
	"time"
)

var (
	// ResponseAppID returns appID from system
	ResponseAppID string

	// BKKLocation contains time location
	BKKLocation *time.Location

	appIDKey  = "request_app_id"
	reqIDKey  = "request_id"
	dateIDKey = "request_datetime"

	// LogKey contains log timing details
	LogKey = "log_request"
)

// GetRequestID function returns request ID
func GetRequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(reqIDKey).(string)
	return requestID
}

// GetAppID function returns app ID
func GetAppID(ctx context.Context) string {
	appID, _ := ctx.Value(appIDKey).(string)
	return appID
}

// GetRequestDate function returns request date
func GetRequestDate(ctx context.Context) string {
	date, _ := ctx.Value(dateIDKey).(string)
	return date
}

// ParsePagination return current page and total page
func ParsePagination(limit, total, offset int) (currentPage int, totalPage int) {
	if limit == 0 {
		limit = 20
	}

	currentPage = offset/limit + 1
	return currentPage, total
}

// DecodeToBody function is decoding all general requests from all API
func DecodeToBody(kerror *kerrors.KError, requestInfo *models.RequestInfo, body interface{}, r *http.Request) (err error) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	if err = decoder.Decode(&requestInfo); err != nil {
		err = kerror.Wrap(err, kerrors.CannotDecodeInputRequest, nil)
		return
	}

	if err = mapper.ConvertMapToModel(body, requestInfo.Body); err != nil {
		err = kerror.Wrap(err, kerrors.MarshalFail, nil)
		return
	}
	requestInfo.Body = body
	return
}

// ConvertToModel is converting byte data to model
func ConvertToModel(data interface{}, model interface{}) error {
	byte, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(byte, model); err != nil {
		return err
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func GenerateTransID() string {
	b := make([]byte, 10)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", b)
}

func GetPtrIntData(input *int) int {
	if input == nil {
		return 0
	}
	return *input
}

func GetPtrStrData(input *string) string {
	if input == nil {
		return ""
	}
	return *input
}

func GetPtrFloat64Data(input *float64) float64 {
	if input == nil {
		return 0.0
	}
	return *input
}

func SplitCardDate(month, year string) string {
	if len(month) == 0 || len(year) == 0 {
		return ""
	}
	arrStr := []rune(year)
	return string(arrStr[len(arrStr)-2:]) + month
}
