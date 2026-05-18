package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func loadServerContext(appCtx *appCtx) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), appCtxKey, appCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func logRequest() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			appCtx, ok := r.Context().Value(appCtxKey).(*appCtx)
			if !ok || appCtx == nil || appCtx.logger == nil {
				fmt.Println("Unable to fetch the logger from the server context")
				return
			}

			status := ww.Status()
			duration := time.Since(start)

			event := appCtx.logger.Info()
			switch {
			case status >= 500:
				event = appCtx.logger.Error()
			case status >= 400:
				event = appCtx.logger.Warn()
			case status >= 300:
				event = appCtx.logger.Info()
			case status >= 200:
				event = appCtx.logger.Debug()
			}

			event.
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", status).
				Dur("duration", duration).
				Msg("HTTP request")
		})
	}
}
