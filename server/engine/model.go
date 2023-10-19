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
	invitedUsers int
	addedToTeam  int
	errorUsers   int

	notInvitedGuest         int
	notInvitedNonTeamMember int
}

func (bir *bulkInviteResult) NotInvitedCount() int {
	return bir.notInvitedGuest + bir.notInvitedNonTeamMember
}

func (bir bulkInviteResult) String() string {
	return fmt.Sprintf("%d users were added. %d had errors (check logs) and %d were not invited, %d were added to the team.", bir.invitedUsers, bir.errorUsers, bir.NotInvitedCount(), bir.addedToTeam)
}

func (bir bulkInviteResult) PrettyString() string {
	prettyString := "Results:\n"

	prettyString += fmt.Sprintf("- **Total users to invite**: %d\n", bir.invitedUsers)

	if bir.errorUsers > 0 {
		prettyString += fmt.Sprintf("- **Errors**: %d (check logs)\n", bir.errorUsers)
	}

	if bir.NotInvitedCount() > 0 {
		prettyString += fmt.Sprintf("- **Not invited**: %d\n", bir.NotInvitedCount())

		if bir.notInvitedGuest > 0 {
			prettyString += fmt.Sprintf("  - **Due to being a guest**: %d\n", bir.notInvitedGuest)
		}

		if bir.notInvitedNonTeamMember > 0 {
			prettyString += fmt.Sprintf("  - **Due to not being a team member**: %d\n", bir.notInvitedNonTeamMember)
		}
	}

	if bir.addedToTeam > 0 {
		prettyString += fmt.Sprintf("- **Added to team**: %d\n", bir.addedToTeam)
	}

	return prettyString
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
