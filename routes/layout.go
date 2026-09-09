package routes

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/mbaklor/maze/paths"
)

type LayoutInfo struct {
	Title        string
	BasePath     string
	FrontendPath string
	Links        []paths.Link
}

func NewLayoutInfo(title string, basePath string, frontendPath string) LayoutInfo {
	return LayoutInfo{Title: title, BasePath: basePath, FrontendPath: frontendPath}
}

func (l *LayoutInfo) GenerateLinks() error {
	links, err := paths.GenerateBaseLinks(l.FrontendPath)
	if err != nil {
		return err
	}
	l.Links = links
	return nil
}

func Layout(c templ.Component, info LayoutInfo) http.Handler {
	return templ.Handler(Page(c, info))
}
