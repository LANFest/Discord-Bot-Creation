package user

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/bwmarrin/discordgo"
)

func LFGCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {

	if strings.HasPrefix(message.Message.Content, data.Constants().CommandPrefix+"lfg ") {
		specialCharsRegex, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			session.MessageReactionAdd(message.Message.ChannelID, message.Message.ID, "632789111958274068")
		}
		gameSuggestion := strings.ToLower(specialCharsRegex.ReplaceAllString(strings.TrimPrefix(message.Message.Content, "!lfg "), ""))
		session.ChannelMessageSend(message.Message.ChannelID, fmt.Sprintf("%s has been added to the LFG queue for %s for the next 4 hours.", message.Message.Author.Username, gameSuggestion))
	}
}
