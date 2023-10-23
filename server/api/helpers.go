package api

import "net/http"

const mattermostUserHeader = "Mattermost-User-ID"

// getMattermostUserIDFromRequest extracts the mattermost user ID from the Mattermost-User-ID header
func getMattermostUserIDFromRequest(r *http.Request) string {
	return r.Header.Get(mattermostUserHeader)
}
