package engine

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-plugin-bulk-invite/server/kvstore"
	"github.com/mattermost/mattermost-plugin-bulk-invite/server/perror"
	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

type Engine struct {
	API plugin.API

	lockStore kvstore.LockStore

	// botUserID the bot user ID to set when sending messages
	botUserID string
}

func (e *Engine) onError(config *Config, _ error) {
	e.API.SendEphemeralPost(config.UserID, &model.Post{
		Message: "⚠️ Error bulk inviting users. Please check logs for more information.",
	})
}

func (e *Engine) checkPermissionsForUser(config *Config) *perror.PError {
	switch config.channel.Type {
	case model.ChannelTypePrivate:
		if !e.API.HasPermissionToChannel(config.UserID, config.ChannelID, model.PermissionManagePrivateChannelMembers) {
			return perror.NewPError(fmt.Errorf("insufficient_private_channel_permissions__add_user"), "You dont have permission to add users to this channel")
		}
	case model.ChannelTypeOpen:
		if !e.API.HasPermissionToChannel(config.UserID, config.ChannelID, model.PermissionManagePublicChannelMembers) {
			return perror.NewPError(fmt.Errorf("insufficient_public_channel_permissions__add_user"), "You dont have permission to add users to this channel")
		}
	case model.ChannelTypeGroup:
		if !e.API.HasPermissionToChannel(config.UserID, config.ChannelID, model.PermissionManageCustomGroupMembers) {
			return perror.NewPError(fmt.Errorf("insufficient_group_channel_permissions__add_user"), "You dont have permission to add users to this channel")
		}
	}

	if config.AddToTeam && !e.API.HasPermissionToTeam(config.UserID, config.channel.TeamId, model.PermissionAddUserToTeam) {
		return perror.NewPError(fmt.Errorf("insufficient_team_permissions__add_user"), "You dont have enough permissions to add users to this team")
	}

	// TODO: Invite users to teams
	// if !e.API.HasPermissionToTeam(config.UserID, config.channel.TeamId, model.PermissionInviteUser) {
	// 	return perror.NewPError(fmt.Errorf("insufficient_team_permissions__invite_user"), "You dont have permission to invite users to this team")
	// }

	return nil
}

func (e *Engine) StartJob(ctx context.Context, config *Config) *perror.PError {
	if e.lockStore.IsLocked(config.ChannelID) {
		return perror.NewPError(fmt.Errorf("channel_locked"), "A bulk operation is already running on this channel. Please wait until it finishes.")
	}

	var appErr *model.AppError
	config.channel, appErr = e.API.GetChannel(config.ChannelID)
	if appErr != nil {
		e.API.LogError("error getting channnel information", "channel_id", config.ChannelID, "err", appErr.Error())
		return perror.NewPError(
			fmt.Errorf("error getting channel: %w", appErr),
			fmt.Sprintf("Error getting channel information. Does channel `%s` exist?", config.ChannelID),
		)
	}

	if err := e.checkPermissionsForUser(config); err != nil {
		return perror.NewPError(
			fmt.Errorf("insufficient permissions: %w", err),
			"Insufficient permissions to add users to channel",
		)
	}

	if err := e.lockStore.Lock(config.ChannelID); err != nil {
		return perror.NewInternalServerPError(
			fmt.Errorf("error locking channel: %w", err),
		)
	}

	go e.start(ctx, config)

	return nil
}

func (e *Engine) start(_ context.Context, config *Config) {
	defer func() {
		if err := e.lockStore.Unlock(config.ChannelID); err != nil {
			e.API.LogError("error unlocking channel. channel will be automatically unlocked after ttl expired", "channel_id", config.ChannelID, "err", err.Error())
		}
	}()

	var appErr *model.AppError
	user, appErr := e.API.GetUser(config.UserID)
	if appErr != nil {
		e.API.LogError("error getting user information", "user_id", config.UserID, "err", appErr.Error())
		e.onError(config, appErr)
		return
	}

	if _, appErr = e.API.CreatePost(&model.Post{
		ChannelId: config.channel.Id,
		UserId:    e.botUserID,
		Message:   fmt.Sprintf("Starting bulk add of %d users (triggered by @%s)", len(config.Users), user.Username),
	}); appErr != nil {
		e.API.LogError("error creating initial post in channel", "channel_id", config.ChannelID, "err", appErr.Error())
	}

	result := e.addUsersToChannel(config)

	post, appErr := e.API.CreatePost(&model.Post{
		ChannelId: config.ChannelID,
		UserId:    e.botUserID,
		Message:   "Bulk add process finished.",
	})
	if appErr != nil {
		e.API.LogError("error creating result post in channel", "channel_id", config.ChannelID, "err", appErr.Error())
		e.onError(config, appErr)
	}

	if _, err := e.API.CreatePost(&model.Post{
		ChannelId: config.ChannelID,
		UserId:    e.botUserID,
		RootId:    post.Id,
		Message:   result.PrettyString(),
	}); err != nil {
		e.API.LogError("error creating threaded result post in channel", "channel_id", config.ChannelID, "err", err.Error())
		e.onError(config, err)
	}
}

