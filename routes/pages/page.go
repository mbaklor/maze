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

func Handler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filename, err := paths.ParseFileFromUrl(r.URL.Path, paths.FrontendPath("pages"))
		if err != nil {
			s := servererror.Handler(logger, err)
			w.WriteHeader(http.StatusInternalServerError)
			s.ServeHTTP(w, r)
			return
		}

		m, err := md.ParseMarkdownFile(filename)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				s := notfound.Handler(logger)
				w.WriteHeader(http.StatusNotFound)
				s.ServeHTTP(w, r)
				return
			}
			s := servererror.Handler(logger, err)
			s.ServeHTTP(w, r)
			return
		}
		info := routes.NewLayoutInfo(m.Info.Title, paths.BasePathFromUrl(r.URL.Path))
		err = info.GenerateLinks()
		if err != nil {
			logger.Error("Failed to generate Base Path links! rendering header with Home", "error", err.Error())
			info.Links = []paths.Link{paths.NewLink("/", "Home")}
		}
		l := routes.Layout(view(m), info)
		l.ServeHTTP(w, r)
	}
}
