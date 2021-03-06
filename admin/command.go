/*
 * This is for the Bot Owners themselves.
 */

package admin

import (
	"fmt"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/bwmarrin/discordgo"
)

// WriteConfigDMCommandHandler : Command handler for !writeConfig
func WriteConfigDMCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	if message.Content != fmt.Sprintf("%swriteConfig", config.Constants().DMCommandPrefix) {
		return false // wrong handler
	}

	if config.IsOwner(message.Author) && utils.IsDM(session, message.Message) {
		config.WriteConfig()
		session.ChannelMessageSend(message.ChannelID, "Config Data Recorded")
	}
	return true
}

// ShutdownDMCommandHandler : Command handler for !shutdown
func ShutdownDMCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	if message.Content != fmt.Sprintf("%sshutdown", config.Constants().DMCommandPrefix) {
		return false // wrong handler
	}

	if config.IsOwner(message.Author) && utils.IsDM(session, message.Message) {
		session.ChannelMessageSend(message.ChannelID, "Buh-bye!")
		utils.Shutdown(session)
	}
	return true
}
