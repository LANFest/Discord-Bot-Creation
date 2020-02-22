package config

import (
	"time"
)

// GuildData : All Data for a particular Guild (Discord Server / LAN Chapter)
type GuildData struct {
	GuildID            string `json:"guildID"`
	LANFestURL         string `json:"lanFestURL"`
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

type LFGData struct {
	ChannelID  string    `json:"channelID"`
	Capacity   int       `json:"capacity"`
	OwnerID    string    `json:"ownerID"`
	CreateDate time.Time `json:"createDate"`
}
