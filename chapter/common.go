package chapter

import (
	"fmt"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/bwmarrin/discordgo"
)

// PromptSetupSteps : Iterates through all guilds that need setup and pushes the appropriate DMs.
func PromptSetupSteps() {
	session := config.Globals().Session

	for _, setup := range config.Globals().OwnerSetups {

		user, userErr := session.User(setup.OwnerID)
		utils.Assert("Missing UserID in PromptSetupSteps! - "+setup.OwnerID, userErr, false)
		if user == nil {
			continue
		}

		guildSetup := setup.GuildSetups[0]
		guild, guildErr := session.Guild(guildSetup.GuildID)
		utils.Assert("Missing GuildID in PromptSetupSteps! - "+guild.ID, guildErr, false)
		if guild == nil {
			continue
		}

		PromptSetupStepByUser(user, guild, guildSetup.SetupStep)
	}
}

// PromptSetupStepByUser : Sends the appropriate setup step to the user.
func PromptSetupStepByUser(user *discordgo.User, guild *discordgo.Guild, step config.GuildSetupStep) {
	message := ""

	switch step {
	case config.GuildSetupStepAnnouncementChannel:
		message = fmt.Sprintf("In **%s**, what is the name of the channel you would like to make public announcements in? (usually #announcements)", guild.Name)
		break
	case config.GuildSetupStepAttendeeRole:
		message = fmt.Sprintf("In **%s**, what is the name of the role you would like to grant to current attendees? (usually @Attendee)", guild.Name)
		break
	case config.GuildSetupStepChapterURL:
		message = fmt.Sprintf("In **%s**, what is the URL for your LanFest chapter? (ex: https://lanfest.com/emeraldcity)", guild.Name)
		break
	case config.GuildSetupStepConfirmAuthorizedUser:
		message = fmt.Sprintf("We need your help to set up correctly for **%s**!\nIf you can help, type **me!**.  If you are unable to help us set up, type the following into one of your server channels: **!authorizeUser <User>** where <User> is a user on the server.", guild.Name)
		break
	case config.GuildSetupStepNewsletterURL:
		message = fmt.Sprintf("In **%s**, what is the signup URL for your chapter newsletter? (ex: https://tempuri.org/signup.asp or noNews! if you don't have one.)", guild.Name)
		break
	case config.GuildSetupStepPastAttendeeRole:
		message = fmt.Sprintf("In **%s**, what is the name of the role you would like to grant to previous attendees? (usually @PastAttendee)", guild.Name)
		break
	}

	utils.SendDMToUser(user.ID, message)
}
