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
	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"
)

func LFGCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if !strings.HasPrefix(message.Message.Content, data.Constants().CommandPrefix+"lfg ") {
		return
	}

	// Should we only let Attendees do it?

	// Should we only allow it during an event?

	commandRegex := regexp.MustCompile("^" + data.Constants().CommandPrefix + "lfg \"(.+)\" (\\d+)$")
	commandArgs := commandRegex.FindStringSubmatch(message.Content)
	if len(commandArgs) < 3 {
		session.ChannelMessageSend(message.ChannelID, "Usage: !lfg \"<Game Name>\" <NumberOfPlayers>")
		return
	}

	gameNameRaw := commandArgs[1]
	capacityRaw := commandArgs[2]

	specialCharsRegex := regexp.MustCompile("[^a-zA-Z0-9]")
	gameName := strings.ToLower(specialCharsRegex.ReplaceAllString(gameNameRaw, ""))
	capacity, _ := strconv.Atoi(capacityRaw)
	if capacity < 2 {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Incorrect number of players! - %s", capacityRaw))
		return
	}

	guildModel := utils.FindGuildByID(message.GuildID)
	guild, _ := session.Guild(guildModel.GuildID)

	rawCategory := linq.From(guild.Channels).FirstWithT(func(c *discordgo.Channel) bool { return c.ID == guildModel.LFGCategoryID })
	if rawCategory == nil {
		// No category! Let's make one.
		if !utils.HasGuildPermission(session, guild.ID, uint(data.PERMISSION_MANAGE_CHANNELS)) {
			session.ChannelMessageSend(message.ChannelID, "I can't create a channel! Please contact a server administrator.")
			utils.LPrint("Unable to create LFG Category: No permissions")
			return
		}

		newCategory, newCatErr := session.GuildChannelCreate(guild.ID, "LFG", discordgo.ChannelTypeGuildCategory)
		if newCatErr != nil {
			session.ChannelMessageSend(message.ChannelID, "Unable to create channel. Please contact a server administrator.")
			utils.LPrintf("Unable to create LFG Category: %s", newCatErr)
			return
		}

		guildModel.LFGCategoryID = newCategory.ID
	}

	createData := discordgo.GuildChannelCreateData{Name: "lfg_" + gameName, Type: discordgo.ChannelTypeGuildText, ParentID: guildModel.LFGCategoryID}

	newChannel, newChannelErr := session.GuildChannelCreateComplex(guild.ID, createData)
	if newChannelErr != nil {
		session.ChannelMessageSend(message.ChannelID, "I can't create the channel! Please contact a server administrator.")
		utils.LPrintf("Unable to create LFG Channel: %s", newChannelErr)
		return
	}

	lfgData := config.LFGData{ChannelID: newChannel.ID, Capacity: capacity, OwnerID: message.Author.ID, CreateDate: time.Now()}
	guildModel.LFGData = append(guildModel.LFGData, lfgData)
	utils.WriteConfig()

	session.MessageReactionAdd(message.ChannelID, message.ID, "%F0%9F%91%8D") // This is url-encoded emoji for thumbs-up
	session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("%s has been added to the LFG queue for %s for the next 4 hours.", message.Message.Author.Username, gameName))
	session.ChannelMessageSend(newChannel.ID, fmt.Sprintf("Hey @everyone! <@%s> is looking for %s players for %s - click the :thumbsup: to join in!", message.Author.ID, capacityRaw, gameName))

}
