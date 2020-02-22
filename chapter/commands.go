/*
 * This is for Chapter admins (usually server owners of their chapter Discord)
 */

package chapter

import (
	"fmt"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"
)

// PartyOnCommandHandler : Command Handler for !partyon
func PartyOnCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	handled := false
	if !strings.HasPrefix(message.Content, fmt.Sprintf("%spartyon ", config.Constants().GuildCommandPrefix)) {
		return handled // Not this command
	}

	handled = true
	channel, _ := session.Channel(message.ChannelID)
	if channel.Type != discordgo.ChannelTypeGuildText {
		return handled // We only respond to text channels, not DMs
	}

	if !utils.HasGuildPermission(session, message.GuildID, discordgo.PermissionManageRoles) {
		session.ChannelMessageSend(message.ChannelID, "Sorry, you don't have permission to grant roles.")
		return handled
	}

	guildData := config.FindGuildByID(message.GuildID)
	if guildData.AttendeeRoleID == "" {
		session.ChannelMessageSend(message.ChannelID, "We can't party!  I don't know what role to assign.")
		return handled
	}

	if len(message.Mentions) == 0 {
		session.ChannelMessageSend(message.ChannelID, "You have to tell me who is ready to party!")
		return handled
	}

	if utils.FindRole(message.GuildID, guildData.AttendeeRoleID) == nil {
		session.ChannelMessageSend(message.ChannelID, "Invalid Attendee Role Stored!")
		return handled
	}

	var responseMessage = "Party On "
	for _, mention := range message.Mentions {
		err := session.GuildMemberRoleAdd(message.GuildID, mention.ID, guildData.AttendeeRoleID)
		utils.Assert("Unable to add role!", err, false)
		responseMessage += "<@" + mention.ID + ">! "
	}
	responseMessage += config.Constants().PartyOnLink

	session.ChannelMessageSend(message.ChannelID, responseMessage)
	return handled
}

// ConfigResponseDMHandler : parses DMs from Chapter Heads to listen for answers to configuration questions.
func ConfigResponseDMHandler(user *discordgo.User, message *discordgo.MessageCreate) {
	session := config.Globals().Session

	// Validate our data
	ownerSetup, ok := linq.From(config.Globals().OwnerSetups).FirstWithT(func(os config.OwnerSetups) bool { return os.OwnerID == user.ID }).(config.OwnerSetups)
	if !ok {
		utils.LogErrorf("ConfigResponseDMHandler", "Unable to find OwnerSetup for user %s.  How did we get here?!?!", user.ID)
		return
	}

	if len(ownerSetup.GuildSetups) < 1 {
		utils.LogErrorf("ConfigResponseDMHandler", "Found OwnerSetup but no GuildSetups for user %s.  How did we get here?!?!", user.ID)
		return
	}

	guildSetup := ownerSetup.GuildSetups[0]
	guild, guildErr := session.Guild(guildSetup.GuildID)
	if guildErr != nil {
		utils.LogErrorf("ConfigResponseDMHandler", "Guild %s doesn't exist!", guildSetup.GuildID)
		return
	}

	guildData := config.FindGuildByID(guild.ID)

	// Handle the different steps
	switch guildSetup.SetupStep {
	case config.GuildSetupStepAnnouncementChannel:
		channel := validateChannelForGuild(guild, message.Content, "#announcements", user)
		if channel == nil {
			// validateChannelForGuild sends error messages, so there's nothing to do but bail.
			return
		}
		guildData.AnnounceChannelID = channel.ID
		break

	case config.GuildSetupStepAttendeeRole:
		return
	case config.GuildSetupStepChapterURL:
		return
	case config.GuildSetupStepComplete:
		utils.LogErrorf("ConfigResponseDMHandler", "Found GuildSetupStepComplete on a current GuildSetup! - %s", guild.ID)
		return
	case config.GuildSetupStepConfirmAuthorizedUser:
		if message.Content != "me!" {
			utils.SendDMToUser(user.ID, fmt.Sprintf("ERROR: Invalid.  Are you the correct person to configure server setttings for **%s**?  If you are, type **me!**. If not, type the following into one of your server channels: **!authorizeUser <User>** where <User> is a user on the server.", guild.Name))
			return
		}
		guildData.AuthorizedUserID = user.ID
		break
	case config.GuildSetupStepNewsletterURL:
		return
	case config.GuildSetupStepPastAttendeeRole:
		return
	}

	guildSetup.SetupStep = config.GetNextGuildSetupStep(guildData)
	PromptSetupStepByUser(user, guild, guildSetup.SetupStep)
}

func validateChannelForGuild(guild *discordgo.Guild, channelName string, suggestion string, target *discordgo.User) *discordgo.Channel {
	if !strings.HasPrefix(channelName, "#") {
		utils.SendDMToUser(target.ID, fmt.Sprintf("ERROR: Please enter a channel name! (ex: %s)", suggestion))
		return nil
	}

	channel := utils.FindChannelByName(guild, discordgo.ChannelTypeGuildText, channelName)
	if channel == nil {
		utils.SendDMToUser(target.ID, fmt.Sprintf("ERROR: Discord Server %s doesn't have a channel named %s", guild.Name, channelName))
	}

	return channel
}

func validateChannelForGuild(guild *discordgo.Guild, channelName string, suggestion string, target *discordgo.User) *discordgo.Channel {
	if !strings.HasPrefix(channelName, "#") {
		utils.SendDMToUser(target.ID, fmt.Sprintf("ERROR: Please enter a channel name! (ex: %s)", suggestion))
		return nil
	}

	channel := utils.FindChannelByName(guild, discordgo.ChannelTypeGuildText, channelName)
	if channel == nil {
		utils.SendDMToUser(target.ID, fmt.Sprintf("ERROR: Discord Server %s doesn't have a channel named %s", guild.Name, channelName))
	}

	return channel
}
