package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-plugin-bulk-invite/server/engine"
	"github.com/mattermost/mattermost-plugin-bulk-invite/server/perror"
)

const maxFileSizeKiloBytes = 256

func Init(handler *Handler, engine *engine.Engine) {
	handlersRouter := handler.Router.PathPrefix("/handlers").Subrouter()
	handlersRouter.HandleFunc(
		"/channel_bulk_add",
		checkAuthenticatedUser(injectEngine(handler.channelBulkAddHandler, engine)),
	).Methods("POST")
}

type bulkAddChannelPayload struct {
	ChannelID string           `json:"channel_id"`
	AddToTeam bool             `json:"add_to_team"`
	Users     []engine.AddUser `json:"users"`
}

func (bip *bulkAddChannelPayload) IsValid() *perror.PError {
	if bip.ChannelID == "" {
		return perror.NewPError(fmt.Errorf("missing channel_id"), "Channel ID is required.")
	}

	if len(bip.Users) == 0 {
		return perror.NewPError(fmt.Errorf("missing users"), "User list is empty.")
	}

	return nil
}

func (bip *bulkAddChannelPayload) FromRequest(r *http.Request) *perror.PError {
	f, h, err := r.FormFile("file")
	if f == nil {
		return perror.NewPError(fmt.Errorf("missing file"), "File is required.")
	}
	if err != nil {
		return perror.NewPError(err, "error parsing file")
	}

	if h.Size > maxFileSizeKiloBytes*1024 {
		return perror.NewPError(fmt.Errorf("file too large"), fmt.Sprintf("File is too large. Max file size is %dKB.", maxFileSizeKiloBytes))
	}

	if h.Header.Get("Content-Type") != "application/json" {
		return perror.NewPError(fmt.Errorf("invalid content type"), "Invalid file type, only JSON is supported")
	}

	if err := json.NewDecoder(f).Decode(&bip); err != nil {
		return perror.NewPError(err, "Error parsing submitted file")
	}

	bip.ChannelID = r.FormValue("channel_id")
	bip.AddToTeam = r.FormValue("add_to_team") == "true"

	return nil
}

func (h *Handler) channelBulkAddHandler(w http.ResponseWriter, r *http.Request, e *engine.Engine) {
	userID := getMattermostUserIDFromRequest(r)

	defer r.Body.Close()

	// Do not parse data in memory
	if err := r.ParseMultipartForm(0); err != nil {
		h.Logger.LogError("error parsing channel bulk add form", "err", err.Error())
		sendInternalServerError(w)
		return
	}

	var payload bulkAddChannelPayload

	if err := payload.FromRequest(r); err != nil {
		h.Logger.LogError("error parsing channel bulk add form payload", "err", err.Error())
		sendResponse(w,
			withHeader("Content-Type", "application/json"),
			withStatusCode(http.StatusBadRequest),
			withBody(err.AsJSON()),
		)
		return
	}

	if err := payload.IsValid(); err != nil {
		sendResponse(w, withStatusCode(http.StatusBadRequest), withBody(err.AsJSON()))
		return
	}

	engineConfig := &engine.Config{
		UserID:    userID,
		ChannelID: payload.ChannelID,
		AddToTeam: payload.AddToTeam,
		Users:     payload.Users,
	}

	if err := e.StartJob(context.TODO(), engineConfig); err != nil {
		sendResponse(w,
			withHeader("Content-Type", "application/json"),
			withStatusCode(http.StatusBadRequest),
			withBody(err.AsJSON()),
		)
		return
	}

	sendResponse(w, withStatusCode(http.StatusCreated), withBody(`{"message": "bulk add job started"}`))
}
