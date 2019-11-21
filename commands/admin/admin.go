/*
 * This is for the Bot Owners themselves.
 */

package admin

import (
	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/bwmarrin/discordgo"
)

func WriteConfigCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if utils.IsOwner(message.Author) &&
		utils.IsDM(message.Message) &&
		message.Content == data.Constants().CommandPrefix+"writeConfig" {

		utils.WriteConfig()
		session.ChannelMessageSend(message.ChannelID, "Config Data Recorded")
	}
}

func ShutdownCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if utils.IsOwner(message.Author) &&
		utils.IsDM(message.Message) &&
		message.Content == data.Constants().CommandPrefix+"shutdown" {

		session.ChannelMessageSend(message.ChannelID, "Buh-bye!")
		utils.Shutdown(session)
	}
}
