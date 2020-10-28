package middlewares

import (
	"fmt"
	"net/http"
	"pmhb-redis/internal/pkg/klog"
	"runtime"
)

// Recover is a middleware that recover a handler from a panic
func Recover(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err, ok := rec.(error)
				if !ok {
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4<<10) // 4KB
				length := runtime.Stack(stack, false)

				klog.WithPrefix("middleware_recover").KErrorf(r.Context(), "panic recover, err: %v, stack: %s", err, stack[:length])
				http.Error(w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
