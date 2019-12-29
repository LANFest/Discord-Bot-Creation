/*
 * This file is for data access helpers.
 */

package utils

import (
	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/ahmetb/go-linq/v3"
)

// FindGuildByID : Finds a guild based on the guildID
func FindGuildByID(targetGuildID string) *config.GuildData {
	for i := range data.Globals().GuildData {
		if data.Globals().GuildData[i].GuildID == targetGuildID {
			return &data.Globals().GuildData[i]
		}
	}

	guild := new(config.GuildData)
	guild.GuildID = targetGuildID
	data.Globals().GuildData = append(data.Globals().GuildData, *guild)

	return guild
}

// GetLFGDataForChannel : Gets any LFGData available for a channel (returns nil if none)
func GetLFGDataForChannel(guildID string, channelID string) *config.LFGData {
	guildData := FindGuildByID(guildID)
	lfgData, ok := linq.From(guildData.LFGData).WhereT(func(lfg config.LFGData) bool { return lfg.ChannelID == channelID }).First().(config.LFGData)
	if ok {
		return &lfgData
	}

	return nil
}
