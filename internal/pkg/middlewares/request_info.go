package middlewares

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httputil"
	"pmhb-redis/internal/app/utils"
	"pmhb-redis/internal/pkg/klog"
	"strings"
	"time"
)

// RequestInfo is a middleware that will print request info
func RequestInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			next.ServeHTTP(w, r)
			return
		}

		logger := klog.WithPrefix("middleware_request_info")
		start := time.Now()
		logger.WithFields(map[string]interface{}{
			"content_type": r.Header.Get("Content-Type"),
		}).KInfof(r.Context(), "started request")
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			logger.KErrorln(r.Context(), "unable to dump request:", err)
		}
		ct := r.Header.Get("Content-Type")
		if strings.HasPrefix(ct, "application/json") || strings.HasPrefix(ct, "application/xml") {
			logger.KDebugf(r.Context(), "request: %v", string(dump))
		}

		buff := bytes.Buffer{}
		sw := statusWriter{
			ResponseWriter: w,
			body:           &buff,
		}

		next.ServeHTTP(&sw, r)

		// skip write log if response is image or terms
		if strings.HasPrefix(ct, "application/json") || strings.HasPrefix(ct, "application/xml") {
			logger.KDebugf(r.Context(), "response: %s", sw.body.String())
		}

		statusCode := sw.status
		if statusCode == 0 {
			statusCode = http.StatusOK
		}

		requestStatus := r.Context().Value("request_status")
		responseTime := time.Since(start)

		LogTimingServices(r.Context())
		logger.WithFields(map[string]interface{}{
			"status_code":    statusCode,
			"status":         http.StatusText(statusCode),
			"response_time":  responseTime.Seconds(),
			"request_status": requestStatus,
			"request_id":     utils.GetRequestID(r.Context()),
			"url_api_path":   r.URL.Path,
		}).KInfof(r.Context(), "Completed %s %s %v %s in %v",
			r.Method,
			r.URL.Path,
			statusCode,
			http.StatusText(statusCode),
			responseTime,
		)
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	w.body.Write(b)

	return w.ResponseWriter.Write(b)
}

func LogTimingServices(ctx context.Context) {

}
