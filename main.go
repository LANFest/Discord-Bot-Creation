package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/config"

	"github.com/bwmarrin/discordgo"
)

var (
	commandPrefix string
	bot           *discordgo.User
	guilds        []config.GuildData
)

var configFilePath = "configData.json"

func main() {
	discord, err := discordgo.New("Bot NTA1MDg3MDMzOTA4Mzk2MDM3.DvRjqw.2qPeYEkKHtPVBHIOP7Oy6Pkcklk")
	errCheck("error creating discord session", err)
	user, err := discord.User("@me")
	errCheck("error retrieving account", err)

	bot = user
	discord.AddHandler(coreMessageHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "Insaniquarium!")
		if err != nil {
			fmt.Println("Error attempting to set my status")
		}
		servers := discord.State.Guilds
		fmt.Printf("Test Bot has started on %d servers", len(servers))
	})

	readConfig(&guilds)

	err = discord.Open()
	errCheck("Error opening connection to Discord", err)
	defer discord.Close()

	commandPrefix = "!"

	<-make(chan struct{})

}

func errCheck(msg string, err error) {
	if err != nil {
		fmt.Printf("%s: %+v", msg, err)
		panic(err)
	}
}

func coreMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == bot.ID || user.Bot {
		//Do nothing because the bot is talking
		return
	}

	partyOnCommandHandler(session, message)
	writeConfigCommandHandler(session, message)

	fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
}

func partyOnCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if strings.HasPrefix(message.Content, "!partyon ") {
		guild := findGuildByID(&guilds, message.GuildID)
		if guild.AttendeeRoleID == "" {
			session.ChannelMessageSend(message.ChannelID, "We can't party!  I don't know what role to assign.")
			return
		}

		if len(message.Mentions) == 0 {
			session.ChannelMessageSend(message.ChannelID, "You have to tell me who is ready to party!")
			return
		}

		if (*findRole(session, message.GuildID, guild.AttendeeRoleID)).ID == "" {
			session.ChannelMessageSend(message.ChannelID, "Invalid Attendee Role Stored!")
			return
		}

		var responseMessage = "Party On "
		for _, mention := range message.Mentions {
			err := session.GuildMemberRoleAdd(message.GuildID, mention.ID, guild.AttendeeRoleID)
			errCheck("Unable to add role!", err)
			responseMessage += "<@" + mention.ID + ">! "
		}
		responseMessage += "https://giphy.com/gifs/chuber-wayne-waynes-world-d3mlYwpf96kMuFjO"

		session.ChannelMessageSend(message.ChannelID, responseMessage)
	}
}

func writeConfigCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	// We have no way to know if the owner of the bot is the one calling writeConfig.
	// Should look into this in the future.
	if strings.HasPrefix(message.Content, "!writeConfig") {
		file, _ := json.MarshalIndent(guilds, "", " ")
		error := ioutil.WriteFile(configFilePath, file, 0644)
		errCheck("Error writing config data!", error)

		session.ChannelMessageSend(message.ChannelID, "Config Data Recorded")
	}
}

func findGuildByID(guilds *[]config.GuildData, targetGuildID string) *config.GuildData {
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

func readConfig(guilds *[]config.GuildData) {
	file, readError := ioutil.ReadFile(configFilePath)
	errCheck("Error reading Config Data", readError)

	parseError := json.Unmarshal(file, guilds)
	errCheck("Error parsing Config Data", parseError)
}

func findRole(session *discordgo.Session, guildID string, roleID string) *discordgo.Role {
	tempGuild, _ := session.Guild(guildID)
	for _, role := range tempGuild.Roles {
		if role.ID == roleID {
			return role
		}
	}
	return new(discordgo.Role)
}
