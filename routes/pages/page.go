package pages

import (
	"log/slog"
	"net/http"

	"github.com/mbaklor/website/md"
	"github.com/mbaklor/website/paths"
	"github.com/mbaklor/website/routes"
)

func Handler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename, err := paths.ParseFileFromUrl(r.URL.Path, paths.FrontendPath("pages"))
		if err != nil {
			return
		}

		m, err := md.ParseMarkdownFile(filename)
		if err != nil {
			return
		}
		info := routes.NewLayoutInfo(m.Info.Title, paths.BasePathFromUrl(r.URL.Path))
		err = info.GenerateLinks()
		if err != nil {
			return
		}
		l := routes.Layout(view(m), info)
		l.ServeHTTP(w, r)
	}
}
