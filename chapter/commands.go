/*
 * This is for Chapter admins (usually server owners of their chapter Discord)
 */

package chapter

import (
	"fmt"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/bwmarrin/discordgo"
)

func PartyOnCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	handled := false
	if !strings.HasPrefix(message.Content, fmt.Sprintf("%spartyon ", data.Constants().CommandPrefix)) {
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

	guildData := utils.FindGuildByID(message.GuildID)
	if guildData.AttendeeRoleID == "" {
		session.ChannelMessageSend(message.ChannelID, "We can't party!  I don't know what role to assign.")
		return handled
	}

	if len(message.Mentions) == 0 {
		session.ChannelMessageSend(message.ChannelID, "You have to tell me who is ready to party!")
		return handled
	}

	if (*utils.FindRole(message.GuildID, guildData.AttendeeRoleID)).ID == "" {
		session.ChannelMessageSend(message.ChannelID, "Invalid Attendee Role Stored!")
		return handled
	}

	var responseMessage = "Party On "
	for _, mention := range message.Mentions {
		err := session.GuildMemberRoleAdd(message.GuildID, mention.ID, guildData.AttendeeRoleID)
		utils.Assert("Unable to add role!", err, false)
		responseMessage += "<@" + mention.ID + ">! "
	}
	responseMessage += data.Constants().PartyOnLink

	session.ChannelMessageSend(message.ChannelID, responseMessage)
	return handled
}
