package main

import (
	"io/ioutil"
	"strings"

	"github.com/LANFest/Discord-Bot-Creation/commands/admin"
	"github.com/LANFest/Discord-Bot-Creation/commands/chapter"
	"github.com/LANFest/Discord-Bot-Creation/data"
	"github.com/LANFest/Discord-Bot-Creation/utils"

	"github.com/bwmarrin/discordgo"
)

func main() {
	globalData := data.Globals()
	file, readError := ioutil.ReadFile("token.txt")
	utils.Assert("Error reading token file", readError, true)

	globalData.Token = string(file)

	discord, err := discordgo.New("Bot " + globalData.Token)
	utils.Assert("error creating discord session", err, true)

	discord.AddHandler(coreMessageHandler)
	discord.AddHandler(coreReadyHandler)

	utils.ReadConfig()

	// Set up command handlers
	globalData.CommandHandlers = []interface{}{chapter.PartyOnCommandHandler, admin.WriteConfigCommandHandler, admin.ShutdownCommandHandler}

	err = discord.Open()
	utils.Assert("Error opening connection to Discord", err, true)
	defer discord.Close()

	<-make(chan struct{})
}

func coreMessageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.Bot {
		//Do nothing because a bot is talking
		return
	}

	if !strings.HasPrefix(message.Content, data.Constants().CommandPrefix) {
		// It's not a command, nothing to do here.
		return
	}

	for _, handler := range data.Globals().CommandHandlers {
		handler.(func(*discordgo.Session, *discordgo.MessageCreate))(session, message)
	}

	utils.LPrintf("Message: %+v || From: %s\n", message.Message, message.Author)
}

func coreReadyHandler(discord *discordgo.Session, ready *discordgo.Ready) {
	globalData := data.Globals()
	globalData.Session = discord
	err := discord.UpdateStatus(0, data.Constants().StatusMessage)
	if err != nil {
		utils.LPrint("Error attempting to set my status")
	}

	// Who are we?
	globalData.Bot = discord.State.Ready.User
	utils.LPrint("Bot Connected!")

	// Who's the owner?
	application, appError := discord.Application("@me")
	utils.Assert("Could not find application!", appError, true)
	globalData.Owner = application.Owner
	utils.LPrintf("Application: %s - Owner: %s", application.Name, application.Owner.String())
	utils.LPrintf("User: %s -  ID: %s\n-----------------------------------------", globalData.Bot.String(), globalData.Bot.ID)

	// Where are we?
	servers := discord.State.Guilds
	utils.LPrintf("Servers (%d):", len(servers))
	for _, server := range servers {
		utils.LPrintf("%s - %s", server.Name, server.ID)
	}

	utils.WriteConfig()
}
