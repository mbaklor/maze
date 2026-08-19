package notfound

import (
	"log/slog"
	"net/http"

	"github.com/mbaklor/website/paths"
	"github.com/mbaklor/website/routes"
)

func Handler(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := routes.NewLayoutInfo("Not Found", paths.BasePathFromUrl(r.URL.Path))
		err := info.GenerateLinks()
		if err != nil {
			logger.Error("Failed to generate Base Path links! rendering header with Home", "error", err.Error())
			info.Links = []paths.Link{paths.NewLink("/", "Home")}
		}
		l := routes.Layout(view(), info)
		l.ServeHTTP(w, r)
	}
}
