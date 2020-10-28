package models

import (
	"context"
	"time"
)

type (
	RequestInfo struct {
		Header KbankRequestHeader `json:"kbank_header"`
		Body   interface{}        `json:"request_body"`
	}
	KbankRequestHeader struct {
		RequestAppID string `json:"request_app_id" kvalid:"request_app_id"`
		RequestDate  string `json:"request_datetime" kvalid:"request_datetime"`
		RequestUID   string `json:"request_id" kvalid:"request_id"`
	}
	KbankResponseHeader struct {
		ResponseAppID string                   `json:"response_app_id"`
		ResponseUID   string                   `json:"response_id"`
		ResponseDate  time.Time                `json:"response_datetime"`
		StatusCode    string                   `json:"status_code"`
		Errors        ResponseErrorKbankHeader `json:"error"`
	}
	ResponseErrorKbankHeader struct {
		ErrorCode string `json:"error_code"`
		ErrorDesc string `json:"error_desc"`
	}
)

func NewRequestHeader(ctx context.Context) KbankRequestHeader {
	requestID, _ := ctx.Value("request_id").(string)
	appID, _ := ctx.Value("request_app_id").(string)
	now := time.Now()
	return KbankRequestHeader{
		RequestAppID: appID,
		RequestDate:  now.Format(time.RFC3339),
		RequestUID:   requestID,
	}
}
