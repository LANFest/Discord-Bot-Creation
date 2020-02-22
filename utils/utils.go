package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ahmetb/go-linq/v3"
	"github.com/bwmarrin/discordgo"

	"github.com/LANFest/Discord-Bot-Creation/config"
)

// ReadConfig : Reads in the config data from the file and populates the supplied pointer
func ReadConfig() {
	file, readError := ioutil.ReadFile(config.Constants().ConfigFilePath)
	Assert("Error reading Config Data", readError, false)

	parseError := json.Unmarshal(file, &config.Globals().GuildData)
	Assert("Error parsing Config Data", parseError, false)
}

// WriteConfig : Writes the config data to disk
func WriteConfig() {
	file, _ := json.MarshalIndent(config.Globals().GuildData, "", " ")
	error := ioutil.WriteFile(config.Constants().ConfigFilePath, file, 0644)
	Assert("Error writing config data!", error, false)
}

// FindRole : finds the requested Role within the list of Roles from the specified Guild
func FindRole(guildID string, roleID string) *discordgo.Role {
	tempGuild, guildError := config.Globals().Session.Guild(guildID)
	if guildError != nil {
		LogErrorf("FindRole", "Missing Guild: %s", guildID)
		return nil
	}

	role, ok := linq.From(tempGuild.Roles).FirstWithT(func(r *discordgo.Role) bool { return r.ID == roleID }).(*discordgo.Role)
	if !ok {
		LogErrorf("FindRole", "Missing Role %s in Guild %s", roleID, guildID)
		return nil
	}

	return role
}

// FindChannelByName : Finds a Channel by name
func FindChannelByName(guild *discordgo.Guild, channelType discordgo.ChannelType, channelName string) *discordgo.Channel {
	// Looking for a text channel and forgot the #? Let's add it.
	if channelType == discordgo.ChannelTypeGuildText && !strings.HasPrefix(channelName, "#") {
		channelName = fmt.Sprintf("#%s", channelName)
	}

	channel, ok := linq.From(guild.Channels).FirstWithT(func(c *discordgo.Channel) bool {
		return c.Type == channelType && strings.ToLower(c.Name) == strings.ToLower(channelName)
	}).(*discordgo.Channel)

	if !ok {
		LogErrorf("FindChannelByName", "Could not find text Channel named |%s| in Guild %s", channelName, guild.ID)
		return nil
	}

	return channel
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

// Shutdown : Shuts down the bot
func Shutdown(session *discordgo.Session) {
	LPrint("Shutting Down!")
	session.Logout()
	session.Close()
	os.Exit(0)
}

// IsOwner : Are you my daddy?
func IsOwner(user *discordgo.User) bool {
	return user.ID == config.Globals().Owner.ID
}

// IsDM : Is this message a DM?
func IsDM(message *discordgo.Message) bool {
	return IsChannelType(message, discordgo.ChannelTypeDM) || IsChannelType(message, discordgo.ChannelTypeGroupDM)
}

// IsGuildMessage : Is this message a Guild Text message?
func IsGuildMessage(message *discordgo.Message) bool {
	return IsChannelType(message, discordgo.ChannelTypeGuildText)
}

// IsChannelType : Is this message from the specified channel type?
func IsChannelType(message *discordgo.Message, chanType discordgo.ChannelType) bool {
	channel, _ := config.Globals().Session.Channel(message.ChannelID)
	return channel.Type == chanType
}

// LPrint : wrapper for fmt.Print that appends a newline
func LPrint(message string) {
	fmt.Print(message + "\n")
}

// LPrintf : wrapper for fmt.Printf that appends a newline
func LPrintf(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
}

// HasGuildPermission : Do I have the specified permission?
func HasGuildPermission(session *discordgo.Session, guildID string, permissionMask uint) bool {
	// Get the guild member
	guildMember, guildErr := session.GuildMember(guildID, config.Globals().Bot.ID)
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
func SendDMToUser(userID string, message string) {
	session := config.Globals().Session
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

// LogErrorf : LogError but with formatting!
func LogErrorf(componentName string, format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	LogError(componentName, message)
}

// LogError : Logs an error
func LogError(componentName string, message string) {
	LPrintf("$s: %s", componentName, message)
}
