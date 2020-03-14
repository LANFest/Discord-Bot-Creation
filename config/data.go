package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/LANFest/Discord-Bot-Creation/utils"
	"github.com/ahmetb/go-linq"
	"github.com/bwmarrin/discordgo"
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

		newGuildSetup.SetupStep = GetNextGuildSetupStep(guildData)

		if newGuildSetup.SetupStep != GuildSetupStepComplete {
			UpsertGuildSetup(newGuildSetup, newOwnerID)
		}

		if Constants().DebugOutput {
			utils.LogErrorf("BuildOwnerSetupDataList", "Guild Setup for %s : %d", guild.Name, newGuildSetup.SetupStep)
		}
	}
}

// GetNextGuildSetupStep : Takes a GuildData and figures out the next setup step.
func GetNextGuildSetupStep(guildData *GuildData) GuildSetupStep {
	if guildData.AuthorizedUserID == "" {
		return GuildSetupStepConfirmAuthorizedUser
	} else if guildData.LanFestURL == "" {
		return GuildSetupStepChapterURL
	} else if guildData.NewsURL == "" {
		return GuildSetupStepNewsletterURL
	} else if guildData.AnnounceChannelID == "" {
		return GuildSetupStepAnnouncementChannel
	} else if guildData.LFGCategoryID == "" {
		return GuildSetupStepLFGCategory
	} else if guildData.AttendeeRoleID == "" {
		return GuildSetupStepAttendeeRole
	} else if guildData.PastAttendeeRoleID == "" {
		return GuildSetupStepPastAttendeeRole
	}
	return GuildSetupStepComplete
}

// UpsertGuildSetup : Updates the global Owner setups with the specified guildsetup
func UpsertGuildSetup(guildSetup GuildSetupData, ownerID string) {
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

// ReadConfig : Reads in the config data from the file and populates the supplied pointer
func ReadConfig() {
	file, readError := ioutil.ReadFile(Constants().ConfigFilePath)
	utils.Assert("Error reading Config Data", readError, false)

	parseError := json.Unmarshal(file, &Globals().GuildData)
	utils.Assert("Error parsing Config Data", parseError, false)
}

// WriteConfig : Writes the config data to disk
func WriteConfig() {
	file, _ := json.MarshalIndent(Globals().GuildData, "", " ")
	error := ioutil.WriteFile(Constants().ConfigFilePath, file, 0644)
	utils.Assert("Error writing config data!", error, false)
}

// IsOwner : Are you my daddy?
func IsOwner(user *discordgo.User) bool {
	return user.ID == Globals().Owner.ID
}
