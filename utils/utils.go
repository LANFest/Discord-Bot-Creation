package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/bwmarrin/discordgo"
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

// ReadConfig : Reads in the config data from the file and populates the supplied pointer
func ReadConfig() {
	file, readError := ioutil.ReadFile(data.Constants().ConfigFilePath)
	Assert("Error reading Config Data", readError, false)

	parseError := json.Unmarshal(file, &data.Globals().GuildData)
	Assert("Error parsing Config Data", parseError, false)
}

// WriteConfig : Writes the config data to disk
func WriteConfig() {
	file, _ := json.MarshalIndent(data.Globals().GuildData, "", " ")
	error := ioutil.WriteFile(data.Constants().ConfigFilePath, file, 0644)
	Assert("Error writing config data!", error, false)
}

// FindRole : finds the requested Role within the list of Roles from the specified Guild
func FindRole(guildID string, roleID string) *discordgo.Role {
	tempGuild, _ := data.Globals().Session.Guild(guildID)
	for _, role := range tempGuild.Roles {
		if role.ID == roleID {
			return role
		}
	}
	return new(discordgo.Role)
}

// Assert : if error exists, panic.
func Assert(msg string, err error, shouldPanic bool) {
	if err != nil {
		LPrintf("%s: %+v", msg, err)
		if shouldPanic {
			panic(err)
		}
	}
}

func Shutdown(session *discordgo.Session) {
	LPrint("Shutting Down!")
	session.Logout()
	session.Close()
	os.Exit(0)
}

func IsOwner(user *discordgo.User) bool {
	return user.ID == data.Globals().Owner.ID
}

func IsDM(message *discordgo.Message) bool {
	channel, _ := data.Globals().Session.Channel(message.ChannelID)
	return channel.Type == discordgo.ChannelTypeDM
}

func LPrint(message string) {
	fmt.Print(message + "\n")
}

func LPrintf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

func Substring(string source, number start, number end) string {
	return string([]rune(source)[])
}