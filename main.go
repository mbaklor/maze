package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/mbaklor/website/md"
)

type TemplateInfo struct {
	Path string
}

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

	r.Route("/static", w.StaticRouter)

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

func (wa *WebApp) StaticRouter(r chi.Router) {
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(GetStaticDir())))

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		wa.logger.Info("in the static handler", "path", r.URL.Path)
		fs.ServeHTTP(w, r)
	})
}

func (wa *WebApp) RootRouter(r chi.Router) {
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path, err := ParseFileFromUrl(r.URL.Path)
		if err != nil {
			wa.logger.Error("error serving page, can't parse path from url", slog.String("error", err.Error()))
			http.Error(w, "can't get path from url: "+err.Error(), http.StatusInternalServerError)
			return
		}

		wa.logger.Debug("path information", "page path", r.URL.Path, "file", path)
		tmpl, err := template.ParseFiles("frontend/templates/base.html")
		if err != nil {
			wa.logger.Error("error serving page, can't open base path", slog.String("error", err.Error()))
			http.Error(w, "can't open base path: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = wa.parseMD(tmpl, path)
		if err != nil {
			wa.logger.Error("error serving page", slog.String("error", err.Error()))
			http.Error(w, "can't open page: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, TemplateInfo{Path: r.URL.Path})
	})
}

func (wa WebApp) parseMD(tmpl *template.Template, filename string) error {
	wa.logger.Info("parsing markdown file", "filename", filename)
	m, err := md.ParseMarkdownFile(filename)
	s := fmt.Sprintf("{{define \"title\"}} Mic's Web - %s {{end}} {{define \"body\"}}\n%s\n{{end}}", m.Info.Title, m.Content)
	_, err = tmpl.Parse(string(s))
	if err != nil {
		return err
	}
	return nil

}
