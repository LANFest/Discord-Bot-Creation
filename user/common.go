package user

import (
	"errors"
	"fmt"

	"github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"
)

const ThumbsUpEmoji = "%F0%9F%91%8D"

func CreateLFGChannel(session *discordgo.Session, category *discordgo.Channel, gameName string) (*discordgo.Channel, error) {
	if category.Type != discordgo.ChannelTypeGuildCategory {
		return nil, errors.New(fmt.Sprintf("Unable to find LFG Category! - Received: %s", category.Name))
	}

	guild, _ := session.Guild(category.GuildID)
	var categoryChannels []*discordgo.Channel
	linq.From(guild.Channels).WhereT(func(c *discordgo.Channel) bool { return c.ParentID == category.ID }).ToSlice(&categoryChannels)

	for counter := 1; counter < 10; counter++ {
		channelName := newLFGChannelName(gameName, counter)
		if !linq.From(categoryChannels).AnyWithT(func(c *discordgo.Channel) bool { return c.Name == channelName }) {
			createData := discordgo.GuildChannelCreateData{Name: channelName, Type: discordgo.ChannelTypeGuildText, ParentID: category.ID}
			return session.GuildChannelCreateComplex(category.GuildID, createData)
		}
	}

	return nil, errors.New(fmt.Sprintf("Ran out of available channel space for %s -- Perhaps you could join an existing group.", gameName))
}

func newLFGChannelName(gameName string, number int) string {
	return fmt.Sprintf("lfg_%s_%d", gameName, number)
}
