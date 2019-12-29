package user

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"
)

func LFGChannelMessageReactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) bool {
	handled := false

	lfgData := utils.GetLFGDataForChannel(reaction.GuildID, reaction.ChannelID)
	if lfgData == nil {
		return handled
	}

	pinnedMessages, _ := session.ChannelMessagesPinned(reaction.ChannelID)
	if len(pinnedMessages) > 0 && reaction.MessageID != pinnedMessages[0].ID {
		return handled
	}

	handled = true // This is a reaction to the first pinned message of a LFG channel.  This is the correct handler.

	decodedEmoji, _ := url.QueryUnescape(ThumbsUpEmoji)
	if reaction.Emoji.Name != decodedEmoji {
		return handled // Not a thumbs up? Nothing to do here...
	}

	// Give me all real (non-Bot) users who are not the owner and have reacted with the thumbsup emoji and place into readyUsers.
	messageReactions, _ := session.MessageReactions(reaction.ChannelID, reaction.MessageID, decodedEmoji, 100) // 100 is the cap
	var readyUserMentions []string
	linq.From(messageReactions).
		WhereT(func(u *discordgo.User) bool { return !u.Bot && u.ID != lfgData.OwnerID }).
		SelectT(func(u *discordgo.User) string { return u.Mention() }).
		ToSlice(&readyUserMentions)

	// If we've already hit capacity and more reactions come in, we don't want to spam.
	if (lfgData.Capacity - 1) == len(readyUserMentions) { // Capacity includes the owner.
		// We found enough!
		var resultMessage string
		owner, ownerError := session.GuildMember(reaction.GuildID, lfgData.OwnerID)
		if ownerError != nil {
			resultMessage = fmt.Sprintf("Unable to find LFG Owner - %s", lfgData.OwnerID)
			utils.LPrintf("LFGChannelMessageReactionAdd: %s", resultMessage)
		} else {
			resultMessage = fmt.Sprintf("%s -- your group is ready!  Members: %s", owner.Mention(), strings.Join(readyUserMentions, " "))
		}
		session.ChannelMessageSend(reaction.ChannelID, resultMessage)
	}
	return handled
}
