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
		DMCommandPrefix:    "!",
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
	lfgData, ok := linq.From(guildData.LFGData).FirstWithT(func(lfg LFGData) bool { return lfg.ChannelID == channelID }).(LFGData)
	if ok {
		return &lfgData
	}

	return nil
}

// BuildOwnerSetupDataList : Builds a list of GuildSetupData objects
func BuildOwnerSetupDataList() {
	for _, guild := range Globals().Session.State.Guilds {
		guildData := FindGuildByID(guild.ID)

		var newGuildSetup GuildSetupData
		newGuildSetup.GuildID = guild.ID

		newOwnerID := guildData.AuthorizedUserID
		if guildData.AuthorizedUserID == "" {
			newOwnerID = guild.OwnerID
		}

		setupStep := GetNextGuildSetupStep(guildData)

		if setupStep != GuildSetupStepComplete {
			upsertGuildSetup(newGuildSetup, newOwnerID)
		}
	}
}

// GetNextGuildSetupStep : Takes a GuildData and figures out the next setup step.
func GetNextGuildSetupStep(guildData *GuildData) GuildSetupStep {
	if guildData.AuthorizedUserID == "" {
		return GuildSetupStepConfirmAuthorizedUser
	} else if guildData.LanFestURL == "" {
		return GuildSetupStepChapterURL
	} else if guildData.NewsURL != "" {
		return GuildSetupStepNewsletterURL
	} else if guildData.AnnounceChannelID != "" {
		return GuildSetupStepAnnouncementChannel
	} else if guildData.AttendeeRoleID != "" {
		return GuildSetupStepAttendeeRole
	} else if guildData.PastAttendeeRoleID != "" {
		return GuildSetupStepPastAttendeeRole
	}
	return GuildSetupStepComplete
}

func upsertGuildSetup(guildSetup GuildSetupData, ownerID string) {
	ownerSetup, ok := linq.From(Globals().OwnerSetups).FirstWithT(func(o OwnerSetups) bool { return o.OwnerID == ownerID }).(OwnerSetups)
	if ok {
		ownerSetup.GuildSetups = append(ownerSetup.GuildSetups, guildSetup)
	} else {
		var newOwnerSetup OwnerSetups
		newOwnerSetup.OwnerID = ownerID
		newOwnerSetup.GuildSetups = []GuildSetupData{guildSetup}
		Globals().OwnerSetups = append(Globals().OwnerSetups, newOwnerSetup)
	}
}
