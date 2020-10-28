package response

import (
	"context"
	"encoding/json"
	"net/http"
	"pmhb-redis/internal/app/models"
	"pmhb-redis/internal/app/utils"
	"pmhb-redis/internal/kerrors"
	"pmhb-redis/internal/pkg/klog"
	"time"
)

const (
	StatusSuccess = "00"
	StatusError   = "01"
)

var logger = klog.WithPrefix("response")

// SuccessResponseFormat struct contains success response data
type SuccessResponseFormat struct {
	KbankResponseHeader models.KbankResponseHeader `json:"kbank_header"`
	ServiceResponseBody *interface{}               `json:"body"`
}

// FailureResponseFormat struct contains failed response data
type FailureResponseFormat struct {
	KbankResponseHeader models.KbankResponseHeader `json:"kbank_header"`
}

// WriteJSON writes JSON data into responseWriter
func WriteJSON(w http.ResponseWriter) func(resp interface{}, httpStatusCode int) {
	return func(resp interface{}, httpStatusCode int) {
		w.Header().Set("Content-type", "application/json")
		json.NewEncoder(w).Encode(resp)
		w.WriteHeader(httpStatusCode)
	}
}

// HandleError function returns failure response
func HandleError(r *http.Request, headerRequest models.KbankRequestHeader, err error) (FailureResponseFormat, int) {
	ctx := r.Context()
	logger := klog.WithPrefix("response")
	var res FailureResponseFormat

	switch v := err.(type) {
	case kerrors.KError:
		logger.WithFields(v.Extract()).KError(ctx, v.LogMessage)

		errCode, errMsg := v.Code.String(), v.Message.String()
		res = FailureResponseFormat{
			models.KbankResponseHeader{
				ResponseAppID: utils.ResponseAppID,
				ResponseDate:  time.Now(),
				StatusCode:    kerrors.StatusErrorFailed.String(),
				Errors: models.ResponseErrorKbankHeader{
					ErrorCode: errCode,
					ErrorDesc: errMsg,
				},
			},
		}
	default:
		logger.KError(ctx, v.Error())
	}
	c := context.WithValue(ctx, "request_status", kerrors.StatusErrorFailed.String())
	temp := r.WithContext(c)
	*r = *temp

	return res, http.StatusOK
}

// HandleSuccess function responses success response format
func HandleSuccess(r *http.Request, headerRequest models.KbankRequestHeader, data interface{}) (SuccessResponseFormat, int) {
	rs := SuccessResponseFormat{
		KbankResponseHeader: models.KbankResponseHeader{
			ResponseAppID: utils.ResponseAppID,
			ResponseDate:  time.Now(),
			StatusCode:    kerrors.StatusNoneError.String(),
			ResponseUID:   headerRequest.RequestUID,
			Errors: models.ResponseErrorKbankHeader{
				ErrorCode: kerrors.NoneError.String(),
				ErrorDesc: kerrors.ErrDictionary.Get(kerrors.NoneError).String(),
			},
		},
	}
	if data != nil {
		rs.ServiceResponseBody = &data
	}
	c := context.WithValue(r.Context(), "request_status", kerrors.StatusNoneError.String())
	r2 := r.WithContext(c)
	*r = *r2

	return rs, http.StatusOK
}
