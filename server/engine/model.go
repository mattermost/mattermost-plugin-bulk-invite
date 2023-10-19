package engine

import (
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
)

type InviteUser struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

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

	// InviteGuests invite the guests to the team and channel if they do not belong to it
	InviteGuests bool

	// inviteToWorkspace invite the users by email to the workspace if they are not registered
	// inviteToWorkspace bool
}
