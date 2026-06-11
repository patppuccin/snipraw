package server

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patppuccin/snipraw/src/config"
	"github.com/patppuccin/snipraw/src/helpers"
)

//go:embed assets
var assets embed.FS

func router(appCtx *appCtx) (http.Handler, error) {
	r := chi.NewRouter()

	// middleware
	r.Use(loadServerContext(appCtx))
	r.Use(logRequest())
	r.Use(middleware.CleanPath)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.StripSlashes)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Language", "en")
			next.ServeHTTP(w, r)
		})
	})
	r.Use(middleware.Compress(5,
		"text/html",
		"text/css",
		"application/javascript",
		"application/json",
		"application/wasm",
		"application/xml",
		"text/plain",
		"text/javascript",
		"image/svg+xml",
	))

	// static assets
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		return nil, fmt.Errorf("failed to mount assets: %w", err)
	}
	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.FS(sub))))

	// routes
	r.Get("/", homePageHandler)
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/assets/icons/favicon.ico", http.StatusMovedPermanently)
	})
	r.Route("/project/{project}", func(r chi.Router) {
		r.Get("/", renderedProjectPageRedirectHandler)
		r.Get("/view", renderedProjectPageHandler)
		r.Get("/view/", renderedProjectPageHandler)
		r.Get("/view/*", renderedProjectPageHandler)
		r.Get("/blob/*", rawContentHandler)
		r.Get("/download", downloadContentHandler)
		r.Get("/download/", downloadContentHandler)
		r.Get("/download/*", downloadContentHandler)
	})

	// Error handlers
	r.NotFound(notFoundHandler)
	r.MethodNotAllowed(methodNotAllowedHandler)

	return r, nil
}

func Run(ctx context.Context, runtime *Runtime, cfg *config.Config) error {

	logger := helpers.InitLogger(runtime.LogLevel)

	appCtx := &appCtx{
		dir:    runtime.Dir,
		config: cfg,
		logger: &logger,
	}

	router, err := router(appCtx)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", runtime.Host, runtime.Port),
		Handler: router,
	}

	errCh := make(chan error, 1)

	go func() {
		logger.Info().
			Str("host", runtime.Host).
			Int("port", runtime.Port).
			Str("dir", runtime.Dir).
			Msg("starting snipraw")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info().Msg("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
