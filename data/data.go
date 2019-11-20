package data

import (
	"os"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/bwmarrin/discordgo"
)

type ConstantsModel struct {
	CommandPrefix  string
	ConfigFilePath string
	PartyOnLink    string
	StatusMessage  string
	DebugOutput    bool
}

func Constants() ConstantsModel {
	file, err := os.Stat("debug.txt")

	return ConstantsModel{
		CommandPrefix:  "!",
		ConfigFilePath: "configData.json",
		PartyOnLink:    "https://giphy.com/gifs/chuber-wayne-waynes-world-d3mlYwpf96kMuFjO",
		StatusMessage:  "with your heart <3",
		DebugOutput:    !os.IsNotExist(err) && !file.IsDir(),
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

type DiscordPermission uint

const (
	PERMISSION_CREATE_INSTANT_INVITE DiscordPermission = 1 << iota
	PERMISSION_KICK_MEMBERS
	PERMISSION_BAN_MEMBERS
	PERMISSION_ADMINISTRATOR
	PERMISSION_MANAGE_CHANNELS
	PERMISSION_MANAGE_GUILD
	PERMISSION_ADD_REACTIONS
	PERMISSION_VIEW_AUDIT_LOG
	PERMISSION_VIEW_CHANNEL
	PERMISSION_SEND_MESSAGES
	PERMISSION_SEND_TTS_MESSAGES
	PERMISSION_MANAGE_MESSAGES
	PERMISSION_EMBED_LINKS
	PERMISSION_ATTACH_FILES
	PERMISSION_READ_MESSAGE_HISTORY
	PERMISSION_MENTION_EVERYONE
	PERMISSION_USE_EXTERNAL_EMOJIS
	PERMISSION_CONNECT
	PERMISSION_SPEAK
	PERMISSION_MUTE_MEMBERS
	PERMISSION_DEAFEN_MEMBERS
	PERMISSION_MOVE_MEMBERS
	PERMISSION_USE_VAD
	PERMISSION_PRIORITY_SPEAKER
	PERMISSION_STREAM
	PERMISSION_CHANGE_NICKNAME
	PERMISSION_MANAGE_NICKNAMES
	PERMISSION_MANAGE_ROLES
	PERMISSION_MANAGE_WEBHOOKS
	PERMISSION_MANAGE_EMOJIS
)
