package inviter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-plugin-bulk-invite/server/kvstore"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type bulkInviteResult struct {
	invitedUsers    int
	addedToTeam     int
	errorUsers      int
	notInvitedUsers int
}

func (bir bulkInviteResult) String() string {
	return fmt.Sprintf("%d invited: %d with errors and %d not invited. %d were added to the team.", bir.invitedUsers, bir.errorUsers, bir.notInvitedUsers, bir.addedToTeam)
}

type Config struct {
	// ChannelID the channel to invite users to
	ChannelID string
	channel   *model.Channel

	// UserID stores the user ID that is inviting all users to the channel
	UserID string

	// Users are all the Users that require inviting to a channel
	Users []InviteUser

	// InviteToTeam invite the users to the team if they do not belong to it
	InviteToTeam bool

	// inviteToWorkspace invite the users by email to the workspace if they are not registered
	// inviteToWorkspace bool
}

type Engine struct {
	API plugin.API

	lockStore kvstore.LockStore

	// botUserID the bot user ID to set when sending messages
	botUserID string
}

func (e *Engine) onError(config *Config, err error) {
	e.API.SendEphemeralPost(config.UserID, &model.Post{
		Message: fmt.Sprintf("Error bulk inviting users. Error provided: \n ```\n%s\n```.\n\n Please check logs for more information.", err.Error()),
	})
}

func (e *Engine) StartJob(ctx context.Context, config *Config) error {
	if e.lockStore.IsLocked(config.ChannelID) {
		return fmt.Errorf("a bulk invite operation is already running on this channel")
	}

	if err := e.lockStore.Lock(config.ChannelID); err != nil {
		return fmt.Errorf("internal error locking the operation")
	}

	go e.start(ctx, config)

	return nil
}

func (e *Engine) start(_ context.Context, config *Config) {
	var appErr *model.AppError
	config.channel, appErr = e.API.GetChannel(config.ChannelID)
	if appErr != nil {
		e.API.LogError("error getting channnel information", "channel_id", config.ChannelID, "err", appErr.Error())
		e.onError(config, appErr)
		return
	}

	user, appErr := e.API.GetUser(config.UserID)
	if appErr != nil {
		e.API.LogError("error getting user information", "user_id", config.UserID, "err", appErr.Error())
		e.onError(config, appErr)
		return
	}

	if _, appErr := e.API.CreatePost(&model.Post{
		ChannelId: config.channel.Id,
		UserId:    e.botUserID,
		Message:   fmt.Sprintf("Starting bulk invite of %d users (triggered by @%s)", len(config.Users), user.Username),
	}); appErr != nil {
		e.API.LogError("error creating initial post in channel", "channel_id", config.ChannelID, "err", appErr.Error())
	}

	result := e.inviteUsers(config)

	if _, appErr := e.API.CreatePost(&model.Post{
		ChannelId: config.ChannelID,
		UserId:    e.botUserID,
		Message:   fmt.Sprintf("Finished bulk inviting users. %s", result),
	}); appErr != nil {
		e.API.LogError("error creating result post in channel", "channel_id", config.ChannelID, "err", appErr.Error())
		e.onError(config, appErr)
	}
}

func (e *Engine) inviteMattermostByUserID(config *Config, invitee InviteUser, result *bulkInviteResult) error {
	if err := e.invite(invitee.UserID, config, result); err != nil {
		return fmt.Errorf("error inviting user by ID: %w", err)
	}

	return nil
}

func (e *Engine) inviteMattermostByUsername(config *Config, invitee InviteUser, result *bulkInviteResult) error {
	user, appErr := e.API.GetUserByUsername(invitee.Username)
	if appErr != nil {
		e.API.LogError("error getting user by username", "username", invitee.Username, "user_id", config.UserID, "channel_id", config.ChannelID, "err", appErr.Error())
		result.errorUsers++
		return fmt.Errorf("error getting user by username: %w", appErr)
	}

	if err := e.invite(user.Id, config, result); err != nil {
		return fmt.Errorf("error inviting user by username: %w", err)
	}

	return nil
}

func (e *Engine) invite(userID string, config *Config, result *bulkInviteResult) error {
	// Check team membership
	teamMembership, appErr := e.API.GetTeamMember(config.channel.TeamId, userID)
	if appErr != nil && appErr.StatusCode != http.StatusNotFound {
		e.API.LogError("error getting team membership information for user", "invitee_user_id", userID, "user_id", config.UserID, "channel_id", config.ChannelID, "team_id", config.channel.TeamId, "err", appErr.Error())
		result.errorUsers++
		return appErr
	}

	if teamMembership == nil {
		if config.InviteToTeam {
			if _, createAppErr := e.API.CreateTeamMember(config.channel.TeamId, userID); createAppErr != nil {
				e.API.LogError("error creating team membership information for user", "invitee_user_id", userID, "user_id", config.UserID, "channel_id", config.ChannelID, "team_id", config.channel.TeamId, "err", appErr.Error())
				result.errorUsers++
				return createAppErr
			}
			result.addedToTeam++
		} else {
			e.API.LogInfo("not inviting member since it doesn't belong to the team", "invitee_user_id", userID, "user_id", config.UserID, "channel_id", config.ChannelID, "team_id", config.channel.TeamId)
			result.notInvitedUsers++
			return fmt.Errorf("not invited")
		}
	}

	if _, appErr := e.API.AddUserToChannel(config.ChannelID, userID, config.UserID); appErr != nil {
		result.errorUsers++
		e.API.LogError("error adding user to channel", "invitee_user_id", userID, "user_id", config.UserID, "channel_id", config.ChannelID, "err", appErr.Error())
		return appErr
	}

	return nil
}

func (e *Engine) inviteUsers(config *Config) bulkInviteResult {
	var result bulkInviteResult
	for _, invitee := range config.Users {
		if invitee.UserID != "" {
			if err := e.inviteMattermostByUserID(config, invitee, &result); err != nil {
				continue
			}
			result.invitedUsers++
		}

		if invitee.Username != "" {
			if err := e.inviteMattermostByUsername(config, invitee, &result); err != nil {
				continue
			}
			result.invitedUsers++
		}
	}

	return result
}

func NewEngine(pluginAPI plugin.API, lockStore kvstore.LockStore, botUserID string) *Engine {
	return &Engine{
		API:       pluginAPI,
		lockStore: lockStore,
		botUserID: botUserID,
	}
}
