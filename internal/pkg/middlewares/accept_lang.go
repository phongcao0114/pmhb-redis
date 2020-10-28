package middlewares

import (
	"context"
	"net/http"
)

const (
	acceptLanguageContextKey = "accept_language"
	acceptLanguageEN         = "en"
)

// NewAcceptLanguageContext return new context with language value
func NewAcceptLanguageContext(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, acceptLanguageContextKey, lang)
}

// AcceptLanguageFromContext return language value stored in context
func AcceptLanguageFromContext(ctx context.Context) string {
	if lang := ctx.Value(acceptLanguageContextKey); lang != nil {
		return lang.(string)
	}
	return acceptLanguageEN
}

// AcceptLanguage is a middleware to inject Accept-Language into request context
func AcceptLanguage(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get("Accept-Language")
		if lang == "" {
			lang = acceptLanguageEN
		}
		r = r.WithContext(NewAcceptLanguageContext(r.Context(), lang))
		next.ServeHTTP(w, r)
	})
}
