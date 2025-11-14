package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
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
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.PathValue("*"), "/")
		basePath := path[0]
		pagePath := path[len(path)-1]
		file := pagePath + ".md"
		if path[0] == "" {
			file = "index.html"
		}
		wa.logger.Debug("path information", "page path", pagePath, "file", file)
		tmpl, err := template.ParseFiles("templates/base.html")
		if err != nil {
			wa.logger.Error("error serving page, can't open base path", slog.String("error", err.Error()))
			http.Error(w, "can't open base path: "+err.Error(), http.StatusInternalServerError)
			return
		}

		if filepath.Ext(file) == ".md" {
			err = wa.parseMD(tmpl, file)
		} else {
			_, err = tmpl.ParseFiles(filepath.Join("templates", file))
		}
		if err != nil {
			wa.logger.Error("error serving page", slog.String("error", err.Error()))
			http.Error(w, "can't open page: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, TemplateInfo{Path: "/" + basePath})
	})
}

func findHeading(node ast.Node) ast.Node {
	if heading, ok := node.(*ast.Heading); ok {
		h := heading.GetChildren()[0].AsLeaf()
		return h
	}
	children := node.GetChildren()
	if len(children) == 0 {
		return nil
	}
	return findHeading(children[0])
}

func (wa WebApp) parseMD(tmpl *template.Template, file string) error {
	wa.logger.Info("parsing markdown file")
	var f []byte
	f, err := os.ReadFile(filepath.Join("templates", file))
	if err != nil {
		return err
	}
	md := markdown.Parse(f, nil)
	heading := findHeading(md)
	h1 := heading.(*ast.Leaf).Literal
	h := markdown.Render(md, html.NewRenderer(html.RendererOptions{}))
	s := fmt.Sprintf("{{define \"title\"}} Mic's Web - %s {{end}} {{define \"body\"}}\n%s\n{{end}}", h1, h)
	_, err = tmpl.Parse(string(s))
	if err != nil {
		return err
	}
	return nil

}
