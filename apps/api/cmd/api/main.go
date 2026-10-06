package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Giandri/Gidocs/apps/api/internal/config"
	apihandler "github.com/Giandri/Gidocs/apps/api/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(cfg.RequestTimeout))
	r.Use(corsMiddleware(cfg.CORSAllowedOrigins))

	r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/api/merge", apihandler.MergeHandler)
	r.Post("/api/compress", apihandler.CompressHandler)
	r.Post("/api/rotate", apihandler.RotateHandler)
	r.Post("/api/split", apihandler.SplitHandler)
	r.Post("/api/remove-pages", apihandler.RemovePagesHandler)
	r.Post("/api/reorder", apihandler.ReorderHandler)
	r.Post("/api/protect", apihandler.ProtectHandler)
	r.Post("/api/unlock", apihandler.UnlockHandler)
	r.Post("/api/watermark", apihandler.WatermarkHandler)
	r.Post("/api/images-to-pdf", apihandler.ImagesToPDFHandler)

	r.Post("/api/jobs", apihandler.EnqueueJob)
	r.Get("/api/jobs/{id}", apihandler.JobStatus)
	r.Get("/api/jobs/{id}/download", apihandler.JobDownload)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("server berjalan", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func corsMiddleware(allowedOrigins string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if allowedOrigins == "" || origin == "" || matchesOrigin(origin, allowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func matchesOrigin(origin, allowed string) bool {
	for _, o := range splitCSV(allowed) {
		if o == origin || o == "*" {
			return true
		}
	}
	return false
}

func splitCSV(s string) []string {
	var result []string
	for _, part := range splitCSVBytes([]byte(s)) {
		result = append(result, string(part))
	}
	return result
}

func splitCSVBytes(b []byte) [][]byte {
	var parts [][]byte
	var current []byte
	for _, c := range b {
		if c == ',' {
			parts = append(parts, trimSpace(current))
			current = nil
		} else {
			current = append(current, c)
		}
	}
	if len(current) > 0 {
		parts = append(parts, trimSpace(current))
	}
	return parts
}

func trimSpace(b []byte) []byte {
	start := 0
	for start < len(b) && (b[start] == ' ' || b[start] == '\t') {
		start++
	}
	end := len(b)
	for end > start && (b[end-1] == ' ' || b[end-1] == '\t') {
		end--
	}
	return b[start:end]
}