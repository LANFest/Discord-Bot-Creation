package config

import (
	"os"

	"github.com/ahmetb/go-linq"
)

// Constants : Singleton Instance of ConstantsModel
func Constants() ConstantsModel {
	file, err := os.Stat("debug.txt")

	return ConstantsModel{
		GuildCommandPrefix: "!",
		DMCommandPrefix:    "?",
		ConfigFilePath:     "configData.json",
		PartyOnLink:        "https://giphy.com/gifs/chuber-wayne-waynes-world-d3mlYwpf96kMuFjO",
		StatusMessage:      "with your heart <3",
		DebugOutput:        !os.IsNotExist(err) && !file.IsDir(),
	}
}

// Globals : Singleton instance of GlobalDataModel
func Globals() *GlobalDataModel {
	return &data
}

var data GlobalDataModel

// FindGuildByID : Finds a guild based on the guildID
func FindGuildByID(targetGuildID string) *GuildData {
	for i := range Globals().GuildData {
		if Globals().GuildData[i].GuildID == targetGuildID {
			return &Globals().GuildData[i]
		}
	}

	guild := new(GuildData)
	guild.GuildID = targetGuildID
	Globals().GuildData = append(Globals().GuildData, *guild)

	return guild
}

// GetLFGDataForChannel : Gets any LFGData available for a channel (returns nil if none)
func GetLFGDataForChannel(guildID string, channelID string) *LFGData {
	guildData := FindGuildByID(guildID)
	lfgData, ok := linq.From(guildData.LFGData).WhereT(func(lfg LFGData) bool { return lfg.ChannelID == channelID }).First().(LFGData)
	if ok {
		return &lfgData
	}

	return nil
}

// BuildGuildSetupDataList : Builds a list of GuildSetupData objects
func BuildGuildSetupDataList() []GuildSetupData {
	var setups []GuildSetupData
	for _, guild := range Globals().Session.State.Guilds {
		guildData := FindGuildByID(guild.ID)
		var newGuildSetup GuildSetupData
		newGuildSetup.OwnerID = guild.OwnerID

		if guildData.AuthorizedUserID != "" {
			newGuildSetup.AuthorizedUserID = guildData.AuthorizedUserID
		} else {
			setups = append(setups, newGuildSetup)
			continue
		}

		if guildData.LanFestURL != "" {
			newGuildSetup.ChapterURL = guildData.LanFestURL
		} else {
			setups = append(setups, newGuildSetup)
			continue
		}

		if guildData.NewsURL != "" {
			newGuildSetup.NewsletterURL = guildData.NewsURL
		} else {
			setups = append(setups, newGuildSetup)
			continue
		}

		if guildData.AnnounceChannelID != "" {
			newGuildSetup.AnnouncementChannelID = guildData.AnnounceChannelID
		} else {
			setups = append(setups, newGuildSetup)
			continue
		}

		if guildData.AttendeeRoleID != "" {
			newGuildSetup.AttendeeRoleID = guildData.AttendeeRoleID
		} else {
			setups = append(setups, newGuildSetup)
			continue
		}

		if guildData.PastAttendeeRoleID != "" {
			newGuildSetup.PastAttendeeRoleID = guildData.PastAttendeeRoleID
		} else {
			setups = append(setups, newGuildSetup)
			continue
		}

		// All of the setup data is complete, no need to add to the list.
	}

	return setups
}