func (e *Engine) addUserToChannelByUserID(config *Config, u AddUser, result *bulkChannelAddResult) error {
	if err := e.addToChannel(u.UserID, config, result); err != nil {
		return fmt.Errorf("error inviting user by ID: %w", err)
	}

	return nil
}

func (e *Engine) addUserToChannelByUsername(config *Config, u AddUser, result *bulkChannelAddResult) error {
	user, appErr := e.API.GetUserByUsername(u.Username)
	if appErr != nil {
		e.API.LogError("error getting user by username", "username", u.Username, "user_id", config.UserID, "channel_id", config.ChannelID, "err", appErr.Error())
		result.errorUsers++
		return fmt.Errorf("error getting user by username: %w", appErr)
	}

	if err := e.addToChannel(user.Id, config, result); err != nil {
		return fmt.Errorf("error inviting user by username: %w", err)
	}

	return nil
}

func (e *Engine) addToChannel(userID string, config *Config, result *bulkChannelAddResult) error {
	// Get user
	user, appErr := e.API.GetUser(userID)
	if appErr != nil {
		e.API.LogError("error getting user information", "add_user_id", userID, "trigger_user_id", config.UserID, "channel_id", config.ChannelID, "err", appErr.Error())
		result.errorUsers++
		return appErr
	}

	// Check if user is guest
	if user.IsGuest() && !config.AddGuests {
		e.API.LogInfo("not inviting guest user", "add_user_id", userID, "trigger_user_id", config.UserID, "channel_id", config.ChannelID)
		result.notAddedGuest++
		return nil
	}

	// Check team membership
	teamMembership, appErr := e.API.GetTeamMember(config.channel.TeamId, userID)
	if appErr != nil && appErr.StatusCode != http.StatusNotFound {
		e.API.LogError("error getting team membership information for user", "add_user_id", userID, "trigger_user_id", config.UserID, "channel_id", config.ChannelID, "team_id", config.channel.TeamId, "err", appErr.Error())
		result.errorUsers++
		return appErr
	}

	if teamMembership == nil {
		if config.AddToTeam {
			if _, createAppErr := e.API.CreateTeamMember(config.channel.TeamId, userID); createAppErr != nil {
				e.API.LogError("error creating team membership information for user", "add_user_id", userID, "trigger_user_id", config.UserID, "channel_id", config.ChannelID, "team_id", config.channel.TeamId, "err", appErr.Error())
				result.errorUsers++
				return createAppErr
			}
			result.addedToTeam++
		} else {
			e.API.LogInfo("not inviting member since it doesn't belong to the team", "add_user_id", userID, "trigger_user_id", config.UserID, "channel_id", config.ChannelID, "team_id", config.channel.TeamId)
			result.notAddedNonTeamMember++
			return nil
		}
	}

	if _, appErr := e.API.AddUserToChannel(config.ChannelID, userID, config.UserID); appErr != nil {
		result.errorUsers++
		e.API.LogError("error adding user to channel", "add_user_id", userID, "trigger_user_id", config.UserID, "channel_id", config.ChannelID, "err", appErr.Error())
		return appErr
	}

	return nil
}

func (e *Engine) addUsersToChannel(config *Config) bulkChannelAddResult {
	var result bulkChannelAddResult
	for _, u := range config.Users {
		if u.UserID != "" {
			if err := e.addUserToChannelByUserID(config, u, &result); err != nil {
				continue
			}
			result.addedUsers++
		}

		if u.Username != "" {
			if err := e.addUserToChannelByUsername(config, u, &result); err != nil {
				continue
			}
			result.addedUsers++
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
