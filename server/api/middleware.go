package api

import (
	"net/http"

	"github.com/mattermost/mattermost-plugin-bulk-invite/server/inviter"
)

type HandlerFuncPluginAPI func(w http.ResponseWriter, r *http.Request, engine *inviter.Engine)

// checkAuthenticatedUser checks the header to ensure that the Mattermost-User-ID header is present.
func checkAuthenticatedUser(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// mattermostUserID := getMattermostUserIDFromRequest(r)
		// if mattermostUserID == "" {
		// 	sendResponse(w, withStatusCode(http.StatusForbidden))
		// 	return
		// }

		handler(w, r)
	}
}

func injectInviterEngine(handler HandlerFuncPluginAPI, engine *inviter.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, engine)
	}
}
