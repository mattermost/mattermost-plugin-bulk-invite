package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-plugin-bulk-invite/server/engine"
)

func Init(handler *Handler, engine *engine.Engine) {
	apiV1Router := handler.Router.PathPrefix("/api/v1").Subrouter()
	apiV1Router.HandleFunc(
		"/bulk_invite",
		checkAuthenticatedUser(injectInviterEngine(handler.apiBulkInviteHandler, engine)),
	).Methods("POST")
}

type bulkInvitePayload struct {
	ChannelID    string              `json:"channel_id"`
	InviteToTeam bool                `json:"invite_to_team"`
	Users        []engine.InviteUser `json:"users"`
}

func (bip *bulkInvitePayload) IsValid() error {
	if bip.ChannelID == "" {
		return fmt.Errorf("Missing channel ID")
	}

	if len(bip.Users) == 0 {
		return fmt.Errorf("Missing users")
	}

	return nil
}

func (h *Handler) apiBulkInviteHandler(w http.ResponseWriter, r *http.Request, e *engine.Engine) {
	// userID := GetMattermostUserIDFromRequest(r)
	userID := "bfswryyw67ntubga5omo9j5e1o"

	defer r.Body.Close()

	var payload bulkInvitePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		h.Logger.LogError("error parsing payload", "err", err.Error())
		sendResponse(w, withStatusCode(http.StatusBadRequest), withBody(`{"error": "%s"}`, err.Error()))
		return
	}

	if err := payload.IsValid(); err != nil {
		sendResponse(w, withStatusCode(http.StatusBadRequest), withBody(`{"error": "%s"}`, err.Error()))
		return
	}

	engineConfig := &engine.Config{
		UserID:       userID,
		ChannelID:    payload.ChannelID,
		InviteToTeam: payload.InviteToTeam,
		Users:        payload.Users,
	}

	if err := e.StartJob(context.TODO(), engineConfig); err != nil {
		sendResponse(w,
			withHeader("Content-Type", "application/json"),
			withStatusCode(http.StatusBadRequest),
			withBody(`{"error": "%s"}`, err.Message()),
		)
		return
	}

	sendResponse(w, withStatusCode(http.StatusCreated))
}
