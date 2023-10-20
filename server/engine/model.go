package engine

import (
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
)

type AddUser struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
}

type bulkChannelAddResult struct {
	addedUsers  int
	addedToTeam int
	errorUsers  int

	notAddedGuest         int
	notAddedNonTeamMember int
}

func (bir *bulkChannelAddResult) NotAddedCount() int {
	return bir.notAddedGuest + bir.notAddedNonTeamMember
}

func (bir bulkChannelAddResult) String() string {
	return fmt.Sprintf("%d users were added. %d had errors (check logs) and %d were not added, %d were added to the team.", bir.addedUsers, bir.errorUsers, bir.NotAddedCount(), bir.addedToTeam)
}

func (bir bulkChannelAddResult) PrettyString() string {
	prettyString := "Results:\n"

	prettyString += fmt.Sprintf("- **Total users to add**: %d\n", bir.addedUsers)

	if bir.errorUsers > 0 {
		prettyString += fmt.Sprintf("- **Errors**: %d (check logs)\n", bir.errorUsers)
	}

	if bir.NotAddedCount() > 0 {
		prettyString += fmt.Sprintf("- **Not added**: %d\n", bir.NotAddedCount())

		if bir.notAddedGuest > 0 {
			prettyString += fmt.Sprintf("  - **Due to being a guest**: %d\n", bir.notAddedGuest)
		}

		if bir.notAddedNonTeamMember > 0 {
			prettyString += fmt.Sprintf("  - **Due to not being a team member**: %d\n", bir.notAddedNonTeamMember)
		}
	}

	if bir.addedToTeam > 0 {
		prettyString += fmt.Sprintf("- **Added to team**: %d\n", bir.addedToTeam)
	}

	return prettyString
}

type Config struct {
	// ChannelID the channel to add users to
	ChannelID string
	channel   *model.Channel

	// UserID stores the user ID that is inviting all users to the channel
	UserID string

	// Users are all the Users that require inviting to a channel
	Users []AddUser

	// AddToTeam add users to the team if they do not belong to it
	AddToTeam bool

	// AddGuests add guests to the team and channel if they do not belong to it
	AddGuests bool

	// inviteToWorkspace invite the users by email to the workspace if they are not registered
	// inviteToWorkspace bool
}
