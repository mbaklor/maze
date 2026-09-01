package routes

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/mbaklor/maze/paths"
)

type LayoutInfo struct {
	Title    string
	BasePath string
	Links    []paths.Link
}

func NewLayoutInfo(title string, basePath string) LayoutInfo {
	return LayoutInfo{Title: title, BasePath: basePath}
}

func (l *LayoutInfo) GenerateLinks() error {
	links, err := paths.GenerateBaseLinks(paths.FrontendPath("pages"))
	if err != nil {
		return err
	}
	l.Links = links
	return nil
}

func Layout(c templ.Component, info LayoutInfo) http.Handler {
	return templ.Handler(Page(c, info))
}
