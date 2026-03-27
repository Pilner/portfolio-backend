package util

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/middleware"
)

type LoggerHandler struct {
	slog.Handler
}

func (h LoggerHandler) Handle(ctx context.Context, record slog.Record) error {
	// get requestId from middleware and insert it as an attribute, if available
	if reqIdVal := ctx.Value(middleware.RequestIDKey); reqIdVal != nil {
		if reqId, ok := reqIdVal.(string); ok && reqId != "" {
			record.AddAttrs(slog.String("requestId", reqId))
		}
	}

	return h.Handler.Handle(ctx, record)
}

func (h LoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return LoggerHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h LoggerHandler) WithGroup(name string) slog.Handler {
	return LoggerHandler{Handler: h.Handler.WithGroup(name)}
}
