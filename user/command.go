/*
 * This is for all users (Attendees mostly)
 */

package user

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"
)

// LFGCommandHandler : Command handler for !lfg
func LFGCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	handled := false
	if !strings.HasPrefix(message.Message.Content, fmt.Sprintf("%slfg ", config.Constants().GuildCommandPrefix)) {
		return handled
	}

	// Should we only let Attendees do it?

	// Should we only allow it during an event?

	handled = true // Yep, this is ours.
	commandRegex := regexp.MustCompile("^" + config.Constants().GuildCommandPrefix + "lfg \"(.+)\" (\\d+)$")
	commandArgs := commandRegex.FindStringSubmatch(message.Content)
	if len(commandArgs) < 3 {
		session.ChannelMessageSend(message.ChannelID, "Usage: !lfg \"<Game Name>\" <NumberOfPlayers>")
		return handled
	}

	gameNameRaw := commandArgs[1]
	capacityRaw := commandArgs[2]

	specialCharsRegex := regexp.MustCompile("[^a-zA-Z0-9]")
	gameName := strings.ToLower(specialCharsRegex.ReplaceAllString(gameNameRaw, ""))
	capacity, _ := strconv.Atoi(capacityRaw)
	if capacity < 2 {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Incorrect number of players! - %s", capacityRaw))
		return handled
	}

	guildModel := config.FindGuildByID(message.GuildID)
	guild, _ := session.Guild(guildModel.GuildID)

	rawCategory := linq.From(guild.Channels).FirstWithT(func(c *discordgo.Channel) bool { return c.ID == guildModel.LFGCategoryID })
	var category *discordgo.Channel
	if rawCategory == nil {
		// No category! Let's make one.
		if !utils.HasGuildPermission(session, guild.ID, discordgo.PermissionManageChannels) {
			session.ChannelMessageSend(message.ChannelID, "I can't create a channel! Please contact a server administrator.")
			utils.LPrint("Unable to create LFG Category: No permissions")
			return handled
		}

		newCategory, newCatErr := session.GuildChannelCreate(guild.ID, "LFG", discordgo.ChannelTypeGuildCategory)
		if newCatErr != nil {
			session.ChannelMessageSend(message.ChannelID, "Unable to create channel. Please contact a server administrator.")
			utils.LPrintf("Unable to create LFG Category: %s", newCatErr)
			return handled
		}

		guildModel.LFGCategoryID = newCategory.ID
		category = newCategory
	}

	newChannel, newChannelErr := CreateLFGChannel(session, category, gameName)
	if newChannelErr != nil {
		errMsg := fmt.Sprintf("Unable to create LFG Channel: %s", newChannelErr)
		session.ChannelMessageSend(message.ChannelID, errMsg)
		utils.LPrint(errMsg)
		return handled
	}

	lfgData := config.LFGData{ChannelID: newChannel.ID, Capacity: capacity, OwnerID: message.Author.ID, CreateDate: time.Now()}
	guildModel.LFGData = append(guildModel.LFGData, lfgData)
	config.WriteConfig()

	session.MessageReactionAdd(message.ChannelID, message.ID, ThumbsUpEmoji)
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("<@%s> has been added to the LFG queue in %s for the next 4 hours.", message.Message.Author.ID, newChannel.Name))
	pollMessage, _ := session.ChannelMessageSend(newChannel.ID, fmt.Sprintf("Hey <@%s> is looking for %s players for %s - click the :thumbsup: to join in!", message.Author.ID, capacityRaw, gameNameRaw))
	session.ChannelMessagePin(newChannel.ID, pollMessage.ID)
	session.MessageReactionAdd(newChannel.ID, pollMessage.ID, ThumbsUpEmoji)

	return handled
}
