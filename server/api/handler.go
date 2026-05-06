package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-plugin-bulk-invite/server/mattermost"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type Handler struct {
	*mux.Router

	Logger mattermost.LoggerAPI
}

func NewHandler(pluginAPI plugin.API) *Handler {
	h := &Handler{
		Router: mux.NewRouter(),
		Logger: pluginAPI,
	}

	h.HandleFunc("{anything:.*}", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
	return h
}
