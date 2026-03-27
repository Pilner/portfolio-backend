package middlewares

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
)

func SetRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Header.Get(middleware.RequestIDHeader)
		if reqId == "" {
			reqId = uuid.NewString()
		}
		ctx := context.WithValue(r.Context(), middleware.RequestIDKey, reqId)
		w.Header().Set(middleware.RequestIDHeader, reqId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}