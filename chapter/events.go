package chapter

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/araddon/dateparse"
	"github.com/bwmarrin/discordgo"
)

// CreateEventCommandHandler Responds to event creation messages.
func CreateEventCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	handled := false
	responseMessage := ""
	usageMessage := "Usage: !event \"<Event Name>\" on <Date> @/at <Time>"
	cmdPrefix := fmt.Sprintf("%sevent ", data.Constants().CommandPrefix)
	if !strings.HasPrefix(message.Content, cmdPrefix) {
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

	var event config.LANPartyData
	name := strings.Split(strings.TrimSpace(strings.TrimPrefix(message.Content, cmdPrefix)), " ")
	if len(name) == 0 {
		session.ChannelMessageSend(message.ChannelID, usageMessage)
		return handled
	}

	reDate := regexp.MustCompile(`on (.+?) @`)
	date := time.Now()
	dates := reDate.FindStringSubmatch(message.Content)
	if (dates == nil) || (len(dates) < 2) {
		session.ChannelMessageSend(message.ChannelID, "Error: Cannot find date.\n"+usageMessage)
		return handled
	}
	date, _ = dateparse.ParseLocal(strings.TrimSpace(dates[1]))
	if !date.After(time.Now()) {
		session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Error: Date %s from %q doesn't seem to be in the future.\n%s", date.String(), dates[1], usageMessage))
		return handled
	}

	event.ActivateDate = time.Now()
	event.Capacity = 0
	event.LastAnnounceDate = time.Now()
	event.Name = name[0]
	event.StartDate = date
	event.TicketDate = time.Now()
	event.TicketURL = "https://www.lanfest.com/events"

	countdown := int(event.StartDate.Sub(time.Now()).Hours() / 24)
	countdownChannel, _ := session.GuildChannelCreate(message.GuildID, fmt.Sprintf("%s: %ddays", event.Name, countdown), discordgo.ChannelTypeGuildVoice)
	event.CountdownChannelID = countdownChannel.ID

	data.Globals().GuildData[0].PastLANData = append(data.Globals().GuildData[0].PastLANData, data.Globals().GuildData[0].NextLANData)
	data.Globals().GuildData[0].NextLANData = event
	utils.WriteConfig()
	responseMessage = fmt.Sprintf("New event \"%s\" created for %s@%s.", name, date, "time")
	session.ChannelMessageSend(message.ChannelID, responseMessage)
	return handled
}
