package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
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

const version = "0.0.1"

const usage = `Maze: dynamic markdown static page server
version ` + version + `

USAGE
maze [option] [flags]

OPTIONS
serve			start the HTTP server
version			print application version

GLOBAL FLAGS
--help	-h		Print this usage message`

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
	if len(os.Args) == 1 {
		serve()
		return
	}
	command := os.Args[1]
	switch command {
	case "serve":
		serve()
	case "version":
		printVersion()
	default:
		flag.Usage()
	}
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "maze version: %s\n", version)
}

func serve() {
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
