package data

import (
	"os"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/bwmarrin/discordgo"
)

// ConstantsModel : Struct for hard-coded constants across the app
type ConstantsModel struct {
	GuildCommandPrefix string
	DMCommandPrefix    string
	ConfigFilePath     string
	PartyOnLink        string
	StatusMessage      string
	DebugOutput        bool
}

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

// GlobalDataModel : Struct for internal data used across the bot.
type GlobalDataModel struct {
	Bot                    *discordgo.User
	GuildData              []config.GuildData
	Token                  string
	DMCommandHandlers      []interface{}
	GuildCommandHandlers   []interface{}
	ReactionAddHandlers    []interface{}
	ReactionDeleteHandlers []interface{}
	Owner                  *discordgo.User
	Session                *discordgo.Session
}

// Globals : Singleton instance of GlobalDataModel
func Globals() *GlobalDataModel {
	return &data
}

var data GlobalDataModel
