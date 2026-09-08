package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mbaklor/maze/paths"
	"github.com/mbaklor/maze/routes/pages"
)

type TemplateInfo struct {
	Path  string
	Links []paths.Link
}

type WebApp struct {
	logger *slog.Logger
	server *http.Server
}

const serveUsage = `Maze: dynamic markdown static page server
version ` + version + `

USAGE
maze serve [flags]

FLAGS
-c	--config	Path to server config file (default: "config.yml")

GLOBAL FLAGS
-h	--help		Print this usage message`

func serve(args []string) {
	set := flag.NewFlagSet("serve", flag.ExitOnError)
	set.Usage = UsageFunc(serveUsage)
	var configPath string
	set.StringVar(&configPath, "config", "config.yml", "path to server config file")
	set.StringVar(&configPath, "c", "config.yml", "path to server config file")

	set.Parse(args)
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	if err := runServer(logger); err != nil {
		logger.Error("App can't run!", slog.String("error", err.Error()))
	}
}

func runServer(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	r := chi.NewRouter()
	server := &http.Server{
		Handler: r,
		Addr:    ":9753",
	}
	w := WebApp{logger, server}
	r.Use(w.LogRequests)

	r.Route("/static", w.StaticRouter)

	r.Route("/", w.RootRouter)
	w.logger.Info("Started serving", slog.String("address", w.server.Addr))
	var shutdownErr error
	go func() {
		<-ctx.Done()
		w.logger.Info("server shutdown requested")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		shutdownErr = w.server.Shutdown(ctx)
	}()
	err := w.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return errors.Join(err, shutdownErr)
	}
	return nil
}

func (wa *WebApp) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wa.logger.Info("got request from client", slog.String("path", r.URL.Path), slog.String("client_addr", r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}

func (wa *WebApp) StaticRouter(r chi.Router) {
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(paths.FrontendPath("static"))))

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		wa.logger.Info("in the static handler", "path", r.URL.Path)
		fs.ServeHTTP(w, r)
	})
}

func (wa *WebApp) RootRouter(r chi.Router) {
	r.Get("/*", pages.Handler(wa.logger))
}
