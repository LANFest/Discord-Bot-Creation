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
)

func main() {
	file, readError := ioutil.ReadFile("token.txt")
	utils.Assert("Error reading token file", readError)

	token = string(file)

	discord, err := discordgo.New("Bot " + token)
	utils.Assert("error creating discord session", err)
	user, err := discord.User("@me")
	utils.Assert("error retrieving account", err)

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

	utils.ReadConfig(&guilds)

	// Set up command handlers
	commandHandlers = []interface{}{partyOnCommandHandler, writeConfigCommandHandler}

	err = discord.Open()
	utils.Assert("Error opening connection to Discord", err)
	defer discord.Close()

	commandPrefix = "!"

	<-make(chan struct{})

}

func coreMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	user := message.Author
	if user.ID == bot.ID || user.Bot {
		//Do nothing because the bot is talking
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
	// We have no way to know if the owner of the bot is the one calling writeConfig.
	// Should look into this in the future.
	if strings.HasPrefix(message.Content, commandPrefix+"writeConfig") {
		utils.WriteConfig(&guilds)
		session.ChannelMessageSend(message.ChannelID, "Config Data Recorded")
	}
}
