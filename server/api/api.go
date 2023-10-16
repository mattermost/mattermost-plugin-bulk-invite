package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/mattermost/mattermost-plugin-bulk-invite/server/inviter"
)

func Init(handler *Handler, engine *inviter.Engine) {
	apiV1Router := handler.Router.PathPrefix("/api/v1").Subrouter()
	apiV1Router.HandleFunc(
		"/bulk_invite",
		checkAuthenticatedUser(injectInviterEngine(handler.apiBulkInviteHandler, engine)),
	).Methods("POST")
}

type bulkInvitePayload struct {
	ChannelID    string               `json:"channel_id"`
	InviteToTeam bool                 `json:"invite_to_team"`
	Users        []inviter.InviteUser `json:"users"`
}

func (h *Handler) apiBulkInviteHandler(w http.ResponseWriter, r *http.Request, engine *inviter.Engine) {
	// userID := GetMattermostUserIDFromRequest(r)
	userID := "bfswryyw67ntubga5omo9j5e1o"

	defer r.Body.Close()

	var payload bulkInvitePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.Logger.LogError("error parsing payload", "err", err.Error())
		_, _ = w.Write([]byte(err.Error()))
		sendResponse(w, withStatusCode(http.StatusBadRequest))
		return
	}

	inviterConfig := &inviter.Config{
		UserID:       userID,
		ChannelID:    payload.ChannelID,
		InviteToTeam: payload.InviteToTeam,
		Users:        payload.Users,
	}

	if err := engine.StartJob(context.TODO(), inviterConfig); err != nil {
		sendResponse(w,
			withHeader("Content-Type", "application/json"),
			withStatusCode(http.StatusBadRequest),
			withBody(`{"error": "%s"}`, err.Error()),
		)
		return
	}

	sendResponse(w, withStatusCode(http.StatusCreated))
}
