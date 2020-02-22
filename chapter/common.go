package chapter

import (
	"fmt"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/ahmetb/go-linq"
	"github.com/bwmarrin/discordgo"
)

// PromptSetupSteps : Checks for any necessary configuration steps that a particular owner needs to take.  If "", then do ALL owners.
func PromptSetupSteps(ownerID string) {
	session := config.Globals().Session

	var ownerSetups []config.OwnerSetups

	if ownerID == "" {
		// Do all of them.
		ownerSetups = config.Globals().OwnerSetups
	} else {
		// Just do the specified owner.
		linq.From(config.Globals().OwnerSetups).WhereT(func(os config.OwnerSetups) bool { return os.OwnerID == ownerID }).ToSlice(&ownerSetups)
	}

	for _, setup := range ownerSetups {

		user, userErr := session.User(setup.OwnerID)
		utils.Assert(fmt.Sprintf("Missing UserID in PromptSetupSteps! - %s", setup.OwnerID), userErr, false)
		if user == nil {
			continue
		}

		guildSetup := setup.GuildSetups[0]
		guild, guildErr := session.Guild(guildSetup.GuildID)
		utils.Assert(fmt.Sprintf("Missing GuildID in PromptSetupSteps! - %s", guild.ID), guildErr, false)
		if guild == nil {
			continue
		}

		if guildSetup.SetupStep != config.GuildSetupStepComplete {
			PromptSetupStepByUser(user, guild, guildSetup.SetupStep)
		}
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
	case config.GuildSetupStepLFGCategory:
		message = fmt.Sprintf("In **%s**, what category should new LFG channels be created in? (If you don't have one, you'll need to create a category for them to go into.)", guild.Name)
	case config.GuildSetupStepNewsletterURL:
		message = fmt.Sprintf("In **%s**, what is the signup URL for your chapter newsletter? (ex: https://tempuri.org/signup.asp or noNews! if you don't have one.)", guild.Name)
		break
	case config.GuildSetupStepPastAttendeeRole:
		message = fmt.Sprintf("In **%s**, what is the name of the role you would like to grant to previous attendees? (usually @PastAttendee)", guild.Name)
		break
	case config.GuildSetupStepComplete:
		message = fmt.Sprintf("Configuration complete for server **%s**!", guild.Name)
		break
	}

	utils.SendDMToUser(config.Globals().Session, user.ID, message)
}
