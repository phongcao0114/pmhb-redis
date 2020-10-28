package middlewares

import (
	"net/http"
	"pmhb-redis/internal/pkg/klog"
)

const (
	// ServiceName declares PH Srv
	ServiceName = "PH-Service"
)

// LoggerWithRequestMeta is a middleware that inject request information into a logger
func LoggerWithRequestMeta(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := klog.WithFields(map[string]interface{}{
			"path":         r.URL.Path,
			"method":       r.Method,
			"service_name": ServiceName,
		})
		c := klog.NewContext(r.Context(), logger)
		r = r.WithContext(c)
		next.ServeHTTP(w, r)
	})
}
