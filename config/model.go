package config

import (
	"time"

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

// GuildData : All Data for a particular Guild (Discord Server / LAN Chapter)
type GuildData struct {
	GuildID            string `json:"guildID"`
	AuthorizedUserID   string `json:"AuthorizedUserID"`
	LanFestURL         string `json:"lanFestURL"`
	NewsURL            string `json:"newsURL"`
	AnnounceChannelID  string `json:"announceChannelID"`
	AttendeeRoleID     string `json:"attendeeRoleID"`
	PastAttendeeRoleID string `json:"pastAttendeeRoleID"`
	LFGCategoryID      string `json:"lfgCategoryID"`
	LANMode            bool
	NextLANData        LANPartyData   `json:"nextLANData"`
	PastLANData        []LANPartyData `json:"pastLANData"`
	LFGData            []LFGData      `json:"lfgData"`
}

// LANPartyData : Information on the next upcoming LANParty
type LANPartyData struct {
	Name             string    `json:"name"`
	StartDate        time.Time `json:"startDate"`
	ActivateDate     time.Time `json:"activateDate"`
	Capacity         int       `json:"capacity"`
	TicketURL        string    `json:"ticketURL"`
	TicketDate       time.Time `json:"ticketDate"`
	LastAnnounceDate time.Time `json:"lastAnnounceDate"`
}

// LFGData : Information on a current LFG setting
type LFGData struct {
	ChannelID  string    `json:"channelID"`
	Capacity   int       `json:"capacity"`
	OwnerID    string    `json:"ownerID"`
	CreateDate time.Time `json:"createDate"`
}

// OwnerSetups : Dictionary of Guilds to be Set Up by a user.
type OwnerSetups struct {
	OwnerID     string
	GuildSetups []GuildSetupData
}

// GuildSetupData : Data structure for guilds that have not been completely set up yet.
type GuildSetupData struct {
	GuildID   string
	SetupStep GuildSetupStep
}

// GlobalDataModel : Struct for internal data used across the bot.
type GlobalDataModel struct {
	Bot                    *discordgo.User
	GuildData              []GuildData
	Token                  string
	DMCommandHandlers      []interface{}
	GuildCommandHandlers   []interface{}
	ReactionAddHandlers    []interface{}
	ReactionDeleteHandlers []interface{}
	Owner                  *discordgo.User
	Session                *discordgo.Session
	OwnerSetups            []OwnerSetups
}

// GuildSetupStep : Ordered setup steps for a guild
type GuildSetupStep int

// GuildStepStep values
const (
	GuildSetupStepConfirmAuthorizedUser GuildSetupStep = iota
	GuildSetupStepChapterURL
	GuildSetupStepNewsletterURL
	GuildSetupStepAnnouncementChannel
	GuildSetupStepAttendeeRole
	GuildSetupStepPastAttendeeRole
	GuildSetupStepComplete
)
