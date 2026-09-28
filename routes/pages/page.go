package pages

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/mbaklor/maze/md"
	"github.com/mbaklor/maze/paths"
	"github.com/mbaklor/maze/routes"
	"github.com/mbaklor/maze/routes/pages/notfound"
	"github.com/mbaklor/maze/routes/pages/servererror"
)

type RouteHandler struct {
	title     string
	logger    *slog.Logger
	filepaths string
}

func NewRouteHandler(title, filepaths string, logger *slog.Logger) RouteHandler {
	return RouteHandler{title: title, filepaths: filepaths, logger: logger}
}

func (rh RouteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	filename, err := paths.ParseFileFromUrl(r.URL.Path, rh.filepaths)
	if err != nil {
		rh.handleServerError(err, w, r)
		return
	}

	m, err := md.ParseMarkdownFile(filename)
	if errors.Is(err, os.ErrNotExist) {
		rh.handleNotExist(w, r)
		return
	}
	if err != nil {
		rh.handleServerError(err, w, r)
		return
	}
	title := rh.title
	if m.Info.PageTitle != "" {
		title = m.Info.PageTitle + " - " + title
	}
	info := routes.NewLayoutInfo(title, paths.BasePathFromUrl(r.URL.Path), rh.filepaths)
	err = info.GenerateLinks()
	if err != nil {
		rh.logger.Error("Failed to generate Base Path links! rendering header with Home", "error", err.Error())
		info.Links = []paths.Link{paths.NewLink("/", "Home")}
	}
	l := routes.Layout(view(m), info)
	l.ServeHTTP(w, r)
}

func (rh RouteHandler) handleNotExist(w http.ResponseWriter, r *http.Request) {
	info := routes.NewLayoutInfo("Not Found", paths.BasePathFromUrl(r.URL.Path), rh.filepaths)
	err := info.GenerateLinks()
	if err != nil {
		rh.logger.Error("Failed to generate Base Path links! rendering header with Home", "error", err.Error())
		info.Links = []paths.Link{paths.NewLink("/", "Home")}
	}
	l := routes.Layout(notfound.View(), info)
	w.WriteHeader(http.StatusNotFound)
	l.ServeHTTP(w, r)
}

func (rh RouteHandler) handleServerError(serverErr error, w http.ResponseWriter, r *http.Request) {
	info := routes.NewLayoutInfo("Server Error", paths.BasePathFromUrl(r.URL.Path), rh.filepaths)
	err := info.GenerateLinks()
	if err != nil {
		rh.logger.Error("Failed to generate Base Path links! rendering header with Home", "error", err.Error())
		info.Links = []paths.Link{paths.NewLink("/", "Home")}
	}
	l := routes.Layout(servererror.View(serverErr.Error()), info)
	w.WriteHeader(http.StatusInternalServerError)
	l.ServeHTTP(w, r)
}
