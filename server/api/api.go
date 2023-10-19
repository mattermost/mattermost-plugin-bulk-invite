package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-plugin-bulk-invite/server/engine"
	"github.com/mattermost/mattermost-plugin-bulk-invite/server/perror"
)

func Init(handler *Handler, engine *engine.Engine) {
	apiV1Router := handler.Router.PathPrefix("/api/v1").Subrouter()
	apiV1Router.HandleFunc(
		"/bulk_invite",
		injectInviterEngine(handler.apiBulkInviteHandler, engine),
	).Methods("POST")

	handlersRouter := handler.Router.PathPrefix("/handlers").Subrouter()
	handlersRouter.HandleFunc(
		"/channel_bulk_invite",
		injectInviterEngine(handler.channelBulkInviteHandler, engine),
	).Methods("POST")
}

type bulkInvitePayload struct {
	ChannelID    string              `json:"channel_id"`
	InviteToTeam bool                `json:"invite_to_team"`
	InviteGuests bool                `json:"invite_guests"`
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
	// userID := getMattermostUserIDFromRequest(r)
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
			withBody(err.AsJSON()),
		)
		return
	}

	sendResponse(w, withStatusCode(http.StatusCreated))
}

func getPayloadFromFormData(r *http.Request) (*bulkInvitePayload, *perror.PError) {
	var payload bulkInvitePayload

	f, h, err := r.FormFile("file")
	if err != nil {
		return nil, perror.NewPError(err, "error parsing file")
	}

	if h.Header.Get("Content-Type") != "application/json" {
		return nil, perror.NewPError(fmt.Errorf("invalid content type"), "Invalid file type, only JSON is supported")
	}

	if err := json.NewDecoder(f).Decode(&payload); err != nil {
		return nil, perror.NewPError(err, "Error parsing submitted file")
	}

	payload.ChannelID = r.FormValue("channel_id")
	payload.InviteToTeam = r.FormValue("invite_to_team") == "true"

	return &payload, nil
}

func (h *Handler) channelBulkInviteHandler(w http.ResponseWriter, r *http.Request, e *engine.Engine) {
	userID := getMattermostUserIDFromRequest(r)

	defer r.Body.Close()

	// Do not parse data in memory
	if err := r.ParseMultipartForm(0); err != nil {
		h.Logger.LogError("error parsing channel bulk invite form", "err", err.Error())
		sendInternalServerError(w)
		return
	}

	payload, err := getPayloadFromFormData(r)
	if err != nil {
		h.Logger.LogError("error parsing channel bulk invite form payload", "err", err.Error())
		sendResponse(w,
			withHeader("Content-Type", "application/json"),
			withStatusCode(http.StatusBadRequest),
			withBody(err.AsJSON()),
		)
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

	sendResponse(w, withStatusCode(http.StatusCreated), withBody(`{"message": "bulk invite job started"}`))
}
