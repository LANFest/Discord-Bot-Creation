package config

import (
	"time"
)

// GuildData : All Data for a particular Guild (Discord Server / LAN Chapter)
type GuildData struct {
	GuildID            string `json:"guildId"`
	LanFestURL         string `json:"lanFestURL"`
	NewsURL            string `json:"newsURL"`
	AnnounceChannelID  string `json:"announceChannelID"`
	AttendeeRoleID     string `json:"attendeeRoleID"`
	PastAttendeeRoleID string `json:"pastAttendeeRoleID"`
	LANMode            bool
	NextLANData        LANPartyData   `json:"nextLANData"`
	PastLANData        []LANPartyData `json:"pastLANData"`
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
