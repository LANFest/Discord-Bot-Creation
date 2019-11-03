package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/bwmarrin/discordgo"
)

var configFilePath = "configData.json"

// FindGuildByID : Finds a guild based on the guildID
func FindGuildByID(guilds *[]config.GuildData, targetGuildID string) *config.GuildData {
	for i := range *guilds {
		if (*guilds)[i].GuildID == targetGuildID {
			guild := (*guilds)[i]
			return &guild
		}
	}

	guild := new(config.GuildData)
	guild.GuildID = targetGuildID
	*guilds = append(*guilds, *guild)

	return guild

}

// ReadConfig : Reads in the config data from the file and populates the supplied pointer
func ReadConfig(guilds *[]config.GuildData) {
	file, readError := ioutil.ReadFile(configFilePath)
	Assert("Error reading Config Data", readError)

	parseError := json.Unmarshal(file, guilds)
	Assert("Error parsing Config Data", parseError)
}

// WriteConfig : Writes the config data to disk
func WriteConfig(guilds *[]config.GuildData) {
	file, _ := json.MarshalIndent(guilds, "", " ")
	error := ioutil.WriteFile(configFilePath, file, 0644)
	Assert("Error writing config data!", error)
}

// FindRole : finds the requested Role within the list of Roles from the specified Guild
func FindRole(session *discordgo.Session, guildID string, roleID string) *discordgo.Role {
	tempGuild, _ := session.Guild(guildID)
	for _, role := range tempGuild.Roles {
		if role.ID == roleID {
			return role
		}
	}
	return new(discordgo.Role)
}

// Assert : if error exists, panic.
func Assert(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}
