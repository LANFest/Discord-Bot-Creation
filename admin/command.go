/*
 * This is for the Bot Owners themselves.
 */

package admin

import (
	"fmt"

	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/bwmarrin/discordgo"
)

// WriteConfigCommandHandler : Command handler for !writeConfig
func WriteConfigCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	if message.Content != fmt.Sprintf("%swriteConfig", data.Constants().CommandPrefix) {
		return false // wrong handler
	}

	if utils.IsOwner(message.Author) && utils.IsDM(message.Message) {
		utils.WriteConfig()
		session.ChannelMessageSend(message.ChannelID, "Config Data Recorded")
	}
	return true
}

// ShutdownCommandHandler : Command handler for !shutdown
func ShutdownCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	if message.Content != fmt.Sprintf("%sshutdown", data.Constants().CommandPrefix) {
		return false // wrong handler
	}

	if utils.IsOwner(message.Author) && utils.IsDM(message.Message) {
		session.ChannelMessageSend(message.ChannelID, "Buh-bye!")
		utils.Shutdown(session)
	}
	return true
}
