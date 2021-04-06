package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"
)

// FindRole : finds the requested Role within the list of Roles from the specified Guild
func FindRole(session *discordgo.Session, guildID string, roleID string) *discordgo.Role {
	tempGuild, guildError := session.Guild(guildID)
	if guildError != nil {
		log.Printf("Missing Guild: %s", guildID)
		return nil
	}

	role, ok := linq.From(tempGuild.Roles).FirstWithT(func(r *discordgo.Role) bool { return r.ID == roleID }).(*discordgo.Role)
	if !ok {
		log.Printf("Missing Role %s in Guild %s", roleID, guildID)
		return nil
	}

	return role
}

// FindRoleByName : Finds a Role by name (or nil if not found)
func FindRoleByName(guild *discordgo.Guild, roleName string) *discordgo.Role {
	role, ok := linq.From(guild.Roles).FirstWithT(func(r *discordgo.Role) bool { return strings.ToLower(r.Name) == strings.ToLower(roleName) }).(*discordgo.Role)
	if ok {
		return role
	}
	return nil
}

// FindChannelByName : Finds a Channel by name
func FindChannelByName(guild *discordgo.Guild, channelType discordgo.ChannelType, channelName string) *discordgo.Channel {
	// Looking for a text channel? Let's helpfully remove that pesky #
	if channelType == discordgo.ChannelTypeGuildText && strings.HasPrefix(channelName, "#") {
		channelName = strings.TrimPrefix(channelName, "#")
	}

	channel, ok := linq.From(guild.Channels).FirstWithT(func(c *discordgo.Channel) bool {
		return c.Type == channelType && strings.ToLower(c.Name) == strings.ToLower(channelName)
	}).(*discordgo.Channel)

	if !ok {
		log.Printf("Could not find text Channel named |%s| in Guild %s", channelName, guild.ID)
		return nil
	}

	return channel
}

// IsDM : Is this message a DM?
func IsDM(session *discordgo.Session, message *discordgo.Message) bool {
	return IsChannelType(session, message, discordgo.ChannelTypeDM) || IsChannelType(session, message, discordgo.ChannelTypeGroupDM)
}

// IsGuildMessage : Is this message a Guild Text message?
func IsGuildMessage(session *discordgo.Session, message *discordgo.Message) bool {
	return IsChannelType(session, message, discordgo.ChannelTypeGuildText)
}

// IsChannelType : Is this message from the specified channel type?
func IsChannelType(session *discordgo.Session, message *discordgo.Message, chanType discordgo.ChannelType) bool {
	channel, _ := session.Channel(message.ChannelID)
	return channel.Type == chanType
}

// HasGuildPermission : Do I have the specified permission?
func HasGuildPermission(session *discordgo.Session, guildID string, permissionMask uint) bool {
	// Get the guild member
	guildMember, guildErr := session.GuildMember(guildID, "@me")
	if guildErr != nil {
		return false
	}

	guildRoles, guildErr := session.GuildRoles(guildID)
	if guildErr != nil {
		return false
	}

	// linq's Intersect can't handle disparate types, so we build ourselves a collection.
	var myRoles []*discordgo.Role
	linq.From(guildRoles).WhereT(func(r *discordgo.Role) bool {
		return linq.From(guildMember.Roles).AnyWithT(func(r2 string) bool {
			return r.ID == r2
		})
	}).ToSlice(&myRoles)

	// Walk through the list, figure out if we've got the permission.
	for _, role := range myRoles {
		if uint(role.Permissions)&permissionMask == permissionMask {
			return true
		}
	}

	return false
}

// SendDMToUser : Sends a message to a user (opening a DM if none exist.)
func SendDMToUser(session *discordgo.Session, userID string, message string) {
	channels, _ := session.UserChannels()
	channel, ok := linq.From(channels).FirstWithT(func(uc *discordgo.Channel) bool {
		recips := uc.Recipients
		return len(recips) == 1 && recips[0].ID == userID
	}).(*discordgo.Channel)

	if !ok {
		newChannel, channelErr := session.UserChannelCreate(userID)

		if channelErr != nil {
			Assert(fmt.Sprintf("Unable to create DM channel with %s in SendMessageToUser!", userID), channelErr, false)
			return
		}

		channel = newChannel
	}

	session.ChannelMessageSend(channel.ID, message)
}
