package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type WebApp struct {
	logger *slog.Logger
	server *http.Server
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	if err := run(logger); err != nil {
		logger.Error("App can't run!", slog.String("error", err.Error()))
	}
}

func run(logger *slog.Logger) error {
	r := chi.NewRouter()
	server := &http.Server{
		Handler: r,
		Addr:    ":9753",
	}
	w := WebApp{logger, server}
	r.Use(w.LogRequests)

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir("templates/static")))

	r.HandleFunc("/static/*", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("in the static handler")
		fs.ServeHTTP(w, r)
	})

	r.Route("/", w.RootRouter)
	w.logger.Info("Started serving", slog.String("address", w.server.Addr))
	err := w.server.ListenAndServe()
	return err
}

func (wa *WebApp) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wa.logger.Info("got request from client", slog.String("path", r.URL.Path), slog.String("client_addr", r.RemoteAddr))
		next.ServeHTTP(w, r)
	})
}

func (wa *WebApp) StaticHandler() {}

func (wa *WebApp) RootRouter(r chi.Router) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			wa.logger.Error("error serving page", slog.String("error", err.Error()))
			http.Error(w, "can't open page: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})
}
