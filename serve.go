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
	"github.com/mbaklor/maze/settings"
)

type TemplateInfo struct {
	Path  string
	Links []paths.Link
}

type WebApp struct {
	logger *slog.Logger
	server *http.Server
	config settings.Settings
}

const serveUsage = `Maze: dynamic markdown static page server
version ` + version + `

USAGE
maze serve [flags]

FLAGS
-c	--config	Path to server config file							(default: "config.yml")
-f	--files		Path to folder where markdown files are stored		(default: "frontend/pages")
-p	--port		Port for server to listen on						(default "9753")
-t	--title		Main title for the website, to show in <title> tag	(default: "Maze Site")

GLOBAL FLAGS
-h	--help		Print this usage message`

func parseServeFlags(args []string) settings.Settings {
	var s settings.Settings
	set := flag.NewFlagSet("serve", flag.ExitOnError)
	set.Usage = UsageFunc(serveUsage)

	set.StringVar(&s.ConfigFile, "config", "", "path to server config file")
	set.StringVar(&s.ConfigFile, "c", "", "path to server config file")

	set.StringVar(&s.FrontendPath, "files", "", "directory where markdown files are stored")
	set.StringVar(&s.FrontendPath, "f", "", "directory where markdown files are stored")

	set.IntVar(&s.ServerPort, "port", 0, "port for server")
	set.IntVar(&s.ServerPort, "p", 0, "port for server")

	set.StringVar(&s.SiteTitle, "title", "", "title for the website")
	set.StringVar(&s.SiteTitle, "t", "", "title for the website")

	set.Parse(args)
	return s
}

func mergeServeSettings(args []string, logger *slog.Logger) settings.Settings {
	var s settings.Settings
	c := parseServeFlags(args)
	e, err := settings.ReadEnvVars()
	if err != nil {
		logger.Warn("got error while reading env vars", "error message", err)
	}
	merge("config.yml", &s.ConfigFile, c.ConfigFile, e.ConfigFile)

	f, err := settings.ReadSettingsFile(s.ConfigFile)
	if err != nil {
		logger.Warn("got error while reading config file", "error message", err, "config file", s.ConfigFile)
	}

	merge("frontend/pages", &s.FrontendPath, c.FrontendPath, e.FrontendPath, f.FrontendPath)
	merge(9753, &s.ServerPort, c.ServerPort, e.ServerPort, f.ServerPort)
	merge("Maze Site", &s.SiteTitle, c.SiteTitle, e.SiteTitle, f.SiteTitle)
	return s
}

func serve(args []string) {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	s := mergeServeSettings(args, logger)

	if err := runServer(logger, s); err != nil {
		logger.Error("App can't run!", slog.String("error", err.Error()))
	}
}

func runServer(logger *slog.Logger, config settings.Settings) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	r := chi.NewRouter()
	addr := fmt.Sprintf(":%d", config.ServerPort)
	server := &http.Server{
		Handler: r,
		Addr:    addr,
	}

	w := WebApp{logger, server, config}
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
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static")))

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		wa.logger.Info("in the static handler", "path", r.URL.Path)
		fs.ServeHTTP(w, r)
	})
}

func (wa *WebApp) RootRouter(r chi.Router) {
	rh := pages.NewRouteHandler(wa.config.SiteTitle, wa.config.FrontendPath, wa.logger)
	r.Get("/*", rh.Handle)
}
