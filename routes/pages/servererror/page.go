package servererror

import (
	"log/slog"
	"net/http"

	"github.com/mbaklor/maze/paths"
	"github.com/mbaklor/maze/routes"
)

func Handler(logger *slog.Logger, serverErr error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		info := routes.NewLayoutInfo("Server Error", paths.BasePathFromUrl(r.URL.Path))
		err := info.GenerateLinks()
		if err != nil {
			logger.Error("Failed to generate Base Path links! rendering header with Home", "error", err.Error())
			info.Links = []paths.Link{paths.NewLink("/", "Home")}
		}
		l := routes.Layout(view(serverErr.Error()), info)
		l.ServeHTTP(w, r)
	}
}
