package utils

import (
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
)

var session = &discordgo.Session{State: discordgo.NewState(), StateEnabled: true}

var validGuild = &discordgo.Guild{
	ID: "12345",
	Roles: []*discordgo.Role{
		{
			ID:   "67890",
			Name: "valid",
		},
	},
	Channels: []*discordgo.Channel{
		{
			ID:   "12345",
			Name: "general",
			Type: discordgo.ChannelTypeGuildText,
		},
	},
}

func TestFindRoleValid(t *testing.T) {
	session.State.GuildAdd(validGuild)

	foundRole := FindRole(session, validGuild.ID, validGuild.Roles[0].ID)
	if foundRole == nil {
		t.Error()
	}

	session.State.GuildRemove(validGuild)
}

// We can test nil states as long as our Guild exists.
func TestFindRoleInvalidRole(t *testing.T) {
	session.State.GuildAdd(validGuild)

	foundRole := FindRole(session, validGuild.ID, "lolnope")
	if foundRole != nil {
		t.Error()
	}

	session.State.GuildRemove(validGuild)
}

func TestFindRoleByNameValid(t *testing.T) {
	foundRole := FindRoleByName(validGuild, validGuild.Roles[0].Name)
	if foundRole == nil {
		t.Error()
	}
}

func TestFindRoleByNameInValid(t *testing.T) {
	foundRole := FindRoleByName(validGuild, "lolnope")
	if foundRole != nil {
		t.Error()
	}
}

func TestFindChannelByNameValid(t *testing.T) {
	foundChannel := FindChannelByName(validGuild, discordgo.ChannelTypeGuildText, validGuild.Channels[0].Name)
	if foundChannel == nil {
		t.Error()
	}
}

func TestFindChannelByNameValidWithHash(t *testing.T) {
	foundChannel := FindChannelByName(validGuild, discordgo.ChannelTypeGuildText, "#"+validGuild.Channels[0].Name)
	if foundChannel == nil {
		t.Error()
	}
}

func TestFindChannelByNameValidWithUpper(t *testing.T) {
	foundChannel := FindChannelByName(validGuild, discordgo.ChannelTypeGuildText, strings.ToUpper(validGuild.Channels[0].Name))
	if foundChannel == nil {
		t.Error()
	}
}

func TestFindChannelByNameInvalid(t *testing.T) {
	foundChannel := FindChannelByName(validGuild, discordgo.ChannelTypeGuildText, "lolnope")
	if foundChannel != nil {
		t.Error()
	}
}

func TestFindChannelByNameChannelTypeInvalid(t *testing.T) {
	foundChannel := FindChannelByName(validGuild, discordgo.ChannelTypeGuildNews, validGuild.Channels[0].Name)
	if foundChannel != nil {
		t.Error()
	}
}
