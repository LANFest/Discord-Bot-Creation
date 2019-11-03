package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/config"
	"github.com/LANFest/Discord-Bot-Creation/utils"

	"github.com/bwmarrin/discordgo"
)

var (
	commandPrefix   string
	bot             *discordgo.User
	guilds          []config.GuildData
	token           string
	commandHandlers []interface{}
	owner           *discordgo.User
)

func main() {
	file, readError := ioutil.ReadFile("token.txt")
	utils.Assert("Error reading token file", readError)

	token = string(file)

	discord, err := discordgo.New("Bot " + token)
	utils.Assert("error creating discord session", err)

	discord.AddHandler(coreMessageHandler)
	discord.AddHandler(func(discord *discordgo.Session, ready *discordgo.Ready) {
		err = discord.UpdateStatus(0, "Insaniquarium!")
		if err != nil {
			fmt.Println("Error attempting to set my status")
		}

		// Who are we?
		bot = discord.State.Ready.User
		fmt.Print("Bot Connected!\n")

		// Who's the owner?
		application, appError := discord.Application("@me")
		utils.Assert("Could not find application!", appError)
		owner = application.Owner
		fmt.Printf("Application: %s - Owner: %s\n", application.Name, application.Owner.String())
		fmt.Printf("User: %s -  ID: %s\n-----------------------------------------\n", bot.String(), bot.ID)

		// Where are we?
		servers := discord.State.Guilds
		fmt.Printf("Servers (%d):\n", len(servers))
		for _, server := range servers {
			fmt.Printf("%s - %s\n", server.Name, server.ID)
		}

	})

	utils.ReadConfig(&guilds)

	// Set up command handlers
	commandHandlers = []interface{}{partyOnCommandHandler, writeConfigCommandHandler, shutdownCommandHandler}

	err = discord.Open()
	utils.Assert("Error opening connection to Discord", err)
	defer discord.Close()

	commandPrefix = "!"

	<-make(chan struct{})

}

func coreMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == bot.ID || user.Bot {
		//Do nothing because a bot is talking
		return
	}

	if !strings.HasPrefix(message.Content, commandPrefix) {
		// It's not a command, nothing to do here.
		return
	}

	for _, handler := range commandHandlers {
		handler.(func(*discordgo.Session, *discordgo.MessageCreate))(session, message)
	}

	fmt.Printf("Message: %+v || From: %s\n", message.Message, message.Author)
}

func partyOnCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if strings.HasPrefix(message.Content, commandPrefix+"partyon ") {
		guild := utils.FindGuildByID(&guilds, message.GuildID)
		if guild.AttendeeRoleID == "" {
			session.ChannelMessageSend(message.ChannelID, "We can't party!  I don't know what role to assign.")
			return
		}

		if len(message.Mentions) == 0 {
			session.ChannelMessageSend(message.ChannelID, "You have to tell me who is ready to party!")
			return
		}

		if (*utils.FindRole(session, message.GuildID, guild.AttendeeRoleID)).ID == "" {
			session.ChannelMessageSend(message.ChannelID, "Invalid Attendee Role Stored!")
			return
		}

		var responseMessage = "Party On "
		for _, mention := range message.Mentions {
			err := session.GuildMemberRoleAdd(message.GuildID, mention.ID, guild.AttendeeRoleID)
			utils.Assert("Unable to add role!", err)
			responseMessage += "<@" + mention.ID + ">! "
		}
		responseMessage += "https://giphy.com/gifs/chuber-wayne-waynes-world-d3mlYwpf96kMuFjO"

		session.ChannelMessageSend(message.ChannelID, responseMessage)
	}
}

func writeConfigCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if isOwner(message.Author) &&
		isDM(session, message.Message) &&
		strings.HasPrefix(message.Content, commandPrefix+"writeConfig") {

		utils.WriteConfig(&guilds)
		session.ChannelMessageSend(message.ChannelID, "Config Data Recorded")
	}
}

func shutdownCommandHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if isOwner(message.Author) &&
		isDM(session, message.Message) &&
		strings.HasPrefix(message.Content, commandPrefix+"shutdown") {

		session.ChannelMessageSend(message.ChannelID, "Buh-bye!")
		utils.Shutdown(session)
	}
}

func isOwner(user *discordgo.User) bool {
	return user.ID == owner.ID
}

func isDM(session *discordgo.Session, message *discordgo.Message) bool {
	channel, _ := session.Channel(message.ChannelID)
	return channel.Type == discordgo.ChannelTypeDM
}
