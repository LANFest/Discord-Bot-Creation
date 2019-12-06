package data

import (
	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/bwmarrin/discordgo"
)

type ConstantsModel struct {
	CommandPrefix  string
	ConfigFilePath string
	PartyOnLink    string
	StatusMessage  string
}

func Constants() ConstantsModel {
	return ConstantsModel{
		CommandPrefix:  "!",
		ConfigFilePath: "configData.json",
		PartyOnLink:    "https://giphy.com/gifs/chuber-wayne-waynes-world-d3mlYwpf96kMuFjO",
		StatusMessage:  "with your heart <3",
	}
}

type GlobalDataModel struct {
	Bot             *discordgo.User
	GuildData       []config.GuildData
	Token           string
	CommandHandlers []interface{}
	Owner           *discordgo.User
	Session         *discordgo.Session
}

func Globals() *GlobalDataModel {
	return &data
}

var data GlobalDataModel
