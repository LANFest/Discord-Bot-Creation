/*
 * This is for Chapter admins (usually server owners of their chapter Discord)
 */

package chapter

import (
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/bwmarrin/discordgo"
)

func PartyOnCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Are we in a channel?
	if strings.HasPrefix(message.Content, data.Constants().CommandPrefix+"partyon ") {
		guild := utils.FindGuildByID(message.GuildID)
		if guild.AttendeeRoleID == "" {
			session.ChannelMessageSend(message.ChannelID, "We can't party!  I don't know what role to assign.")
			return
		}

		if len(message.Mentions) == 0 {
			session.ChannelMessageSend(message.ChannelID, "You have to tell me who is ready to party!")
			return
		}

		if message.ChannelID == "" {
			session.ChannelMessageSend(message.ChannelID, "I can only do this in a server!")
		}

		if (*utils.FindRole(message.GuildID, guild.AttendeeRoleID)).ID == "" {
			session.ChannelMessageSend(message.ChannelID, "Invalid Attendee Role Stored!")
			return
		}

		var responseMessage = "Party On "
		for _, mention := range message.Mentions {
			err := session.GuildMemberRoleAdd(message.GuildID, mention.ID, guild.AttendeeRoleID)
			utils.Assert("Unable to add role!", err, false)
			responseMessage += "<@" + mention.ID + ">! "
		}
		responseMessage += data.Constants().PartyOnLink

		session.ChannelMessageSend(message.ChannelID, responseMessage)
	}
}
