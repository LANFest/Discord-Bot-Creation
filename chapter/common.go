package chapter

import "github.com/bwmarrin/discordgo"

import "github.com/LANFest/Discord-Bot-Creation/config"

// PromptSetupSteps : Iterates through all guilds that need setup and pushes the appropriate DMs.
func PromptSetupSteps() {
	for _, setup := range config.Globals().GuildSetups {
		var targetUserID string
		if setup.AuthorizedUserID != "" {
			targetUserID = setup.AuthorizedUserID
		} else {
			targetUserID = setup.OwnerID
		}

		switch getSetupStep(setup) {
		case config.GuildSetupStepConfirmAuthorizedUser:

		}
	}
}

// PromptSetupStepByUser : Sends the appropriate setup step to the user.
func PromptSetupStepByUser(user *discordgo.User) {

}

func getSetupStep(data config.GuildSetupData) config.GuildSetupStep {
	if data.AuthorizedUserID == "" {
		return config.GuildSetupStepConfirmAuthorizedUser
	}

	if data.ChapterURL == "" {
		return config.GuildSetupStepChapterURL
	}

	if data.NewsletterURL == "" {
		return config.GuildSetupStepNewsletterURL
	}

	if data.AnnouncementChannelID == "" {
		return config.GuildSetupStepAnnouncementChannel
	}

	if data.AttendeeRoleID == "" {
		return config.GuildSetupStepAttendeeRole
	}

	if data.PastAttendeeRoleID == "" {
		return config.GuildSetupStepPastAttendeeRole
	}

	return config.GuildSetupStepComplete
}
